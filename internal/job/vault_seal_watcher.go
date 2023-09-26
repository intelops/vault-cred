package job

import (
	"fmt"

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
	var followerClients []*client.VaultClient
	var leaderVaultClients []*client.VaultClient
	for index, nodeAddress := range v.conf.NodeAddresses {
		conf := v.conf
		conf.Address = nodeAddress
		vc, err := client.NewVaultClient(v.log, conf)
		if err != nil {
			return err
		}
		// assuming that first node address is the leader
		// Todo: keep leaders and followers in the configuration separate
		if index == 0 {
			leaderVaultClients = append(leaderVaultClients, vc)
		} else {
			followerClients = append(followerClients, vc)
		}
	}

	if leaderSealed, err := v.areNodesSealed(leaderVaultClients); err != nil {
		return err
	} else if leaderSealed {
		v.log.Info("vault is sealed for leader nodes")
		for _, vc := range leaderVaultClients {
			err = vc.Unseal()
			if err != nil {
				return fmt.Errorf("failed to unseal vault for leader node, %v", err)
			}
		}
	}

	if followersSealed, err := v.areNodesSealed(followerClients); err != nil {
		return err
	} else if followersSealed {
		for index, vc := range followerClients {
			err = vc.JoinRaftCluster(v.conf.LeaderAPIAddr)
			if err != nil {
				return fmt.Errorf("failed to join the HA cluster by node index: %v, %v", index+1, err)
			}

			err = vc.Unseal()
			if err != nil {
				return fmt.Errorf("failed to unseal vault for node index: %v, %v", index+1, err)
			}
		}

	}

	return nil
}

func (v *VaultSealWatcher) areNodesSealed(vaultClients []*client.VaultClient) (bool, error) {
	status := false
	for _, vc := range vaultClients {
		res, err := vc.IsVaultSealed()
		v.log.Info("unseal status for %v is %v", vc, res)
		if err != nil {
			return false, fmt.Errorf("failed to get vault seal status for %s, %v", v.conf.Address, err)
		}

		if !res {
			return false, nil
		} else {
			v.log.Info("vault node %s is sealed", v.conf.Address)
		}
		status = res
	}
	return status, nil
}
