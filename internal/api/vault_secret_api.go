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

// 	secretPathsData := map[string][]string{}
// 	propertiesData := map[string][]string{}
// 	secretPaths := []string{}

// 	// Populate the secretPathsData and propertiesData maps
// 	for _, secretPathData := range request.SecretPathData {
// 		secretPathsData[secretPathData.SecretKey] = append(secretPathsData[secretPathData.SecretKey], secretPathData.SecretPath)
// 		secretPaths = append(secretPaths, secretPathData.SecretPath)
// 		if secretPathData.Property != "" {
// 			propertiesData[secretPathData.SecretKey] = append(propertiesData[secretPathData.SecretKey], secretPathData.Property)
// 		} else {
// 			propertiesData[secretPathData.SecretKey] = append(propertiesData[secretPathData.SecretKey], secretPathData.SecretKey)
// 		}
// 	}

// 	// Sort the paths and properties to ensure consistent ordering
// 	for key := range secretPathsData {
// 		sort.Strings(secretPathsData[key])
// 		sort.Strings(propertiesData[key])
// 	}

// 	// Log the sorted maps for debugging purposes
// 	v.log.Info("Sorted Secret Paths Data", secretPathsData)
// 	v.log.Info("Sorted Properties Data", propertiesData)

// 	appRoleName := kadAppRolePrefix + request.SecretName
// 	token, err := v.createAppRoleToken(context.Background(), appRoleName, secretPaths)
// 	if err != nil {
// 		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
// 	}

// 	k8sclient, err := client.NewK8SClient(v.log)
// 	if err != nil {
// 		v.log.Errorf("failed to initialize k8s client, %v", err)
// 		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
// 	}

// 	cred := map[string][]byte{"token": []byte(token)}
// 	vaultTokenSecretName := "vault-token-" + request.SecretName
// 	err = k8sclient.CreateOrUpdateSecret(ctx, request.Namespace, vaultTokenSecretName, v1.SecretTypeOpaque, cred, nil)
// 	if err != nil {
// 		v.log.Errorf("failed to create cluster vault token secret, %v", err)
// 		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
// 	}

// 	vaultAddressStr := fmt.Sprintf(vaultAddress, request.DomainName)
// 	secretStoreName := "ext-store-" + request.SecretName
// 	err = k8sclient.CreateOrUpdateSecretStore(ctx, secretStoreName, request.Namespace, vaultAddressStr, vaultTokenSecretName, "token")
// 	if err != nil {
// 		v.log.Errorf("failed to create secret store, %v", err)
// 		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
// 	}

// 	externalSecretName := "ext-secret-" + request.SecretName
// 	err = k8sclient.CreateOrUpdateExternalSecret(ctx, externalSecretName, request.Namespace, secretStoreName, request.SecretName, "", secretPathsData, propertiesData)
// 	if err != nil {
// 		v.log.Errorf("failed to create vault external secret, %v", err)
// 		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
// 	}

// 	return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_OK}, nil
// }

func (v *VaultCredServ) ConfigureVaultSecret(ctx context.Context, request *vaultcredpb.ConfigureVaultSecretRequest) (*vaultcredpb.ConfigureVaultSecretResponse, error) {
	v.log.Infof("Configure Vault Secret Request received for secret %s", request.SecretName)

	secretPathProperties := []SecretPathProperty{}

	for _, secretPathData := range request.SecretPathData {
		secretPathProperties = append(secretPathProperties, SecretPathProperty{
			SecretKey:  secretPathData.SecretKey,
			SecretPath: secretPathData.SecretPath,
			Property:   secretPathData.Property,
		})
	}

	secretPaths := []string{}
	secretPathsData := map[string][]string{}
	propertiesData := map[string][]string{}

	for _, spp := range secretPathProperties {
		secretPathsData[spp.SecretKey] = append(secretPathsData[spp.SecretKey], spp.SecretPath)
		secretPaths = append(secretPaths, spp.SecretPath)
		if spp.Property != "" {
			propertiesData[spp.SecretKey] = append(propertiesData[spp.SecretKey], spp.Property)
		} else {
			propertiesData[spp.SecretKey] = append(propertiesData[spp.SecretKey], spp.SecretKey)
		}
	}

	for key := range secretPathsData {
		sort.Strings(secretPathsData[key])
		sort.Strings(propertiesData[key])
	}

	appRoleName := "kad-" + request.SecretName
	token, err := v.createAppRoleToken(context.Background(), appRoleName, secretPaths)
	if err != nil {
		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
	}

	k8sClient, err := client.NewK8SClient(v.log)
	if err != nil {
		v.log.Errorf("failed to initialize k8s client, %v", err)
		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
	}

	cred := map[string][]byte{"token": []byte(token)}
	vaultTokenSecretName := "vault-token-" + request.SecretName
	err = k8sClient.CreateOrUpdateSecret(ctx, request.Namespace, vaultTokenSecretName, "Opaque", cred, nil)
	if err != nil {
		v.log.Errorf("failed to create cluster vault token secret, %v", err)
		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
	}

	vaultAddressStr := fmt.Sprintf("http://%s", request.DomainName)
	secretStoreName := "ext-store-" + request.SecretName
	err = k8sClient.CreateOrUpdateSecretStore(ctx, secretStoreName, request.Namespace, vaultAddressStr, vaultTokenSecretName, "token")
	if err != nil {
		v.log.Errorf("failed to create secret store, %v", err)
		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
	}

	externalSecretName := "ext-secret-" + request.SecretName
	v.log.Infof("Sorted Secret Paths Data: %v", secretPathsData)
	v.log.Infof("Properties Data: %v", propertiesData)
	err = k8sClient.CreateOrUpdateExternalSecret(ctx, externalSecretName, request.Namespace, secretStoreName, request.SecretName, "", secretPathsData, propertiesData)
	if err != nil {
		v.log.Errorf("failed to create vault external secret, %v", err)
		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
	}

	return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_OK}, nil
}
