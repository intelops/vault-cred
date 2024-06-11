package api

import (
	"context"
	"fmt"

	"sort"

	"github.com/intelops/vault-cred/internal/client"
	"github.com/intelops/vault-cred/proto/pb/vaultcredpb"
)

var (
	kadAppRolePrefix = "capten-approle-"
	vaultAddress     = "http://vault.%s"
)

type SecretPathProperty struct {
	SecretKey  string
	SecretPath string
	Property   string
}

// func (v *VaultCredServ) ConfigureVaultSecret(ctx context.Context, request *vaultcredpb.ConfigureVaultSecretRequest) (*vaultcredpb.ConfigureVaultSecretResponse, error) {
// 	v.log.Infof("Configure Vault Secret Request received for secret %s", request.SecretName)

// 	secretPathProperties := []SecretPathProperty{}

// 	for _, secretPathData := range request.SecretPathData {
// 		secretPathProperties = append(secretPathProperties, SecretPathProperty{
// 			SecretKey:  secretPathData.SecretKey,
// 			SecretPath: secretPathData.SecretPath,
// 			Property:   secretPathData.Property,
// 		})
// 	}

// 	secretPaths := []string{}
// 	secretPathsData := map[string][]string{}
// 	propertiesData := map[string][]string{}

// 	for _, spp := range secretPathProperties {
// 		secretPathsData[spp.SecretKey] = append(secretPathsData[spp.SecretKey], spp.SecretPath)
// 		secretPaths = append(secretPaths, spp.SecretPath)
// 		if spp.Property != "" {
// 			propertiesData[spp.SecretKey] = append(propertiesData[spp.SecretKey], spp.Property)
// 		} else {
// 			propertiesData[spp.SecretKey] = append(propertiesData[spp.SecretKey], spp.SecretKey)
// 		}
// 	}

// 	for key := range secretPathsData {
// 		sort.Strings(secretPathsData[key])
// 		sort.Strings(propertiesData[key])
// 	}

// 	appRoleName := "kad-" + request.SecretName
// 	token, err := v.createAppRoleToken(context.Background(), appRoleName, secretPaths)
// 	if err != nil {
// 		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
// 	}

// 	k8sClient, err := client.NewK8SClient(v.log)
// 	if err != nil {
// 		v.log.Errorf("failed to initialize k8s client, %v", err)
// 		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
// 	}

// 	cred := map[string][]byte{"token": []byte(token)}
// 	vaultTokenSecretName := "vault-token-" + request.SecretName
// 	err = k8sClient.CreateOrUpdateSecret(ctx, request.Namespace, vaultTokenSecretName, "Opaque", cred, nil)
// 	if err != nil {
// 		v.log.Errorf("failed to create cluster vault token secret, %v", err)
// 		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
// 	}

// 	// vaultAddressStr := fmt.Sprintf("http://%s", request.DomainName)
// 	vaultAddressStr := fmt.Sprintf(vaultAddress, request.DomainName)
// 	log.Println("Vault Address string", vaultAddressStr)
// 	secretStoreName := "ext-store-" + request.SecretName
// 	err = k8sClient.CreateOrUpdateSecretStore(ctx, secretStoreName, request.Namespace, vaultAddressStr, vaultTokenSecretName, "token")
// 	if err != nil {
// 		v.log.Errorf("failed to create secret store, %v", err)
// 		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
// 	}

// 	externalSecretName := "ext-secret-" + request.SecretName
// 	v.log.Infof("Sorted Secret Paths Data: %v", secretPathsData)
// 	v.log.Infof("Properties Data: %v", propertiesData)
// 	err = k8sClient.CreateOrUpdateExternalSecret(ctx, externalSecretName, request.Namespace, secretStoreName, request.SecretName, "", secretPathsData, propertiesData)
// 	if err != nil {
// 		v.log.Errorf("failed to create vault external secret, %v", err)
// 		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
// 	}

// 	return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_OK}, nil
// }

