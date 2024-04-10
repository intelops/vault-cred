package api

import (
	"context"
	"fmt"

	"github.com/intelops/vault-cred/internal/client"
	"github.com/intelops/vault-cred/proto/pb/vaultcredpb"
)

func (v *VaultCredServ) CreateAppRoleToken(ctx context.Context,
	request *vaultcredpb.CreateAppRoleTokenRequest) (*vaultcredpb.CreateAppRoleTokenResponse, error) {
	token, err := v.createAppRoleToken(ctx, request.AppRoleName, request.SecretPaths)
	return &vaultcredpb.CreateAppRoleTokenResponse{Token: token}, err
}

func (v *VaultCredServ) createAppRoleToken(ctx context.Context, appRoleName string, secretPaths []string) (string, error) {
	v.log.Infof("app role token request for vault path %v with role %s", secretPaths, appRoleName)
	vc, err := client.NewVaultClientForTokenFromEnv(v.log, v.conf)
	if err != nil {
		return "", err
	}

	err = vc.EnableAppRoleAuth()
	if err != nil {
		return "", err
	}

	var policyData string
	for _, credPath := range secretPaths {
		credPathPolicy := fmt.Sprintf(vaultPolicyReadPath, credPath)
		policyData = policyData + "\n" + credPathPolicy
	}

	policyName := appRoleName + "-policy"
	err = vc.CreateOrUpdatePolicy(policyName, policyData)
	if err != nil {
		v.log.Errorf("error while creating Vault policy for app role %s", appRoleName, err)
		return "", err
	}

	err = vc.CreateOrUpdateAppRole(appRoleName, []string{policyName})
	if err != nil {
		v.log.Errorf("error while creating Vault policy for app role %s", appRoleName, err)
		return "", err
	}

	token, err := vc.AuthenticateWithAppRole(appRoleName)
	if err != nil {
		return "", err
	}

	v.log.Infof("app role token generated for path %v with role %s", secretPaths, appRoleName)
	return token, nil
}

func (v *VaultCredServ) DeleteAppRole(ctx context.Context, request *vaultcredpb.DeleteAppRoleRequest) (*vaultcredpb.DeleteAppRoleResponse, error) {
	v.log.Infof("app role delete request for vault role %s", request.RoleName)
	vc, err := client.NewVaultClientForTokenFromEnv(v.log, v.conf)
	if err != nil {
		return nil, err
	}

	err = vc.EnableAppRoleAuth()
	if err != nil {
		return nil, err
	}

	err = vc.DeleteRole(request.RoleName)
	if err != nil {
		return nil, err
	}

	policyName := request.RoleName + "-policy"
	err = vc.DeletePolicy(policyName)
	if err != nil {
		return nil, err
	}

	return &vaultcredpb.DeleteAppRoleResponse{
		Status:        vaultcredpb.StatusCode_OK,
		StatusMessage: fmt.Sprintf("app role %s deleted", request.RoleName),
	}, nil
}

func (v *VaultCredServ) GetCredentialWithAppRoleToken(ctx context.Context, request *vaultcredpb.GetCredentialWithAppRoleTokenRequest) (*vaultcredpb.GetCredentialWithAppRoleTokenResponse, error) {
	vc, err := client.NewVaultClientForToken(v.log, v.conf, request.Token)
	if err != nil {
		return nil, err
	}

	credential, err := vc.GetCredential(ctx, CredentialMountPath(), request.SecretPath)
	if err != nil {
		v.log.Error("app role get credential request failed for %s, %v", request.SecretPath, err)
		return nil, err
	}
	v.log.Infof("app role get credential request processed for %s", request.SecretPath)
	return &vaultcredpb.GetCredentialWithAppRoleTokenResponse{Credential: credential}, nil
}
