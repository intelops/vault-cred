package job

import (
	"github.com/intelops/go-common/logging"
	"github.com/intelops/vault-cred/config"
	"github.com/intelops/vault-cred/internal/client"
)

type VaultSealWatcher struct {
	log       logging.Logger
	conf      config.VaultEnv
	frequency string
}

func NewVaultSealWatcher(log logging.Logger, frequency string) (*VaultSealWatcher, error) {
	conf, err := config.GetVaultEnv()
	if err != nil {
		return nil, err
	}
	return &VaultSealWatcher{
		log:       log,
		frequency: frequency,
		conf:      conf,
	}, nil
}

func (v *VaultSealWatcher) CronSpec() string {
	return v.frequency
}

func (v *VaultSealWatcher) Run() {
	v.log.Debug("started vault seal watcher job")
	vc, err := client.NewVaultClient(v.log, v.conf)
	if err != nil {
		v.log.Errorf("%s", err)
		return
	}

	res, err := vc.IsVaultSealed()
	if err != nil {
		v.log.Errorf("failed to get vault seal status, %s", err)
		return
	}
	if v.conf.HAEnabled {
		v.log.Infof("HA ENABLED", v.conf.HAEnabled)
		err := vc.JoinRaftCluster()
		if err != nil {
			v.log.Errorf("Failed to join the HA cluster: %v\n", err)

		}
	}

	if res {
		v.log.Info("vault is sealed, trying to unseal")
		if v.conf.HAEnabled {
			v.log.Infof("HA Enabled", v.conf.HAEnabled)
			err := vc.JoinRaftCluster()
			if err != nil {
				v.log.Errorf("Failed to join the HA cluster: %v\n", err)

			}
		}
		err := vc.Unseal()
		if err != nil {
			v.log.Errorf("failed to unseal vault, %s", err)
			return
		}
		v.log.Info("vault unsealed executed")

		res, err := vc.IsVaultSealed()
		if err != nil {
			v.log.Errorf("failed to get vault seal status, %s", err)
			return
		}
		v.log.Infof("vault sealed status: %v", res)
		return
	} else {
		v.log.Debug("vault is in unsealed status")
	}
}
