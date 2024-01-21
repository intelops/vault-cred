package api

import (
	"context"
	"fmt"

	"github.com/intelops/vault-cred/proto/pb/vaultcredpb"
)

func (v *VaultCredServ) ConfigureClusterK8SAuth(ctx context.Context, request *vaultcredpb.ConfigureClusterK8SAuthRequest) (*vaultcredpb.ConfigureClusterK8SAuthResponse, error) {
	return nil, fmt.Errorf("not supported")
}

func (v *VaultCredServ) CreateK8SAuthRole(ctx context.Context, request *vaultcredpb.CreateK8SAuthRoleRequest) (*vaultcredpb.CreateK8SAuthRoleResponse, error) {
	return nil, nil
}

func (v *VaultCredServ) UpdateK8SAuthRole(ctx context.Context, request *vaultcredpb.UpdateK8SAuthRoleRequest) (*vaultcredpb.UpdateK8SAuthRoleResponse, error) {
	return nil, fmt.Errorf("not supported")
}

func (v *VaultCredServ) DeleteK8SAuthRole(ctx context.Context, request *vaultcredpb.DeleteK8SAuthRoleRequest) (*vaultcredpb.DeleteK8SAuthRoleResponse, error) {
	return nil, fmt.Errorf("not supported")
}
