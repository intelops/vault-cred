package api

import (
	"context"
	"fmt"
	"log"

	//"log"

	"github.com/intelops/vault-cred/internal/client"
	"github.com/intelops/vault-cred/proto/pb/vaultcredpb"
	v1 "k8s.io/api/core/v1"
)

var (
	kadAppRolePrefix = "capten-approle-"
	vaultAddress     = "http://vault.%s"
)

// func (v *VaultCredServ) ConfigureVaultSecret(ctx context.Context, request *vaultcredpb.ConfigureVaultSecretRequest) (*vaultcredpb.ConfigureVaultSecretResponse, error) {
// 	v.log.Infof("Configure Vault Secret Request recieved for secret ", request.SecretName)

// 	secretPathsData := map[string]string{}
// 	secretPaths := []string{}
// 	for _, secretPathData := range request.SecretPathData {
// 		secretPathsData[secretPathData.SecretKey] = secretPathData.SecretPath
// 		secretPaths = append(secretPaths, secretPathData.SecretPath)
// 	}

// 	appRoleName := kadAppRolePrefix + request.SecretName
// 	token, err := v.createAppRoleToken(context.Background(), appRoleName, secretPaths)
// 	if err != nil {
// 		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
// 	}

// 	k8sclient, err := client.NewK8SClient(v.log)
// 	if err != nil {
// 		v.log.Errorf("failed to initalize k8s client, %v", err)
// 		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
// 	}

// 	cred := map[string][]byte{"token": []byte(token)}
// 	vaultTokenSecretName := "vault-token-" + request.SecretName
// 	err = k8sclient.CreateOrUpdateSecret(ctx, request.Namespace, vaultTokenSecretName, v1.SecretTypeOpaque, cred, nil)
// 	if err != nil {
// 		v.log.Errorf("failed to create cluter vault token secret, %v", err)
// 		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
// 	}

// 	vaultAddressStr := fmt.Sprintf(vaultAddress, request.DomainName)
// 	secretStoreName := "ext-store-" + request.SecretName
// 	err = k8sclient.CreateOrUpdateSecretStore(ctx, secretStoreName, request.Namespace, vaultAddressStr, vaultTokenSecretName, "token")
// 	if err != nil {
// 		v.log.Errorf("failed to create cluter vault token secret, %v", err)
// 		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
// 	}
// 	v.log.Infof("created secret store %s/%s", request.Namespace, secretStoreName)

// 	externalSecretName := "ext-secret-" + request.SecretName
// 	err = k8sclient.CreateOrUpdateExternalSecret(ctx, externalSecretName, request.Namespace, secretStoreName,
// 		request.SecretName, "", secretPathsData)
// 	if err != nil {
// 		v.log.Errorf("failed to create vault external secret, %v", err)
// 		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
// 	}
// 	v.log.Infof("created external secret %s/%s", request.Namespace, externalSecretName)
// 	return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_OK}, nil
// }

// func (v *VaultCredServ) ConfigureVaultSecret(ctx context.Context, request *vaultcredpb.ConfigureVaultSecretRequest) (*vaultcredpb.ConfigureVaultSecretResponse, error) {
// 	v.log.Infof("Configure Vault Secret Request received for secret %s", request.SecretName)

// 	secretPathsData := map[string]string{}
// 	secretPaths := []string{}
// 	properties := map[string]string{}
// 	for _, secretPathData := range request.SecretPathData {
// 		secretPathsData[secretPathData.SecretKey] = secretPathData.SecretPath
// 		secretPaths = append(secretPaths, secretPathData.SecretPath)
// 		properties[secretPathData.SecretKey] = secretPathData.Property
// 	}
// 	log.Println("Prop", properties)
// 	log.Println("secretpath data in configuring vault secret", secretPathsData)
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
// 		v.log.Errorf("failed to create cluster vault token secret, %v", err)
// 		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err

// 	}
// 	v.log.Infof("created secret store %s/%s", request.Namespace, secretStoreName)

// 	externalSecretName := "ext-secret-" + request.SecretName
// 	err = k8sclient.CreateOrUpdateExternalSecret(ctx, externalSecretName, request.Namespace, secretStoreName,
// 		request.SecretName, "", secretPathsData, properties)
// 	if err != nil {
// 		v.log.Errorf("failed to create vault external secret, %v", err)
// 		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err

// 	}
// 	v.log.Infof("created external secret %s/%s", request.Namespace, externalSecretName)
// 	return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_OK}, nil
// }

