package api

import (
	"context"
	"fmt"

	"github.com/intelops/vault-cred/config"
	"github.com/intelops/vault-cred/internal/client"
	"github.com/intelops/vault-cred/proto/pb/vaultcredpb"
)

type VaultCredServ struct {
	vaultcredpb.UnimplementedVaultCredServer
	conf config.VaultEnv
}

func NewVaultCredServ() (*VaultCredServ, error) {
	conf, err := config.GetVaultEnv()
	if err != nil {
		return nil, err
	}

	return &VaultCredServ{
		conf: conf,
	}, nil
}

func getCredentialPath(credentialType, credEntityName, credIdentifier string) string {
	return fmt.Sprintf("%s/%s/%s", credentialType, credEntityName, credIdentifier)
}

func (v *VaultCredServ) GetCredRequest(ctx context.Context, request *vaultcredpb.GetCredRequest) (*vaultcredpb.GetCredResponse, error) {
	vc, err := client.NewVaultClientForServiceAccount(ctx, v.conf, request.VaultRole, request.ServiceAccountToken)
	if err != nil {
		return nil, err
	}

	secretPath := getCredentialPath(request.CredentialType, request.CredEntityName, request.CredIdentifier)
	credentail, err := vc.GetCredential(ctx, secretPath)
	if err != nil {
		return nil, err
	}
	return &vaultcredpb.GetCredResponse{Credentail: credentail}, nil
}

func (v *VaultCredServ) PutCredRequest(ctx context.Context, request *vaultcredpb.PutCredRequest) (*vaultcredpb.PutCredResponse, error) {
	vc, err := client.NewVaultClientForServiceAccount(ctx, v.conf, request.VaultRole, request.ServiceAccountToken)
	if err != nil {
		return nil, err
	}

	secretPath := getCredentialPath(request.CredentialType, request.CredEntityName, request.CredIdentifier)
	err = vc.PutCredential(ctx, secretPath, request.Credentail)
	if err != nil {
		return nil, err
	}
	return &vaultcredpb.PutCredResponse{}, nil
}

func (v *VaultCredServ) DeleteCredRequest(ctx context.Context, request *vaultcredpb.PutCredRequest) (*vaultcredpb.DeleteCredResponse, error) {
	vc, err := client.NewVaultClientForServiceAccount(ctx, v.conf, request.VaultRole, request.ServiceAccountToken)
	if err != nil {
		return nil, err
	}

	secretPath := getCredentialPath(request.CredentialType, request.CredEntityName, request.CredIdentifier)
	err = vc.DeleteCredential(ctx, secretPath)
	if err != nil {
		return nil, err
	}

	return &vaultcredpb.DeleteCredResponse{}, nil
}
