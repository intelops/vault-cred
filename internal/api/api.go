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

func getCredentialMountPath(credentialType, credEntityName string) string {
	return fmt.Sprintf("%s/%s", credentialType, credEntityName)
}

func (v *VaultCredServ) GetCred(ctx context.Context, request *vaultcredpb.GetCredRequest) (*vaultcredpb.GetCredResponse, error) {
	fmt.Printf("here")
	vc, err := client.NewVaultClientForServiceAccount(ctx, v.conf, request.VaultRole, request.ServiceAccountToken)
	if err != nil {
		return nil, err
	}

	mountPath := getCredentialMountPath(request.CredentialType, request.CredEntityName)
	credentail, err := vc.GetCredential(ctx, mountPath, request.CredIdentifier)
	if err != nil {
		return nil, err
	}
	return &vaultcredpb.GetCredResponse{Credentail: credentail}, nil
}

func (v *VaultCredServ) PutCred(ctx context.Context, request *vaultcredpb.PutCredRequest) (*vaultcredpb.PutCredResponse, error) {
	vc, err := client.NewVaultClientForServiceAccount(ctx, v.conf, request.VaultRole, request.ServiceAccountToken)
	if err != nil {
		return nil, err
	}

	mountPath := getCredentialMountPath(request.CredentialType, request.CredEntityName)
	err = vc.PutCredential(ctx, mountPath, request.CredIdentifier, request.Credentail)
	if err != nil {
		return nil, err
	}
	return &vaultcredpb.PutCredResponse{}, nil
}

func (v *VaultCredServ) DeleteCred(ctx context.Context, request *vaultcredpb.DeleteCredRequest) (*vaultcredpb.DeleteCredResponse, error) {
	vc, err := client.NewVaultClientForServiceAccount(ctx, v.conf, request.VaultRole, request.ServiceAccountToken)
	if err != nil {
		return nil, err
	}

	mountPath := getCredentialMountPath(request.CredentialType, request.CredEntityName)
	err = vc.DeleteCredential(ctx, mountPath, request.CredIdentifier)
	if err != nil {
		return nil, err
	}

	return &vaultcredpb.DeleteCredResponse{}, nil
}