// func (v *VaultCredServ) ConfigureVaultSecret(ctx context.Context, request *vaultcredpb.ConfigureVaultSecretRequest) (*vaultcredpb.ConfigureVaultSecretResponse, error) {
// 	v.log.Infof("Configure Vault Secret Request received for secret %s", request.SecretName)

// 	secretPathsData := make(map[string][]string)
// 	propertiesData := make(map[string][]string)
// 	secretPaths := []string{}

// 	log.Println("Request path data", request.SecretPathData)

// 	for _, secretPathData := range request.SecretPathData {
// 		secretPathsData[secretPathData.SecretKey] = append(secretPathsData[secretPathData.SecretKey], secretPathData.SecretPath)
// 		secretPaths = append(secretPaths, secretPathData.SecretPath)
// 		if secretPathData.Property != "" {
// 			propertiesData[secretPathData.SecretKey] = append(propertiesData[secretPathData.SecretKey], secretPathData.Property)
// 		} else {
// 			propertiesData[secretPathData.SecretKey] = append(propertiesData[secretPathData.SecretKey], secretPathData.SecretKey) // default to secretKey if property is not provided
// 		}
// 	}

// 	log.Println("Secret Paths data while configuring", secretPathsData)
// 	log.Println("Properties while configuring", propertiesData)

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
// 		v.log.Errorf("failed to create cluster vault secret store, %v", err)
// 		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
// 	}
// 	v.log.Infof("created secret store %s/%s", request.Namespace, secretStoreName)

// 	externalSecretName := "ext-secret-" + request.SecretName
// 	err = k8sclient.CreateOrUpdateExternalSecret(ctx, externalSecretName, request.Namespace, secretStoreName, request.SecretName, "", secretPathsData, propertiesData)
// 	if err != nil {
// 		v.log.Errorf("failed to create vault external secret, %v", err)
// 		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
// 	}
// 	v.log.Infof("created external secret %s/%s", request.Namespace, externalSecretName)
// 	return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_OK}, nil
// }

func (v *VaultCredServ) ConfigureVaultSecret(ctx context.Context, request *vaultcredpb.ConfigureVaultSecretRequest) (*vaultcredpb.ConfigureVaultSecretResponse, error) {
	v.log.Infof("Configure Vault Secret Request received for secret %s", request.SecretName)

	secretPathsData := map[string][]string{}
	propertiesData := map[string][]string{}
	secretPaths := []string{}

	for _, secretPathData := range request.SecretPathData {
		secretPathsData[secretPathData.SecretKey] = append(secretPathsData[secretPathData.SecretKey], secretPathData.SecretPath)
		secretPaths = append(secretPaths, secretPathData.SecretPath)
		if secretPathData.Property != "" {
			propertiesData[secretPathData.SecretKey] = append(propertiesData[secretPathData.SecretKey], secretPathData.Property)
		} else {
			propertiesData[secretPathData.SecretKey] = append(propertiesData[secretPathData.SecretKey], secretPathData.SecretKey)
		}
	}

	appRoleName := kadAppRolePrefix + request.SecretName
	token, err := v.createAppRoleToken(context.Background(), appRoleName, secretPaths)
	if err != nil {
		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
	}

	k8sclient, err := client.NewK8SClient(v.log)
	if err != nil {
		v.log.Errorf("failed to initialize k8s client, %v", err)
		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
	}

	cred := map[string][]byte{"token": []byte(token)}
	vaultTokenSecretName := "vault-token-" + request.SecretName
	err = k8sclient.CreateOrUpdateSecret(ctx, request.Namespace, vaultTokenSecretName, v1.SecretTypeOpaque, cred, nil)
	if err != nil {
		v.log.Errorf("failed to create cluster vault token secret, %v", err)
		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
	}

	vaultAddressStr := fmt.Sprintf(vaultAddress, request.DomainName)
	secretStoreName := "ext-store-" + request.SecretName
	err = k8sclient.CreateOrUpdateSecretStore(ctx, secretStoreName, request.Namespace, vaultAddressStr, vaultTokenSecretName, "token")
	if err != nil {
		v.log.Errorf("failed to create secret store, %v", err)
		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
	}

	externalSecretName := "ext-secret-" + request.SecretName
	log.Println("Secret Paths Data", secretPathsData)
	log.Println("Properties data", propertiesData)
	err = k8sclient.CreateOrUpdateExternalSecret(ctx, externalSecretName, request.Namespace, secretStoreName, request.SecretName, "", secretPathsData, propertiesData)
	if err != nil {
		v.log.Errorf("failed to create vault external secret, %v", err)
		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
	}

	return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_OK}, nil
}
