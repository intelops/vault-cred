package job

import (
	"fmt"

	"github.com/intelops/go-common/logging"
	"github.com/intelops/vault-cred/config"
	"github.com/intelops/vault-cred/internal/client"
)

// this var is set to true after leader created for the very first time
// remove this later, after learning the correct usecase of Leader() api
var leaderCreated bool

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
	v.log.Debugf("started vault seal watcher job with vault HA: %v", v.conf.HAEnabled)

	if v.conf.HAEnabled {
		if len(v.conf.NodeAddresses) != 3 {
			v.log.Errorf("vault HA node count %d is not valid", len(v.conf.NodeAddresses))
			return
		}

		if err := v.handleUnsealForHAVault(); err != nil {
			v.log.Errorf("%s", err)
		}
	} else {
		if err := v.handleUnsealForNonHAVault(); err != nil {
			v.log.Errorf("%s", err)
		}
	}
}

func (v *VaultSealWatcher) handleUnsealForNonHAVault() error {
	vc, err := client.NewVaultClient(v.log, v.conf)
	if err != nil {
		return err
	}

	res, err := vc.IsVaultSealed()
	if err != nil {
		return fmt.Errorf("failed to get vault seal status, %s", err)
	}

	if res {
		v.log.Info("vault is sealed, trying to unseal")
		err := vc.Unseal()
		if err != nil {
			return fmt.Errorf("failed to unseal vault, %s", err)
		}
		v.log.Info("vault unsealed executed")

		res, err := vc.IsVaultSealed()
		if err != nil {
			return fmt.Errorf("failed to get vault seal status, %s", err)
		}
		v.log.Infof("vault sealed status: %v", res)
	} else {
		v.log.Debug("vault is in unsealed status")
	}
	return nil
}

func (v *VaultSealWatcher) handleUnsealForHAVault() error {
	var vaultClients []*client.VaultClient
	for _, nodeAddress := range v.conf.NodeAddresses {
		conf := v.conf
		conf.Address = nodeAddress
		vc, err := client.NewVaultClient(v.log, conf)
		if err != nil {
			return err
		}
		vaultClients = append(vaultClients, vc)
	}

	for index, vc := range vaultClients {
		leaderNode := v.conf.LeaderAPIAddr

		if leader, err := vc.Leader(); err == nil && leader != "" {
			leaderNode = leader
		}

		if leaderCreated {
			err := vc.JoinRaftCluster(leaderNode)
			if err != nil {
				return fmt.Errorf("failed to join the HA cluster by node index: %v, %v", index+1, err)
			}
		}

		err := vc.Unseal()
		if err != nil {
			return fmt.Errorf("failed to unseal vault for node index: %v, %v", index+1, err)
		}
		leaderCreated = true
	}

	return nil
}
