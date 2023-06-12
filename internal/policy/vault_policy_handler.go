package policy

import (
	"context"
	"strings"

	"github.com/intelops/go-common/logging"
	"github.com/intelops/vault-cred/config"
	"github.com/intelops/vault-cred/internal/client"
	"github.com/pkg/errors"
)

type VaultPolicyHandler struct {
	log      logging.Logger
	conf     config.VaultEnv
	policies map[string]struct{}
}

func NewVaultPolicyHandler(log logging.Logger, conf config.VaultEnv) *VaultPolicyHandler {
	return &VaultPolicyHandler{log: log, conf: conf,
		policies: map[string]struct{}{}}
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

	p.policies = map[string]struct{}{}
	for _, cmData := range allConfigMapData {
		policyName := cmData["policyName"]
		policyData := cmData["policyData"]
		err = vc.CreateOrUpdatePolicy(policyName, policyData)
		if err != nil {
			return errors.WithMessagef(err, "error while creating vault policy %s, %v", policyName, cmData)
		}
		p.policies[policyName] = struct{}{}
	}
	return nil
}

func (p *VaultPolicyHandler) UpdateVaultRoles(ctx context.Context) error {
	vc, err := client.NewVaultClientForVaultToken(p.log, p.conf)
	if err != nil {
		return err
	}

	allConfigMapData, err := p.getVaultConfigMaps(ctx, "vault-role-")
	if err != nil {
		return errors.WithMessagef(err, "error while getting vault role configmaps")
	}

	for _, cmData := range allConfigMapData {
		roleName := cmData["roleName"]
		policyNames := cmData["policyNames"]
		servieAccounts := cmData["servieAccounts"]
		servieAccountNameSpaces := cmData["servieAccountNameSpaces"]

		policyNameList := strings.Split(policyNames, ",")
		policiesExist := true
		for _, policyName := range policyNameList {
			if _, ok := p.policies[policyName]; !ok {
				policiesExist = false
				break
			}
		}

		if !policiesExist {
			p.log.Errorf("all polices are not exist to map to the role, %v", cmData)
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