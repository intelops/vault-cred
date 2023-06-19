package client

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
	"k8s.io/client-go/rest"
)

type VaultPolicyData struct {
	VaultRoleName           string              `json:"vaultRoleName"`
	PolicyName              string              `json:"policyName"`
	ServiceAccount          string              `json:"serviceAccount"`
	ServiceAccountNameSpace string              `json:"serviceAccountNameSpace"`
	CredentialAccessList    map[string][]string `json:"credentialAccessList"`
}

func (v *VaultClient) CreateOrUpdatePolicy(policyName, rules string) error {
	err := v.c.Sys().PutPolicy(policyName, rules)
	if err != nil {
		return err
	}

	v.log.Infof("Updated policy %s", policyName)
	return nil
}

func (v *VaultClient) DeletePolicy(policyName string) error {
	err := v.c.Sys().DeletePolicy(policyName)
	if err != nil {
		return err
	}
	v.log.Infof("Deleted policy %s", policyName)
	return nil
}

func (v *VaultClient) CreateOrUpdateRole(roleName string, serviceAccounts, namespaces, policies []string) error {
	roleData := make(map[string]interface{})

	sa := strings.Join(serviceAccounts, ",")
	ns := strings.Join(namespaces, ",")
	roleData["bound_service_account_names"] = sa
	roleData["bound_service_account_namespaces"] = ns
	roleData["policies"] = policies
	roleData["max_ttl"] = 1800000

	path := fmt.Sprintf("/auth/kubernetes/role/%s", roleName)
	_, err := v.c.Logical().Write(path, roleData)
	if err != nil {
		return err
	}
	v.log.Infof("Updated role %s", roleName)
	return nil
}

func (v *VaultClient) DeleteRole(roleName string) error {
	path := fmt.Sprintf("/auth/kubernetes/role/%s", roleName)
	_, err := v.c.Logical().Delete(path)
	if err != nil {
		return err
	}
	v.log.Infof("Deleted role %s", roleName)
	return nil
}

func (v *VaultClient) ListPolicies() ([]string, error) {
	return v.c.Sys().ListPolicies()
}

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
		return nil
	}

	k8s, err := NewK8SClient(v.log)
	if err != nil {
		return errors.WithMessage(err, "error initializing k8s client")
	}

	config, err := k8s.GetClusterConfig()
	if err != nil {
		return errors.WithMessage(err, "failed to get k8s api server endpoint")
	}

	err = rest.LoadTLSFiles(config)
	if err != nil {
		return errors.WithMessage(err, "failed to load k8s tls config")
	}
	base64CertData := base64.StdEncoding.EncodeToString([]byte(config.CAData))

	options := api.MountInput{
		Type:        "kubernetes",
		Description: "Kubernetes authentication",
		Options: map[string]string{
			"kubernetes_host":            config.Host,
			"kubernetes_ca_cert":         base64CertData,
			"token_reviewer_jwt":         config.BearerToken,
			"kubernetes_skip_tls_verify": "false",
		},
	}

	err = v.c.Sys().EnableAuthWithOptions("kubernetes", &options)
	if err != nil {
		return errors.WithMessage(err, "failed to enable auth kubernetes")
	}
	v.log.Infof("auth kubernetes enabled")
	return nil
}
