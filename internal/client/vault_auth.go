package client

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
	"k8s.io/client-go/rest"
)

func (v *VaultClient) CheckAndEnableK8sAuth() error {
	mountPoints, err := v.c.Sys().ListAuth()
	if err != nil {
		return errors.WithMessage(err, "failed to get auth mount points")
	}

	v.log.Debugf("%d auth mountpoint found", len(mountPoints))
	authEnabled := false
	for _, mountPoint := range mountPoints {
		if mountPoint.Type == "kubernetes" {
			v.log.Debug("auth kubernetes mount point found")
			authEnabled = true
			break
		}
	}

	if authEnabled {
		v.log.Debug("auth kubernetes already enabled")
		return v.updateK8SAuthConfig()
	}

	options := api.MountInput{
		Type:        "kubernetes",
		Description: "Kubernetes authentication",
	}

	err = v.c.Sys().EnableAuthWithOptions("kubernetes", &options)
	if err != nil {
		return errors.WithMessage(err, "failed to enable auth kubernetes")
	}

	err = v.updateK8SAuthConfig()
	if err != nil {
		return err
	}

	v.log.Infof("auth kubernetes enabled")
	return nil
}

func (v *VaultClient) updateK8SAuthConfig() error {
	configData, err := v.loadK8SConfigData()
	if err != nil {
		return err
	}

	_, err = v.c.Logical().Write("/auth/kubernetes/config", configData)
	if err != nil {
		return errors.WithMessage(err, "failed to write to k8s auth config")
	}
	return nil
}

func (v *VaultClient) loadK8SConfigData() (map[string]interface{}, error) {
	k8s, err := NewK8SClient(v.log)
	if err != nil {
		return nil, errors.WithMessage(err, "error initializing k8s client")
	}

	config, err := k8s.GetClusterConfig()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get k8s api server endpoint")
	}

	err = rest.LoadTLSFiles(config)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to load k8s tls config")
	}

	return map[string]interface{}{
		"kubernetes_host":            "https://kubernetes.default.svc",
		"kubernetes_ca_cert":         string(config.CAData),
		"token_reviewer_jwt":         config.BearerToken,
		"kubernetes_skip_tls_verify": "true",
	}, nil
}

func (v *VaultClient) EnableAppRoleAuth() error {
	err := v.c.Sys().EnableAuth("approle", "approle", "")
	if err != nil && strings.Contains(err.Error(), "already in use") {
		return nil
	}
	return err
}

func (v *VaultClient) createOrUpdateAppRole(roleName string) (string, string, error) {
	roleIDResponse, err := v.c.Logical().Read("auth/approle/role/" + roleName + "/role-id")
	if err != nil {
		return "", "", err
	}

	secretIDResponse, err := v.c.Logical().Write("auth/approle/role/"+roleName+"/secret-id", nil)
	if err != nil {
		return "", "", err
	}

	roleID := roleIDResponse.Data["role_id"].(string)
	secretID := secretIDResponse.Data["secret_id"].(string)

	return roleID, secretID, nil
}

func (v *VaultClient) AuthenticateWithAppRole(roleName string) (string, error) {
	roleID, secretID, err := v.createOrUpdateAppRole(roleName)
	if err != nil {
		return "", err
	}

	data := map[string]interface{}{
		"role_id":   roleID,
		"secret_id": secretID,
	}

	secret, err := v.c.Logical().Write("auth/approle/login", data)
	if err != nil {
		return "", err
	}

	if secret == nil || secret.Auth == nil || secret.Auth.ClientToken == "" {
		return "", fmt.Errorf("authentication failed")
	}

	return secret.Auth.ClientToken, nil
}

func (v *VaultClient) ClusterEnableK8sAuth(clusterName, host, caCert, jwtToken string) error {
	options := api.MountInput{
		Type:        "kubernetes",
		Description: "Kubernetes authentication",
	}

	authPath := "k8s-" + clusterName
	err := v.c.Sys().EnableAuthWithOptions(authPath, &options)
	if err != nil {
		if !strings.Contains(err.Error(), "path is already in use") {
			return errors.WithMessage(err, "failed to enable auth kubernetes")
		}
	}

	configData := map[string]interface{}{
		"kubernetes_host":            host,
		"kubernetes_ca_cert":         caCert,
		"token_reviewer_jwt":         jwtToken,
		"kubernetes_skip_tls_verify": "true",
		"issuer":                     "https://kubernetes.default.svc.cluster.local",
	}

	configPath := fmt.Sprintf("/auth/%s/config/", authPath)
	_, err = v.c.Logical().Write(configPath, configData)
	if err != nil {
		return errors.WithMessage(err, "failed to write to k8s auth config")
	}

	v.log.Infof("cluster %s auth kubernetes enabled", clusterName)
	return nil
}
