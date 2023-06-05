package client

import (
	//"encoding/json"

	"encoding/json"
	"fmt"
	"log"

	//"path/filepath"
	"strings"
	//"golang.org/x/exp/event/keys"
	//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//	"github.com/hashicorp/vault/api"
	//
	// "k8s.io/client-go/kubernetes"
	// "k8s.io/client-go/tools/clientcmd"
	// "k8s.io/client-go/util/homedir"
)

type Policy struct {
	Path         string   `json:"path"`
	Capabilities []string `json:"capabilities"`
}

func (v VaultClient) CreateVaultPolicy(policyName string, policyRules []Policy) error {

	//	v.c.SetToken("hvs.Meh6YP3vpWtVp2ono8PL1THQ")

	//
	// Construct the policy string
	var policyBuilder strings.Builder
	for _, rule := range policyRules {
		policyBuilder.WriteString(fmt.Sprintf("path \"%s\" {\n", rule.Path))
		policyBuilder.WriteString(fmt.Sprintf("\tcapabilities = %s\n", formatCapabilities(rule.Capabilities)))
		policyBuilder.WriteString("}\n\n")
	}

	policy := policyBuilder.String()
	fmt.Println("Policy", policy)
	// Write the policy to Vault

	err := v.c.Sys().PutPolicy(policyName, policy)
	// err = vaultClient.Sys().PutPolicy(policyName, policy)
	if err != nil {
		return fmt.Errorf("failed to create Vault policy: %v", err)
	}

	log.Printf("Vault policy '%s' created successfully!", policyName)
	return nil
}

func formatCapabilities(capabilities []string) string {
	formatted := make([]string, len(capabilities))
	for i, capability := range capabilities {
		formatted[i] = fmt.Sprintf("\"%s\"", capability)
	}
	return fmt.Sprintf("[%s]", strings.Join(formatted, ", "))
}

func (v VaultClient) CheckVaultPolicyExists(policyName string) (bool, error) {

	//	v.c.SetToken("hvs.Meh6YP3vpWtVp2ono8PL1THQ")

	// Get the list of existing policies
	policies, err := v.c.Sys().ListPolicies()
	if err != nil {
		return false, fmt.Errorf("failed to retrieve Vault policies: %v", err)
	}

	// Check if the desired policy exists in the list
	for _, policy := range policies {
		if policy == policyName {
			return true, nil
		}
	}

	return false, nil
}
func (v VaultClient) CreateVaultRole(roleName, policyName, serviceAccount string) error {

	//	v.c.SetToken("hvs.Meh6YP3vpWtVp2ono8PL1THQ")

	// Configure the role

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

// func GetConfigMap() (string, string, string) {
// 	namespace := "default"
// 	configMapName := "my-configmap"
// 	clientset := Config()
// 	configMap, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), configMapName, metav1.GetOptions{})
// 	if err != nil {
// 		log.Fatalf("Failed to retrieve ConfigMap: %v", err)
// 	}

//		// Read the values from the ConfigMap data
//		serviceAccount := configMap.Data["serviceAccount"]
//		credentialTypes := configMap.Data["credentialTypes"]
//		accessType := configMap.Data["accessType"]
//		fmt.Println("Credential type0", credentialTypes)
//		fmt.Println("Access type", accessType)
//		return serviceAccount, credentialTypes, accessType
//	}
func (v VaultClient) DeleteVaultPolicy(policyName string) error {
	// Set the Vault address

	//	v.c.SetToken("hvs.Meh6YP3vpWtVp2ono8PL1THQ")

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

	//	v.c.SetToken("hvs.Meh6YP3vpWtVp2ono8PL1THQ")

	// Delete the role from Vault
	_, err := v.c.Logical().Delete(fmt.Sprintf("auth/kubernetes/role/%s", roleName))
	if err != nil {
		return fmt.Errorf("failed to delete Vault role: %v", err)
	}

	log.Printf("Vault role '%s' deleted successfully!", roleName)
	return nil
}
