package client

import (
	"fmt"
	"strings"
)

type VaultPolicyData struct {
	VaultRoleName           string              `json:"vaultRoleName"`
	PolicyName              string              `json:"policyName"`
	ServiceAccount          string              `json:"serviceAccount"`
	ServiceAccountNameSpace string              `json:"serviceAccountNameSpace"`
	CredentialAccessList    map[string][]string `json:"credentialAccessList"`
}

func (v *VaultClient) CreateOrUpdatePolicy(policyName, data string) error {
	policyData := make(map[string]interface{})
	policyData["policy"] = data

	path := fmt.Sprintf("/sys/policy/%s", policyName)
	_, err := v.c.Logical().Write(path, policyData)
	if err != nil {
		return err
	}

	v.log.Infof("Created policy: %s", policyName)
	return nil
}

func (v *VaultClient) DeletePolicy(policyName string) error {
	path := fmt.Sprintf("/sys/policy/%s", policyName)
	_, err := v.c.Logical().Delete(path)
	if err != nil {
		return err
	}
	v.log.Infof("Deleted policy: %s", policyName)
	return nil
}

func (v *VaultClient) CreateOrUpdateRole(roleName string, serviceAccounts, namespaces, policies []string) error {
	roleData := make(map[string]interface{})

	sa := strings.Join(serviceAccounts, ",")
	ns := strings.Join(namespaces, ",")
	roleData["bound_service_account_names"] = sa
	roleData["bound_service_account_namespaces"] = ns
	roleData["policies"] = policies
	roleData["max_ttl"] = 1800000

	path := fmt.Sprintf("/auth/kubernetes/role/%s", roleName)
	_, err := v.c.Logical().Write(path, roleData)
	if err != nil {
		return err
	}
	v.log.Infof("Created role mapping: %s", roleName)
	return nil
}

func (v *VaultClient) DeleteRole(roleName string) error {
	path := fmt.Sprintf("/auth/kubernetes/role/%s", roleName)
	_, err := v.c.Logical().Delete(path)
	if err != nil {
		return err
	}
	v.log.Infof("Deleted role mapping: %s", roleName)
	return nil
}
