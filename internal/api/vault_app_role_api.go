package api

import (
	"context"
	"fmt"

	"github.com/intelops/vault-cred/internal/client"
	"github.com/intelops/vault-cred/proto/pb/vaultcredpb"
)

func (v *VaultCredServ) CreateAppRoleToken(ctx context.Context, request *vaultcredpb.CreateAppRoleTokenRequest) (*vaultcredpb.CreateAppRoleTokenResponse, error) {
	v.log.Infof("app role token request for vault path %v with role %s", request.SecretPaths, request.AppRoleName)
	vc, err := client.NewVaultClientForTokenFromEnv(v.log, v.conf)
	if err != nil {
		return nil, err
	}

	err = vc.EnableAppRoleAuth()
	if err != nil {
		return nil, err
	}

	var policyData string
	for _, credPath := range request.SecretPaths {
		credPathPolicy := fmt.Sprintf(vaultPolicyReadPath, credPath)
		policyData = policyData + "\n" + credPathPolicy
	}

	policyName := request.AppRoleName + "-policy"
	err = vc.CreateOrUpdatePolicy(policyName, policyData)
	if err != nil {
		v.log.Errorf("error while creating Vault policy for app role %s", request.AppRoleName, err)
		return nil, err
	}

	err = vc.CreateOrUpdateAppRole(request.AppRoleName, []string{policyName})
	if err != nil {
		v.log.Errorf("error while creating Vault policy for app role %s", request.AppRoleName, err)
		return nil, err
	}

	token, err := vc.AuthenticateWithAppRole(request.AppRoleName)
	if err != nil {
		return nil, err
	}

	v.log.Infof("app role token generated for path %v with role %s", request.SecretPaths, request.AppRoleName)
	return &vaultcredpb.CreateAppRoleTokenResponse{Token: token}, nil
}

func (v *VaultCredServ) DeleteAppRole(ctx context.Context, request *vaultcredpb.DeleteAppRoleRequest) (*vaultcredpb.DeleteAppRoleResponse, error) {
	return nil, nil
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
