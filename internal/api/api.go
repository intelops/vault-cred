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

func credentialMountPath() string {
	return "secret"
}

func prepareCredentialSecretPath(credentialType, credEntityName, credIdentifier string) string {
	return fmt.Sprintf("%s/%s/%s", credentialType, credEntityName, credIdentifier)
}

func (v *VaultCredServ) GetCred(ctx context.Context, request *vaultcredpb.GetCredRequest) (*vaultcredpb.GetCredResponse, error) {
	vc, err := client.NewVaultClientForServiceAccount(ctx, v.conf, request.VaultRole, request.ServiceAccountToken)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to initiize vault client")
	}

	secretPath := prepareCredentialSecretPath(request.CredentialType, request.CredEntityName, request.CredIdentifier)
	credentail, err := vc.GetCredential(ctx, credentialMountPath(), secretPath)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get credential")
	}

	v.log.Infof("get credential request processed for %s", secretPath)
	return &vaultcredpb.GetCredResponse{Credentail: credentail}, nil
}

func (v *VaultCredServ) PutCred(ctx context.Context, request *vaultcredpb.PutCredRequest) (*vaultcredpb.PutCredResponse, error) {
	vc, err := client.NewVaultClientForServiceAccount(ctx, v.conf, request.VaultRole, request.ServiceAccountToken)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to initiize vault client")
	}

	secretPath := prepareCredentialSecretPath(request.CredentialType, request.CredEntityName, request.CredIdentifier)
	err = vc.PutCredential(ctx, credentialMountPath(), secretPath, request.Credentail)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to write credential")
	}

	v.log.Infof("write credential request processed for %s", secretPath)
	return &vaultcredpb.PutCredResponse{}, nil
}

func (v *VaultCredServ) DeleteCred(ctx context.Context, request *vaultcredpb.DeleteCredRequest) (*vaultcredpb.DeleteCredResponse, error) {
	vc, err := client.NewVaultClientForServiceAccount(ctx, v.conf, request.VaultRole, request.ServiceAccountToken)
	if err != nil {
		return nil, err
	}

	secretPath := prepareCredentialSecretPath(request.CredentialType, request.CredEntityName, request.CredIdentifier)
	err = vc.DeleteCredential(ctx, credentialMountPath(), secretPath)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to delete credential")
	}

	v.log.Infof("delete credential request processed for %s", secretPath)
	return &vaultcredpb.DeleteCredResponse{}, nil
}
