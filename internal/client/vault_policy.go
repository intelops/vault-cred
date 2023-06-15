package client

import (
	"fmt"
	"strings"
)

type VaultPolicyData struct {
	VaultRoleName           string   `json:"vaultRoleName"`
	PolicyName              string   `json:"policyName"`
	ServiceAccount          string   `json:"serviceAccount"`
	ServiceAccountNameSpace string   `json:"serviceAccountNameSpace"`
	CredentialAccessList    []string `json:"credentialAccessList"`
	// CredentialAccessList    map[string][]string `json:"credentialAccessList"`
}

func (v *VaultClient) CreateOrUpdatePolicy(policyName, rules string) error {

	err := v.c.Sys().PutPolicy(policyName, rules)
	if err != nil {
		return err
	}

	v.log.Infof("Created policy: %s", policyName)
	return nil
}

func (v *VaultClient) DeletePolicy(policyName string) error {

	err := v.c.Sys().DeletePolicy(policyName)
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

func (v *VaultClient) ListPolicies() ([]string, error) {

	return v.c.Sys().ListPolicies()
}

func (v *VaultClient) CheckAndEnableK8sAuth() error {

	return v.c.Sys().EnableAuth("kubernetes", "kubernetes", "kubernetes authentication")
}

//

func (v *VaultClient) PreparePolicyData(saDataList []VaultPolicyData) map[string][]string {
	policyData := make(map[string][]string)

	for _, saData := range saDataList {
		credentialAccessList := saData.CredentialAccessList
		for _, credentialAccess := range credentialAccessList {
			accessPairs := strings.Split(credentialAccess, ",")

			var parts []string
			for _, pair := range accessPairs {
				parts = append(parts, strings.Split(pair, ":")...)

			}

			if len(parts) == 2 || len(parts) >= 2 {
				var accessTypes []string

				credentialType := strings.TrimSpace(parts[0])

				for i := 1; i < len(parts); i++ {
					accessTypes = append(accessTypes, strings.Split(parts[i], ",")...)
				}

				trimmedAccessTypes := make([]string, len(accessTypes))
				for i, accessType := range accessTypes {
					trimmedAccessTypes[i] = strings.TrimSpace(accessType)

				}

				if _, exists := policyData[credentialType]; exists {
					// Append the access types to the existing credential type in the policy data
					policyData[credentialType] = append(policyData[credentialType], trimmedAccessTypes...)

				} else {
					// Create a new entry for the credential type in the policy data
					policyData[credentialType] = trimmedAccessTypes

				}

			}
		}
	}

	return policyData
}
func formatCapabilities(capabilities []string) string {
	formatted := make([]string, len(capabilities))
	for i, capability := range capabilities {
		formatted[i] = fmt.Sprintf("\"%s\"", capability)
	}
	return strings.Join(formatted, ", ")
}

func (v *VaultClient) CreateOrUpdPolicy(policyName string, policyData map[string][]string) error {

	rules := ""
	for path, access := range policyData {
		rules += fmt.Sprintf("path \"secret/data/%s/*\" {\n  capabilities = [%s]\n}\n", path, formatCapabilities(access))
	}

	_, err := v.c.Logical().Write("sys/policy/"+policyName, map[string]interface{}{
		"rules": rules,
	})
	if err != nil {
		return err
	}

	fmt.Println("Policy Data", policyData)
	return nil
}