func (v *VaultCredServ) ConfigureVaultSecret(ctx context.Context, request *vaultcredpb.ConfigureVaultSecretRequest) (*vaultcredpb.ConfigureVaultSecretResponse, error) {
	v.log.Infof("Configure Vault Secret Request received for secret %s", request.SecretName)

	// Initialize a slice to hold SecretPathProperty structs
	secretPathProperties := []SecretPathProperty{}

	// Populate the secretPathProperties slice with data from the request
	for _, secretPathData := range request.SecretPathData {
		secretPathProperties = append(secretPathProperties, SecretPathProperty{
			SecretKey:  secretPathData.SecretKey,
			SecretPath: secretPathData.SecretPath,
			Property:   secretPathData.Property,
		})
	}

	// Log initial secret path properties
	v.log.Infof("Initial Secret Path Properties: %v", secretPathProperties)

	// Sort the secretPathProperties slice by SecretKey, then by SecretPath, and then by Property
	sort.SliceStable(secretPathProperties, func(i, j int) bool {
		if secretPathProperties[i].SecretKey != secretPathProperties[j].SecretKey {
			return secretPathProperties[i].SecretKey < secretPathProperties[j].SecretKey
		}
		if secretPathProperties[i].SecretPath != secretPathProperties[j].SecretPath {
			return secretPathProperties[i].SecretPath < secretPathProperties[j].SecretPath
		}
		return secretPathProperties[i].Property < secretPathProperties[j].Property
	})

	// Log sorted secret path properties
	v.log.Infof("Sorted Secret Path Properties: %v", secretPathProperties)

	// Initialize slices and maps to hold secret paths and properties
	secretPaths := []string{}
	secretPathsData := map[string][]string{}
	propertiesData := map[string][]string{}

	// Populate the maps and slice with data from sorted secretPathProperties
	for _, spp := range secretPathProperties {
		secretPathsData[spp.SecretKey] = append(secretPathsData[spp.SecretKey], spp.SecretPath)
		secretPaths = append(secretPaths, spp.SecretPath)
		if spp.Property != "" {
			propertiesData[spp.SecretKey] = append(propertiesData[spp.SecretKey], spp.Property)
		} else {
			propertiesData[spp.SecretKey] = append(propertiesData[spp.SecretKey], spp.SecretKey)
		}
	}

	// Log secret paths data and properties data after sorting and population
	v.log.Infof("Secret Paths Data after sorting and population: %v", secretPathsData)
	v.log.Infof("Properties Data after sorting and population: %v", propertiesData)

	// Generate an AppRole name using the secret name
	appRoleName := "kad-" + request.SecretName

	// Create an AppRole token using the secret paths
	token, err := v.createAppRoleToken(context.Background(), appRoleName, secretPaths)
	if err != nil {
		v.log.Errorf("Error creating AppRole token: %v", err)
		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
	}

	// Initialize a Kubernetes client
	k8sClient, err := client.NewK8SClient(v.log)
	if err != nil {
		v.log.Errorf("Failed to initialize k8s client: %v", err)
		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
	}

	// Create a map to hold the token data
	cred := map[string][]byte{"token": []byte(token)}

	// Generate a name for the Kubernetes secret to store the Vault token
	vaultTokenSecretName := "vault-token-" + request.SecretName

	// Create or update the Kubernetes secret with the Vault token
	err = k8sClient.CreateOrUpdateSecret(ctx, request.Namespace, vaultTokenSecretName, "Opaque", cred, nil)
	if err != nil {
		v.log.Errorf("Failed to create cluster vault token secret: %v", err)
		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
	}

	// Format the Vault address string using the domain name from the request
	vaultAddressStr := fmt.Sprintf(vaultAddress, request.DomainName)
	v.log.Infof("Vault Address string: %s", vaultAddressStr)

	// Generate a name for the SecretStore
	secretStoreName := "ext-store-" + request.SecretName

	// Create or update the SecretStore in Kubernetes
	err = k8sClient.CreateOrUpdateSecretStore(ctx, secretStoreName, request.Namespace, vaultAddressStr, vaultTokenSecretName, "token")
	if err != nil {
		v.log.Errorf("Failed to create secret store: %v", err)
		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
	}

	// Generate a name for the ExternalSecret
	externalSecretName := "ext-secret-" + request.SecretName

	// Log the sorted secret paths and properties for debugging
	v.log.Infof("Final Secret Paths Data: %v", secretPathsData)
	v.log.Infof("Final Properties Data: %v", propertiesData)

	// Create or update the ExternalSecret in Kubernetes using the sorted data
	err = k8sClient.CreateOrUpdateExternalSecret(ctx, externalSecretName, request.Namespace, secretStoreName, request.SecretName, "", secretPathsData, propertiesData)
	if err != nil {
		v.log.Errorf("Failed to create vault external secret: %v", err)
		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
	}

	// Return a successful response
	return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_OK}, nil
}
