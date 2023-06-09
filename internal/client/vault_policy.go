package client

import (
	"context"
	"encoding/json"
	"fmt"

	//vault "github.com/hashicorp/vault/api"
	"log"
	//"path/filepath"
	"strings"

	"github.com/intelops/vault-cred/config"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type SAData struct {
	ServiceAccount   string
	CredentialAccess []string
}

func ReturnToken() (string, error) {
	token, err := readFileContent(config.VaultEnv{}.VaultTokenPath)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve token: %v", err)
	}

	return token, nil
}

func filterVaultConfigMaps(clientset *kubernetes.Clientset) ([]corev1.ConfigMap, error) {
	configMaps := []corev1.ConfigMap{}

	// Get all namespaces in the cluster
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %v", err)
	}

	// Iterate over namespaces and fetch ConfigMaps
	for _, ns := range namespaces.Items {

		cmList, err := clientset.CoreV1().ConfigMaps(ns.Name).List(context.TODO(), metav1.ListOptions{})

		if err != nil {
			return nil, fmt.Errorf("failed to list ConfigMaps in namespace %s: %v", ns.Name, err)
		}

		configMaps = append(configMaps, cmList.Items...)
	}

	vaultConfigMaps := []corev1.ConfigMap{}
	for _, cm := range configMaps {
		if strings.HasPrefix(cm.Name, "vault-policy") {
			vaultConfigMaps = append(vaultConfigMaps, cm)
		}

	}
	return vaultConfigMaps, nil
}

func (v VaultClient) RetrieveSA(clientset *kubernetes.Clientset) ([]SAData, string, string, error) {
	res, err := filterVaultConfigMaps(clientset)
	if err != nil {
		return nil, "", "", fmt.Errorf("error while filtering: %v", err)
	}

	var saDataList []SAData
	var roleName string
	var policyName string
	for _, cm := range res {
		serviceAccount := cm.Data["serviceAccount"]
		credentialAccess := cm.Data["credentialAccess"]
		roleName := serviceAccount + "-role"
		policyName := serviceAccount + "-policy"
		fmt.Println(roleName)
		fmt.Println(policyName)
		accessList := strings.Split(credentialAccess, "\n")
		var credentialAccessList []string
		for _, access := range accessList {
			parts := strings.SplitN(access, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				credentialAccessList = append(credentialAccessList, fmt.Sprintf("%s: %s", key, value))
			}
		}

		saData := SAData{
			ServiceAccount:   serviceAccount,
			CredentialAccess: credentialAccessList,
		}
		// fmt.Println("SAData", saData)
		// fmt.Println("CredAccess", saData.CredentialAccess)
		saDataList = append(saDataList, saData)
		policyData := preparePolicyData(saDataList)
		res, err := v.CheckVaultPolicyExists(policyName)
		if err != nil {
			return nil, "", "", fmt.Errorf("failed to checking the policy: %v", err)
		}
		if res {
			v.DeleteVaultPolicy(policyName)
			err = v.CreatePolicy(policyName, policyData)
			if err != nil {
				log.Println("Error while creating policy", err)
			}
			v.DeleteVaultRole(roleName)
			err = v.CreateVaultRole(roleName, policyName, serviceAccount)
			if err != nil {
				log.Println("Error while creating vault role", err)
			}
		} else {
			err = v.CreatePolicy(policyName, policyData)
			if err != nil {
				log.Println("Error while creating policy", err)
			}
			err = v.CreateVaultRole(roleName, policyName, serviceAccount)
			if err != nil {
				log.Println("Error while creating vault role", err)
			}

		}

	}

	fmt.Println("SADatalist", saDataList)
	return saDataList, roleName, policyName, nil
}
func preparePolicyData(saDataList []SAData) map[string][]string {
	policyData := make(map[string][]string)

	for _, saData := range saDataList {
		credentialAccessList := saData.CredentialAccess
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

	fmt.Println("PolicyDataPreparation", policyData)

	return policyData
}

func (v *VaultClient) CreatePolicy(policyName string, policyData map[string][]string) error {
	token, _ := ReturnToken()
	v.c.SetToken(token)

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

	return nil
}
func formatCapabilities(capabilities []string) string {
	formatted := make([]string, len(capabilities))
	for i, capability := range capabilities {
		formatted[i] = fmt.Sprintf("\"%s\"", capability)
	}
	return strings.Join(formatted, ", ")
}
func (v *VaultClient) CreateVaultRole(roleName, policyName, serviceAccount string) error {
	// Configure the role
	token, _ := ReturnToken()
	v.c.SetToken(token)

	roleData := map[string]interface{}{
		"bound_service_account_names":      []string{serviceAccount},
		"bound_service_account_namespaces": []string{"default"},
		"policies":                         []string{policyName},
		"token_bound_cidrs":                []string{"0.0.0.0/0"},
		"token_explicit_max_ttl":           "24h",
		"token_no_default_policy":          true,
		"token_num_uses":                   10,
		"token_period":                     "1h",
		"token_type":                       "service",
	}
	roleJSON, err := json.Marshal(roleData)
	if err != nil {
		log.Fatalf("Failed to marshal Vault role: %v", err)
	}
	fmt.Println("RoleJson", string(roleJSON))
	// Write the Vault role
	_, err = v.c.Logical().Write(fmt.Sprintf("auth/kubernetes/role/%s", roleName), roleData)
	if err != nil {
		log.Fatalf("Failed to create Vault role: %v", err)
		return err
	}

	log.Println("Vault role and policy created successfully!")
	return nil
}

func (v *VaultClient) CheckVaultPolicyExists(policyName string) (bool, error) {

	//v.c.SetToken("hvs.Meh6YP3vpWtVp2ono8PL1THQ")
	token, _ := ReturnToken()

	v.c.SetToken(token)
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
func (v VaultClient) DeleteVaultPolicy(policyName string) error {
	// Set the Vault address

	// v.c.SetToken("hvs.Meh6YP3vpWtVp2ono8PL1THQ")
	token, _ := ReturnToken()

	v.c.SetToken(token)

	// Delete the policy from Vault
	err := v.c.Sys().DeletePolicy(policyName)
	if err != nil {
		return fmt.Errorf("failed to delete Vault policy: %v", err)
	}

	log.Printf("Vault policy '%s' deleted successfully!", policyName)
	return nil
}
func (v VaultClient) DeleteVaultRole(roleName string) error {
	// Set the Vault address
	token, _ := ReturnToken()

	v.c.SetToken(token)

	//	v.client.SetToken("hvs.Meh6YP3vpWtVp2ono8PL1THQ")

	// Delete the role from Vault
	_, err := v.c.Logical().Delete(fmt.Sprintf("auth/kubernetes/role/%s", roleName))
	if err != nil {
		return fmt.Errorf("failed to delete Vault role: %v", err)
	}

	log.Printf("Vault role '%s' deleted successfully!", roleName)
	return nil
}
