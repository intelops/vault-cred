package client

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	//vault "github.com/hashicorp/vault/api"
	"log"
	//"path/filepath"
	"strings"

	//	"github.com/intelops/vault-cred/config"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type SAData struct {
	ServiceAccount   string
	CredentialAccess []string
}

func FilterVaultConfigMaps(clientset *kubernetes.Clientset) ([]corev1.ConfigMap, error) {
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


func PreparePolicyData(saDataList []SAData) map[string][]string {
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
	// token, _ := ReturnToken()
	// v.c.SetToken(token)

	rules := ""
	for path, access := range policyData {
		rules += fmt.Sprintf("path \"secret/data/%s/*\" {\n  capabilities = [%s]\n}\n", path, formatCapabilities(access))
	}

	_, err := v.C.Logical().Write("sys/policy/"+policyName, map[string]interface{}{
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
	// token, _ := ReturnToken()
	// v.c.SetToken(token)

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
	_, err = v.C.Logical().Write(fmt.Sprintf("auth/kubernetes/role/%s", roleName), roleData)
	if err != nil {
		log.Fatalf("Failed to create Vault role: %v", err)
		return err
	}

	log.Println("Vault role and policy created successfully!")
	return nil
}

func (v VaultClient) DeleteVaultPolicy(policyName string) error {
	// Set the Vault address

	// Delete the policy from Vault
	err := v.C.Sys().DeletePolicy(policyName)
	if err != nil {
		return fmt.Errorf("failed to delete Vault policy: %v", err)
	}

	log.Printf("Vault policy '%s' deleted successfully!", policyName)
	return nil
}
func (v VaultClient) DeleteVaultRole(roleName string) error {
	// Set the Vault address

	//	v.client.SetToken("hvs.Meh6YP3vpWtVp2ono8PL1THQ")

	// Delete the role from Vault
	_, err := v.C.Logical().Delete(fmt.Sprintf("auth/kubernetes/role/%s", roleName))
	if err != nil {
		return fmt.Errorf("failed to delete Vault role: %v", err)
	}

	log.Printf("Vault role '%s' deleted successfully!", roleName)
	return nil
}
func (v VaultClient) CheckVaultPolicyExistsAndEqual(policyName string, desiredPolicyData map[string][]string) (bool, error) {
	// Create a new Vault client

	// Set the Vault token for authentication

	// Get the existing policy rules
	policy, err := v.C.Sys().GetPolicy(policyName)
	if err != nil {
		return false, err
	}

	// Convert the existing policy rules to map[string][]string format
	existingPolicyData := make(map[string][]string)
	if policy != "" {
		existingPolicyData = convertPolicyToData(policy)
	}

	// Compare the existing policy data with the desired policy data
	areEqual := reflect.DeepEqual(existingPolicyData, desiredPolicyData)

	return areEqual, nil
}

// Helper function to convert policy rules to map[string][]string format
func convertPolicyToData(policy string) map[string][]string {
	data := make(map[string][]string)

	// Parse the policy rules line by line
	lines := strings.Split(policy, "\n")
	for _, line := range lines {
		// Skip empty lines or comments
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Extract the path and access list
		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 2 {
			path := parts[0]
			accessList := parts[1]
			data[path] = strings.Fields(accessList)
		}
	}

	return data
}
