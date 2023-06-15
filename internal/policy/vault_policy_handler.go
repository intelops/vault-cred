package policy

import (
	"context"
	"fmt"
	"strings"

	"github.com/intelops/go-common/logging"
	"github.com/intelops/vault-cred/config"
	"github.com/intelops/vault-cred/internal/client"
	"github.com/pkg/errors"
)

type VaultPolicyHandler struct {
	log  logging.Logger
	conf config.VaultEnv
}

func NewVaultPolicyHandler(log logging.Logger, conf config.VaultEnv) *VaultPolicyHandler {
	return &VaultPolicyHandler{log: log, conf: conf}
}

func (p *VaultPolicyHandler) getVaultConfigMaps(ctx context.Context, prefix string) (map[string]map[string]string, error) {
	k8s, err := client.NewK8SClient(p.log)
	if err != nil {
		return nil, err
	}

	allConfigMapData, err := k8s.GetConfigMapsHasPrefix(ctx, prefix)
	if err != nil {
		return nil, errors.WithMessagef(err, "error while getting vault policy configmaps")
	}
	return allConfigMapData, nil
}

func (p *VaultPolicyHandler) UpdateVaultPolicies(ctx context.Context) error {
	vc, err := client.NewVaultClientForVaultToken(p.log, p.conf)
	if err != nil {
		return err
	}

	allConfigMapData, err := p.getVaultConfigMaps(ctx, "vault-policy-")
	if err != nil {
		return errors.WithMessagef(err, "error while getting vault policy configmaps")
	}
	p.log.Infof("found %d policy config maps", len(allConfigMapData))

	var credentialAccessList []string
	var saDataList []client.VaultPolicyData
	for _, cmData := range allConfigMapData {
		serviceAccount := cmData["serviceAccount"]
		credentialAccess := cmData["credentialAccess"]
		policyName := serviceAccount + "-policy"
		p.log.Info("Policy Name", policyName)
		accessList := strings.Split(credentialAccess, "\n")

		for _, access := range accessList {
			parts := strings.SplitN(access, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				credentialAccessList = append(credentialAccessList, fmt.Sprintf("%s: %s", key, value))
			}
		}

		saData := client.VaultPolicyData{
			ServiceAccount:       serviceAccount,
			CredentialAccessList: credentialAccessList,
		}
		saDataList = append(saDataList, saData)
		policydata := vc.PreparePolicyData(saDataList)
		p.log.Info("PolicyData", policydata)
		err := vc.CreateOrUpdPolicy(policyName, policydata)
		if err != nil {
			return errors.WithMessagef(err, "error while creating vault policy %s, %v", policyName, cmData)
		}
	}
	return nil
}

func (p *VaultPolicyHandler) UpdateVaultRoles(ctx context.Context) error {
	allConfigMapData, err := p.getVaultConfigMaps(ctx, "vault-role-")
	if err != nil {
		return errors.WithMessagef(err, "error while getting vault role configmaps")
	}

	if len(allConfigMapData) == 0 {
		p.log.Infof("no vault roles found %d to configure")
		return nil
	}

	vc, err := client.NewVaultClientForVaultToken(p.log, p.conf)
	if err != nil {
		return err
	}

	err = vc.CheckAndEnableK8sAuth()
	if err != nil {
		return err
	}

	existingPolicies, err := vc.ListPolicies()
	if err != nil {
		return err
	}

	p.log.Infof("found %d role config maps", len(allConfigMapData))
	for _, cmData := range allConfigMapData {
		roleName := cmData["roleName"]
		policyNames := cmData["policyNames"]
		servieAccounts := cmData["servieAccounts"]
		servieAccountNameSpaces := cmData["servieAccountNameSpaces"]

		policyNameList := strings.Split(policyNames, ",")
		policiesExist := true
		for _, policyName := range policyNameList {
			found := false
			for _, existingPolicyName := range existingPolicies {
				if existingPolicyName == policyName {
					found = true
					break
				}
			}
			if !found {
				policiesExist = false
				break
			}
		}

		if !policiesExist {
			p.log.Errorf("all polices are not exist to map to the role %s, %v", roleName, cmData)
			continue
		}

		err = vc.CreateOrUpdateRole(roleName, policyNameList,
			strings.Split(servieAccounts, ","),
			strings.Split(servieAccountNameSpaces, ","))
		if err != nil {
			return errors.WithMessagef(err, "error while creating vault role %s, %v", roleName, cmData)
		}
	}
	return nil
}
