package vault

import (
	"github.com/hashicorp/vault/api"
	"github.com/intelops/vault-cred/internal/k8s"
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
	k8s, err := k8s.NewK8SClient(v.log)
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
