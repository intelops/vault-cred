package job

import (
	"context"

	"github.com/intelops/go-common/logging"
	"github.com/intelops/vault-cred/config"
	"github.com/intelops/vault-cred/internal/policy"
)

type VaultPolicyWatcher struct {
	log       logging.Logger
	handler   *policy.VaultPolicyHandler
	frequency string
}

func NewVaultPolicyWatcher(log logging.Logger, frequency string) (*VaultPolicyWatcher, error) {
	conf, err := config.GetVaultEnv()
	if err != nil {
		return nil, err
	}
	return &VaultPolicyWatcher{
		log:       log,
		frequency: frequency,
		handler:   policy.NewVaultPolicyHandler(log, conf),
	}, nil
}

func (v *VaultPolicyWatcher) CronSpec() string {
	return v.frequency
}

func (v *VaultPolicyWatcher) Run() {
	v.log.Info("started vault policy watcher")
	ctx := context.Background()
	if err := v.handler.UpdateVaultPolicies(ctx); err != nil {
		v.log.Errorf("failed to update vault policies, %v", err)
	}

	if err := v.handler.UpdateVaultRoles(ctx); err != nil {
		v.log.Errorf("failed to update roles, %v", err)
	}
}
