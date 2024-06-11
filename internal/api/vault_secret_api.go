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

	sort.SliceStable(secretPathProperties, func(i, j int) bool {
		if secretPathProperties[i].SecretKey != secretPathProperties[j].SecretKey {
			return secretPathProperties[i].SecretKey < secretPathProperties[j].SecretKey
		}
		if secretPathProperties[i].SecretPath != secretPathProperties[j].SecretPath {
			return secretPathProperties[i].SecretPath < secretPathProperties[j].SecretPath
		}
		return secretPathProperties[i].Property < secretPathProperties[j].Property
	})

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

	appRoleName := "kad-" + request.SecretName

	token, err := v.createAppRoleToken(context.Background(), appRoleName, secretPaths)
	if err != nil {
		v.log.Errorf("Error creating AppRole token: %v", err)
		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
	}

	k8sClient, err := client.NewK8SClient(v.log)
	if err != nil {
		v.log.Errorf("Failed to initialize k8s client: %v", err)
		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
	}

	cred := map[string][]byte{"token": []byte(token)}

	vaultTokenSecretName := "vault-token-" + request.SecretName

	err = k8sClient.CreateOrUpdateSecret(ctx, request.Namespace, vaultTokenSecretName, "Opaque", cred, nil)
	if err != nil {
		v.log.Errorf("Failed to create cluster vault token secret: %v", err)
		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
	}

	vaultAddressStr := fmt.Sprintf(vaultAddress, request.DomainName)

	secretStoreName := "ext-store-" + request.SecretName

	err = k8sClient.CreateOrUpdateSecretStore(ctx, secretStoreName, request.Namespace, vaultAddressStr, vaultTokenSecretName, "token")
	if err != nil {
		v.log.Errorf("Failed to create secret store: %v", err)
		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
	}

	externalSecretName := "ext-secret-" + request.SecretName

	v.log.Infof("Secret Paths Data: %v", secretPathsData)

	err = k8sClient.CreateOrUpdateExternalSecret(ctx, externalSecretName, request.Namespace, secretStoreName, request.SecretName, "", secretPathsData, propertiesData)
	if err != nil {
		v.log.Errorf("Failed to create vault external secret: %v", err)
		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
	}

	return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_OK}, nil
}
