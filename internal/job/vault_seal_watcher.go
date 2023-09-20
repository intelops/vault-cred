package job

import (
	"context"

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

	if v.conf.HAEnabled {
		v.log.Infof(" Vault HA ENABLED", v.conf.HAEnabled)
		err := v.Unseal_HA_Enabled()
		if err != nil {
			v.log.Errorf("%s", err)
			return
		}

	} else {
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

		if res {
			v.log.Info("vault is sealed, trying to unseal")
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
}
func (v *VaultSealWatcher) Unseal_HA_Enabled() error {

	addresses := []string{
		v.conf.Address,
		v.conf.Address2,
		v.conf.Adddress3,
	}
	k8sclient, err := client.NewK8SClient(v.log)
	if err != nil {
		v.log.Errorf("Error while connecting with k8s %s", err)
		return err
	}
	podname, err := k8sclient.GetVaultPodInstances(context.Background())
	if err != nil {
		v.log.Errorf("Error while retrieving vault instances %s", err)
		return err
	}

	var vc *client.VaultClient

	var vaultClients []*client.VaultClient
	for _, address := range addresses {
		conf := v.conf
		conf.Address = address
		vc, err := client.NewVaultClient(v.log, conf)

		if err != nil {
			v.log.Errorf("%s", err)
			return err
		}

		vaultClients = append(vaultClients, vc)
	}

	for i, svc := range podname {

		vc = vaultClients[i]

		res, err := vc.IsVaultSealedForAllInstances(svc)
		if err != nil {
			v.log.Errorf("failed to get vault seal status, %s", err)
			return err
		}
		v.log.Info("Seal Status for  %v", svc, res)
		if res {

			if i == 0 {

				v.log.Info("Unsealing for first instance")

				v.conf.LeaderPodIp = svc

				if err != nil {
					v.log.Errorf("failed to retrieve pod ip, %s", err)
					return err
				}
				err = vc.Unseal()
				if err != nil {
					v.log.Errorf("failed to unseal vault, %s", err)
					return err
				}

			} else {

				leaderaddr, err := vc.LeaderAPIAddr(v.conf.LeaderPodIp)
				if err != nil {
					v.log.Errorf("failed to retrieve leader address, %s", err)
					return err
				}
				v.log.Debug("Leader Address", leaderaddr)

				v.log.Infof("Unsealing for  %v instance", podname)

				err = vc.JoinRaftCluster(svc, leaderaddr)
				if err != nil {
					v.log.Errorf("Failed to join the HA cluster: %v\n", err)
					return err

				}
				err = vc.Unseal()
				if err != nil {
					v.log.Errorf("failed to unseal vault, %s", err)
					return err
				}

			}

		}

	}
	return nil
}
