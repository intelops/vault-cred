package policy

import (
	"context"
	"fmt"

	"strings"

	"github.com/intelops/go-common/logging"
	"github.com/intelops/vault-cred/internal/client"
	"github.com/pkg/errors"
)

type VaultPolicyHandler struct {
	log               logging.Logger
	policyConfigCache vaultConfigData
	roleConfigCache   vaultConfigData
}

func NewVaultPolicyHandler(log logging.Logger) *VaultPolicyHandler {
	return &VaultPolicyHandler{log: log,
		policyConfigCache: newVaultConfigMapCache(),
		roleConfigCache:   newVaultConfigMapCache()}
}

func (p *VaultPolicyHandler) getVaultConfigMaps(ctx context.Context, prefix string) ([]client.ConfigMapData, error) {
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

func (p *VaultPolicyHandler) UpdateVaultPolicies(ctx context.Context, vc *client.VaultClient) error {
	allConfigMapData, err := p.getVaultConfigMaps(ctx, "vault-policy-")
	if err != nil {
		return errors.WithMessagef(err, "error while getting vault policy configmaps")
	}

	p.log.Debugf("found %d policy config maps", len(allConfigMapData))
	for _, cmData := range allConfigMapData {
		policyName := cmData.Data["policyName"]
		policyData := cmData.Data["policyData"]

		cmKey := fmt.Sprintf("%s:%s", cmData.Name, cmData.Namespace)
		existingCmData, found := p.policyConfigCache.Get(cmKey)
		if found {
			if existingCmData.LastUpdatedTime != cmData.LastUpdatedTime {
				err = vc.CreateOrUpdatePolicy(policyName, policyData)
				if err != nil {
					p.log.Errorf("error while updating Vault policy %s: %v", policyName, err)
					continue
				}
				p.policyConfigCache.Put(cmKey, cmData)
			} else {
				p.log.Debugf("no update needed for vault policy %s", policyName)
			}
		} else {
			err = vc.CreateOrUpdatePolicy(policyName, policyData)
			if err != nil {
				p.log.Errorf("error while creating Vault policy %s: %v", policyName, err)
				continue
			}
			p.policyConfigCache.Put(cmKey, cmData)
		}
	}
	return nil
}

func (p *VaultPolicyHandler) UpdateVaultRoles(ctx context.Context, vc *client.VaultClient) error {
	allConfigMapData, err := p.getVaultConfigMaps(ctx, "vault-role-")
	if err != nil {
		return errors.WithMessagef(err, "error while getting vault role configmaps")
	}

	if len(allConfigMapData) == 0 {
		p.log.Debugf("no vault roles found")
		return nil
	}

	err = vc.CheckAndEnableK8sAuth()
	if err != nil {
		return err
	}

	existingPolicies, err := vc.ListPolicies()
	if err != nil {
		return err
	}

	p.log.Debugf("found %d role config maps", len(allConfigMapData))
	for _, cmData := range allConfigMapData {
		roleName := cmData.Data["roleName"]
		policyNames := cmData.Data["policyNames"]
		servieAccounts := cmData.Data["servieAccounts"]
		servieAccountNameSpaces := cmData.Data["servieAccountNameSpaces"]
		policyNameList := strings.Split(policyNames, ",")

		cmKey := fmt.Sprintf("%s:%s", cmData.Name, cmData.Namespace)
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

		existingCmData, found := p.roleConfigCache.Get(cmKey)
		if found {
			if existingCmData.LastUpdatedTime != cmData.LastUpdatedTime {
				err = vc.CreateOrUpdateRole(roleName,
					strings.Split(servieAccounts, ","),
					strings.Split(servieAccountNameSpaces, ","),
					policyNameList)
				if err != nil {
					p.log.Errorf("error while updating Vault role %s: %v", roleName, err)
					continue
				}
				p.roleConfigCache.Put(cmKey, cmData)
			} else {
				p.log.Debugf("no update needed for vault role %s", roleName)
			}
		} else {
			err = vc.CreateOrUpdateRole(roleName,
				strings.Split(servieAccounts, ","),
				strings.Split(servieAccountNameSpaces, ","),
				policyNameList)
			if err != nil {
				p.log.Errorf("error while creating Vault role %s: %v", roleName, err)
				continue
			}
			p.roleConfigCache.Put(cmKey, cmData)
		}
	}
	return nil
}

func (p *VaultPolicyHandler) EnsureKVMounted(ctx context.Context, vc *client.VaultClient) error {
	return vc.CheckAndMountKVMount("secret/")
}
