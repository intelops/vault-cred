package api

import (
	"context"
	"fmt"

	"github.com/intelops/go-common/logging"
	"github.com/intelops/vault-cred/config"
	"github.com/intelops/vault-cred/internal/client"
	"github.com/intelops/vault-cred/proto/pb/vaultcredpb"
	"github.com/pkg/errors"
)

const (
	vaultPolicyReadPath = `path "secret/data/%s" {capabilities = ["read"]}`
)

type VaultCredServ struct {
	vaultcredpb.UnimplementedVaultCredServer
	conf config.VaultEnv
	log  logging.Logger
}

func NewVaultCredServ(log logging.Logger) (*VaultCredServ, error) {
	conf, err := config.GetVaultEnv()
	if err != nil {
		return nil, err
	}

	return &VaultCredServ{
		conf: conf,
		log:  log,
	}, nil
}

func CredentialMountPath() string {
	return "secret"
}

func PrepareCredentialSecretPath(credentialType, credEntityName, credIdentifier string) string {
	return fmt.Sprintf("%s/%s/%s", credentialType, credEntityName, credIdentifier)
}

func (v *VaultCredServ) GetCred(ctx context.Context, request *vaultcredpb.GetCredRequest) (*vaultcredpb.GetCredResponse, error) {
	vc, err := client.NewVaultClientForServiceAccount(ctx, v.log, v.conf)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to initiize vault client")
	}

	secretPath := PrepareCredentialSecretPath(request.CredentialType, request.CredEntityName, request.CredIdentifier)
	credentail, err := vc.GetCredential(ctx, CredentialMountPath(), secretPath)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get credential")
	}

	v.log.Infof("get credential request processed for %s", secretPath)
	return &vaultcredpb.GetCredResponse{Credential: credentail}, nil
}

func (v *VaultCredServ) PutCred(ctx context.Context, request *vaultcredpb.PutCredRequest) (*vaultcredpb.PutCredResponse, error) {
	vc, err := client.NewVaultClientForServiceAccount(ctx, v.log, v.conf)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to initiize vault client")
	}

	secretPath := PrepareCredentialSecretPath(request.CredentialType, request.CredEntityName, request.CredIdentifier)
	err = vc.PutCredential(ctx, CredentialMountPath(), secretPath, request.Credential)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to write credential")
	}

	v.log.Infof("write credential request processed for %s", secretPath)
	return &vaultcredpb.PutCredResponse{}, nil
}

func (v *VaultCredServ) DeleteCred(ctx context.Context, request *vaultcredpb.DeleteCredRequest) (*vaultcredpb.DeleteCredResponse, error) {
	vc, err := client.NewVaultClientForServiceAccount(ctx, v.log, v.conf)
	if err != nil {
		return nil, err
	}

	secretPath := PrepareCredentialSecretPath(request.CredentialType, request.CredEntityName, request.CredIdentifier)
	err = vc.DeleteCredential(ctx, CredentialMountPath(), secretPath)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to delete credential")
	}

	v.log.Infof("delete credential request processed for %s", secretPath)
	return &vaultcredpb.DeleteCredResponse{}, nil
}

func (v *VaultCredServ) GetAppRoleToken(ctx context.Context, request *vaultcredpb.GetAppRoleTokenRequest) (*vaultcredpb.GetAppRoleTokenResponse, error) {
	v.log.Infof("app role token request for vault path %s with role %s", request.CredentialPath, request.AppRoleName)
	vc, err := client.NewVaultClientForTokenFromEnv(v.log, v.conf)
	if err != nil {
		return nil, err
	}

	err = vc.EnableAppRoleAuth()
	if err != nil {
		return nil, err
	}

	policyData := fmt.Sprintf(vaultPolicyReadPath, request.CredentialPath)
	v.log.Infof("creating policy %s", policyData)
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

	v.log.Infof("app role token generated for path %s with role %s", request.CredentialPath, request.AppRoleName)
	return &vaultcredpb.GetAppRoleTokenResponse{Token: token}, nil
}

func (v *VaultCredServ) GetCredentialWithAppRoleToken(ctx context.Context, request *vaultcredpb.GetCredentialWithAppRoleTokenRequest) (*vaultcredpb.GetCredentialWithAppRoleTokenResponse, error) {
	vc, err := client.NewVaultClientForToken(v.log, v.conf, request.Token)
	if err != nil {
		return nil, err
	}

	credential, err := vc.GetCredential(ctx, CredentialMountPath(), request.CredentialPath)
	if err != nil {
		v.log.Error("app role get credential request failed for %s, %v", request.CredentialPath, err)
		return nil, err
	}
	v.log.Infof("app role get credential request processed for %s", request.CredentialPath)
	return &vaultcredpb.GetCredentialWithAppRoleTokenResponse{Credential: credential}, nil
}
