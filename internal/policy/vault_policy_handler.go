package policy

import (
	"context"
	"strings"

	"github.com/intelops/go-common/logging"
	"github.com/intelops/vault-cred/config"
	"github.com/intelops/vault-cred/internal/client"
	"github.com/pkg/errors"
)

type VaultPolicyHandler struct {
	log  logging.Logger
	conf config.VaultEnv
}

func NewVaultPolicyHandler(log logging.Logger, conf config.VaultEnv) *VaultPolicyHandler {
	return &VaultPolicyHandler{log: log, conf: conf}
}

func (p *VaultPolicyHandler) getVaultConfigMaps(ctx context.Context, prefix string) (map[string]map[string]string, error) {
	k8s, err := client.NewK8SClient(p.log)
	if err != nil {
		return nil, err
	}

	allConfigMapData, err := k8s.GetConfigMapsHasPrefix(ctx, prefix)
	if err != nil {
		return nil, errors.WithMessagef(err, "error while getting vault policy configmaps")
	}
	return allConfigMapData, nil
}

func (p *VaultPolicyHandler) UpdateVaultPolicies(ctx context.Context) error {
	vc, err := client.NewVaultClientForVaultToken(p.log, p.conf)
	if err != nil {
		return err
	}

	allConfigMapData, err := p.getVaultConfigMaps(ctx, "vault-policy-")
	if err != nil {
		return errors.WithMessagef(err, "error while getting vault policy configmaps")
	}

	for _, cmData := range allConfigMapData {
		policyName := cmData["policyName"]
		policyData := cmData["policyData"]
		err = vc.CreateOrUpdatePolicy(policyName, policyData)
		if err != nil {
			return errors.WithMessagef(err, "error while creating vault policy %s, %v", policyName, cmData)
		}
	}
	return nil
}

func (p *VaultPolicyHandler) UpdateVaultRoles(ctx context.Context) error {
	vc, err := client.NewVaultClientForVaultToken(p.log, p.conf)
	if err != nil {
		return err
	}

	allConfigMapData, err := p.getVaultConfigMaps(ctx, "vault-role-")
	if err != nil {
		return errors.WithMessagef(err, "error while getting vault role configmaps")
	}

	for _, cmData := range allConfigMapData {
		roleName := cmData["policyNames"]
		policyNames := cmData["policyNames"]
		servieAccounts := cmData["servieAccounts"]
		servieAccountNameSpaces := cmData["servieAccountNameSpaces"]
		err = vc.CreateOrUpdateRole(roleName,
			strings.Split(policyNames, ","),
			strings.Split(servieAccounts, ","),
			strings.Split(servieAccountNameSpaces, ","))
		if err != nil {
			return errors.WithMessagef(err, "error while creating vault role %s, %v", roleName, cmData)
		}
	}
	return nil
}
