package job

import (
	"context"

	"github.com/intelops/go-common/logging"
	"github.com/intelops/vault-cred/config"
	"github.com/intelops/vault-cred/internal/client"
	"github.com/intelops/vault-cred/internal/policy"
)

type VaultPolicyWatcher struct {
	log       logging.Logger
	frequency string
	conf      config.VaultEnv
	handler   *policy.VaultPolicyHandler
}

func NewVaultPolicyWatcher(log logging.Logger, frequency string) (*VaultPolicyWatcher, error) {
	conf, err := config.GetVaultEnv()
	if err != nil {
		return nil, err
	}
	return &VaultPolicyWatcher{
		log:       log,
		frequency: frequency,
		conf:      conf,
		handler:   policy.NewVaultPolicyHandler(log),
	}, nil
}

func (v *VaultPolicyWatcher) CronSpec() string {
	return v.frequency
}

func (v *VaultPolicyWatcher) Run() {
	v.log.Debug("started vault policy watcher")
	vc, err := client.NewVaultClientForTokenFromEnv(v.log, v.conf)
	if err != nil {
		v.log.Errorf("%s", err)
		return
	}

	ctx := context.Background()
	if err := v.handler.EnsureKVMounted(ctx, vc); err != nil {
		v.log.Errorf("failed to check vault kv secret mount, %v", err)
	}

	if err := v.handler.UpdateVaultPolicies(ctx, vc); err != nil {
		v.log.Errorf("failed to update vault policies, %v", err)
	}

	if err := v.handler.UpdateVaultRoles(ctx, vc); err != nil {
		v.log.Errorf("failed to update roles, %v", err)
	}
}
