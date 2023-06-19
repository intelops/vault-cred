package policy

import (
	"context"

	"strings"
	"time"

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
func (p *VaultPolicyHandler) GetConfigMapTime(ctx context.Context, configMapName, namespace string) (time.Time, error) {
	k8s, err := client.NewK8SClient(p.log)
	if err != nil {
		return time.Time{}, err
	}

	allConfigMapData, err := k8s.GetConfigMapCreationTimestamp(ctx, configMapName, namespace)
	if err != nil {
		return time.Time{}, errors.WithMessagef(err, "error while getting vault policy configmaps")
	}
	return allConfigMapData, nil
}
func extractNamespaceAndName(cmname string) (string, string) {
	parts := strings.SplitN(cmname, ":", 2)
	if len(parts) != 2 {
		return "", ""
	}
	namespace := parts[0]
	configname := parts[1]
	return namespace, configname
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

	for cmname, cmData := range allConfigMapData {
		policyName := cmData["policyName"]
		policyData := cmData["policyData"]

		nsname, configname := extractNamespaceAndName(cmname)

		policyExists, err := vc.CheckVaultPolicyExists(policyName)

		if err != nil {
			p.log.Errorf("Error while checking if Vault policy exists: %v", err)
			continue
		}
		originalTimestamp, err := p.GetConfigMapTime(ctx, configname, nsname)
		if err != nil {
			p.log.Errorf("Error while getting timestamp", err)
		}
		if !policyExists {
			// Create Vault policy if it doesn't exist
			err = vc.CreateOrUpdatePolicy(policyName, policyData)
			if err != nil {
				p.log.Errorf("Error while creating Vault policy %s: %v", policyName, err)
				continue
			}
			p.log.Infof("Created Vault policy %s", policyName)
		} else {
			//	p.log.Infof("Vault policy %s already exists", policyName)
		}

		updatedTimestamp, err := p.GetConfigMapTime(ctx, configname, nsname)
		if err != nil {
			p.log.Errorf("Error while checking timestamp: %v", err)
			continue
		}

		if originalTimestamp != updatedTimestamp {
			// Update the Vault policy based on the updated ConfigMap
			err = vc.CreateOrUpdatePolicy(policyName, policyData)
			if err != nil {
				return errors.WithMessagef(err, "error while creating vault policy %s, %v", policyName, cmData)
			}
			p.log.Infof("Updated Vault policy %s", policyName)
		} else {
			//p.log.Infof("No update needed for Vault policy %s", policyName)
			continue
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
	for cmname, cmData := range allConfigMapData {
		roleName := cmData["roleName"]
		policyNames := cmData["policyNames"]
		servieAccounts := cmData["servieAccounts"]
		servieAccountNameSpaces := cmData["servieAccountNameSpaces"]
		nsname, configname := extractNamespaceAndName(cmname)
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
		exists, err := vc.RoleExists(roleName)
		if err != nil {
			return err
		}
		originalTimestamp, err := p.GetConfigMapTime(ctx, configname, nsname)
		if err != nil {
			p.log.Errorf("Error while getting timestamp", err)
		}
		if !exists {
			p.log.Infof("Vault role  %s not  already exists", roleName)
			err = vc.CreateOrUpdateRole(roleName, policyNameList,
				strings.Split(servieAccounts, ","),
				strings.Split(servieAccountNameSpaces, ","))
			if err != nil {
				return errors.WithMessagef(err, "error while creating vault role %s", roleName)
			}

		} else {
			p.log.Infof("Vault role %s already exists", roleName)
		}

		updatedTimestamp, err := p.GetConfigMapTime(ctx, configname, nsname)
		if err != nil {
			p.log.Errorf("Error while checking timestamp: %v", err)
			continue
		}

		if originalTimestamp != updatedTimestamp {
			// Update the Vault policy based on the updated ConfigMap
			err = vc.CreateOrUpdateRole(roleName, policyNameList,
				strings.Split(servieAccounts, ","),
				strings.Split(servieAccountNameSpaces, ","))
			if err != nil {
				p.log.Errorf("Error while updating Vault role %s: %v", roleName, err)
				continue
			}
			p.log.Infof("Updated Vault role %s", roleName)
		} else {
			p.log.Infof("No update needed for Vault role %s", roleName)
		}

		//
	}
	return nil
}
