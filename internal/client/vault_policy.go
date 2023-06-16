package client

import (
	"fmt"
	"log"
	"strings"
	"github.com/hashicorp/vault/api"
)

type VaultPolicyData struct {
	VaultRoleName           string              `json:"vaultRoleName"`
	PolicyName              string              `json:"policyName"`
	ServiceAccount          string              `json:"serviceAccount"`
	ServiceAccountNameSpace string              `json:"serviceAccountNameSpace"`
	CredentialAccessList    map[string][]string `json:"credentialAccessList"`
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
	v.c.SetToken("hvs.UFiGRSW9Y8OK2miUygIDw7Np")
	roleData := make(map[string]interface{})

	sa := strings.Join(serviceAccounts, ",")
	ns := strings.Join(namespaces, ",")
	roleData["bound_service_account_names"] = sa
	roleData["bound_service_account_namespaces"] = ns
	roleData["policies"] = policies
	roleData["max_ttl"] = 1800000
	log.Println("RoleData", roleData)
	path := fmt.Sprintf("/auth/kubernetes/role/%s", roleName)
	//
	_, err := v.c.Logical().Read(path)
	if err != nil {
		// If the role doesn't exist, create it
		_, err = v.c.Logical().Write(path, roleData)
		if err != nil {
			
			return err
		}
		v.log.Infof("Created role mapping: %s", roleName)
	} else {
		// If the role already exists, update it
		_, err = v.c.Logical().Write(path+"/", roleData)
		if err != nil {
			return err
		}
		v.log.Infof("Updated role mapping: %s", roleName)
	}
	// _, err := v.c.Logical().Write(path, roleData)
	// if err != nil {
	// 	return err
	// }
	//v.log.Infof("Created role mapping: %s", roleName)
	return nil
}
func (vc *VaultClient) RoleExists(roleName string) (bool, error) {
	// Construct the API endpoint for retrieving role details
 
    path := fmt.Sprintf("/auth/kubernetes/role/%s", roleName)
	// Send a GET request to the Vault API to retrieve role details
	_, err := vc.c.Logical().Read(path)
	if err != nil {
		// Check if the error is due to the role not existing
		if apiErr, ok := err.(*api.ResponseError); ok && apiErr.StatusCode == 404 {
			return false, nil // Role does not exist
		}
		return false, err // Error occurred while retrieving role details
	}

	return true, nil // Role exists
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

func (v *VaultClient)  CheckVaultPolicyExists(policyName string) (bool, error) {
	
	// Get the list of existing policies
	policies, err := v.c.Sys().ListPolicies()
	if err != nil {
		return false, fmt.Errorf("failed to retrieve Vault policies: %v", err)
	}

	// Check if the desired policy exists in the list
	for _, policy := range policies {
		if policy == policyName {
			fmt.Print("POlicy", policy)
			return true, nil
		}
	}

	return false, nil
}
