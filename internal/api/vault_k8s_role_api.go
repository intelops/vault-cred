package api

import (
	"context"
	"fmt"

	"github.com/intelops/vault-cred/internal/client"
	"github.com/intelops/vault-cred/proto/pb/vaultcredpb"
)

func (v *VaultCredServ) AddClusterK8SAuth(ctx context.Context, request *vaultcredpb.AddClusterK8SAuthRequest) (*vaultcredpb.AddClusterK8SAuthResponse, error) {
	v.log.Infof("add k8s auth request for cluster %s", request.ClusterName)

	vc, err := client.NewVaultClientForTokenFromEnv(v.log, v.conf)
	if err != nil {
		v.log.Infof("error in getting vault client, %v", err)
		return &vaultcredpb.AddClusterK8SAuthResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
	}

	err = vc.ClusterEnableK8sAuth(request.ClusterName, request.Host, request.CaCert, request.JwtToken)
	if err != nil {
		v.log.Infof("error in adding k8s auth request for cluster %s, %v", request.ClusterName, err)
		return &vaultcredpb.AddClusterK8SAuthResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
	}
	v.log.Infof("k8s auth request for cluster %s added", request.ClusterName)
	return &vaultcredpb.AddClusterK8SAuthResponse{Status: vaultcredpb.StatusCode_OK}, nil
}

func (v *VaultCredServ) DeleteClusterK8SAuth(ctx context.Context, request *vaultcredpb.DeleteClusterK8SAuthRequest) (*vaultcredpb.DeleteClusterK8SAuthResponse, error) {
	return nil, fmt.Errorf("not supported")
}

func (v *VaultCredServ) CreateK8SAuthRole(ctx context.Context, request *vaultcredpb.CreateK8SAuthRoleRequest) (*vaultcredpb.CreateK8SAuthRoleResponse, error) {
	v.log.Infof("create k8s auth role %s for cluster %s", request.RoleName, request.ClusterName)

	vc, err := client.NewVaultClientForTokenFromEnv(v.log, v.conf)
	if err != nil {
		v.log.Infof("error in getting vault client, %v", err)
		return &vaultcredpb.CreateK8SAuthRoleResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
	}

	var policyData string
	for _, secretPolicy := range request.SecretPolicy {
		var credPathPolicy string
		if secretPolicy.Access == vaultcredpb.SecretAccess_READ {
			credPathPolicy = fmt.Sprintf(vaultPolicyReadPath, secretPolicy.SecretPath)
		} else if secretPolicy.Access == vaultcredpb.SecretAccess_WRITE {
			credPathPolicy = fmt.Sprintf(vaultPolicyWritePath, secretPolicy.SecretPath)
		} else {
			return &vaultcredpb.CreateK8SAuthRoleResponse{Status: vaultcredpb.StatusCode_INVALID_ARGUMENT}, fmt.Errorf("invalid security policy")
		}
		policyData = policyData + "\n" + credPathPolicy
	}

	policyName := "policy-" + request.ClusterName + "-" + request.RoleName
	err = vc.CreateOrUpdatePolicy(policyName, policyData)
	if err != nil {
		v.log.Infof("error in creating k8s auth policy %s, %v", policyName, err)
		return &vaultcredpb.CreateK8SAuthRoleResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
	}

	v.log.Infof("k8s auth policy %s created", policyName)

	roleName := "role-" + request.ClusterName + "-" + request.RoleName
	err = vc.CreateOrUpdateClusterRole(request.ClusterName, roleName, request.ServiceAccounts, request.Namespaces, []string{policyName})
	if err != nil {
		v.log.Infof("error in creating k8s auth policy %s, %v", policyName, err)
		return &vaultcredpb.CreateK8SAuthRoleResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
	}

	return &vaultcredpb.CreateK8SAuthRoleResponse{Status: vaultcredpb.StatusCode_OK}, nil
}

func (v *VaultCredServ) UpdateK8SAuthRole(ctx context.Context, request *vaultcredpb.UpdateK8SAuthRoleRequest) (*vaultcredpb.UpdateK8SAuthRoleResponse, error) {
	return nil, fmt.Errorf("not supported")
}

func (v *VaultCredServ) DeleteK8SAuthRole(ctx context.Context, request *vaultcredpb.DeleteK8SAuthRoleRequest) (*vaultcredpb.DeleteK8SAuthRoleResponse, error) {
	return nil, fmt.Errorf("not supported")
}
