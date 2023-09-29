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

	sealed, err := v.isAnyNodeSealed(vaultClients)
	if err != nil {
		return err
	}

	if !sealed {
		v.log.Debug("All nodes are unsealed")
		return nil
	}

	var leaderNode string
	for _, vc := range vaultClients {
		if leader, err := vc.Leader(); err == nil && leader != "" {
			leaderNode = leader
			break
		}
	}
	v.log.Infof("Found leader node: %v", leaderNode)

	for index, vc := range vaultClients {
		if len(leaderNode) > 0 {
			err := vc.JoinRaftCluster(leaderNode)
			if err != nil {
				return fmt.Errorf("failed to join the HA cluster by node index: %v, error: %v", index, err)
			}
			v.log.Info("Node %s joined leader %s", v.conf.Address, leaderNode)
		}

		if err := vc.Unseal(); err != nil {
			return fmt.Errorf("failed to unseal vault on node %s, %v", v.conf.Address, err)
		}

		v.log.Info("Node %s successfully Unsealed", v.conf.Address)
	}
	return nil
}

func (v *VaultSealWatcher) isAnyNodeSealed(vaultClients []*client.VaultClient) (bool, error) {
	sealedStatus := false
	for _, vc := range vaultClients {
		res, err := vc.IsVaultSealed()
		if err != nil {
			return false, fmt.Errorf("failed to get vault seal status for %s, %v", v.conf.Address, err)
		}

		if res {
			sealedStatus = true
		}

		v.log.Debugf("vault node %s seal status %s", v.conf.Address, res)
	}
	return sealedStatus, nil
}
