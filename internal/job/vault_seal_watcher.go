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
	addresses := []string{
		v.conf.Address,
		v.conf.Address2,
		v.conf.Adddress3,
	}
	k8sclient, err := client.NewK8SClient(v.log)
	if err != nil {
		v.log.Errorf("Error while connecting with k8s %s", err)
		return
	}
	podname, err := k8sclient.GetVaultPodInstances(context.Background())
	if err != nil {
		v.log.Errorf("Error while retrieving vault instances %s", err)
		return
	}

	var vc *client.VaultClient

	var vaultClients []*client.VaultClient
	for _, address := range addresses {
		conf := v.conf
		conf.Address = address
		vc, err := client.NewVaultClient(v.log, conf)

		if err != nil {
			v.log.Errorf("%s", err)
			return
		}

		vaultClients = append(vaultClients, vc)
	}

	if v.conf.HAEnabled {

		v.log.Infof("HA ENABLED", v.conf.HAEnabled)

		for _, svc := range podname {

			switch svc {
			case "vault-hash-0":
				vc = vaultClients[0]

				podip, err := vc.GetPodIP(svc, v.conf.VaultSecretNameSpace)
				if err != nil {
					v.log.Errorf("failed to retrieve pod ip, %s", err)
					return
				}

				v.conf.LeaderPodIp = podip
			case "vault-hash-1":
				vc = vaultClients[1]

			case "vault-hash-2":
				vc = vaultClients[2]

			default:

			}

			podip, err := vc.GetPodIP(svc, v.conf.VaultSecretNameSpace)
			if err != nil {
				v.log.Errorf("failed to retrieve pod ip, %s", err)
				return
			}

			res, err := vc.IsVaultSealedForAllInstances(podip)
			if err != nil {
				v.log.Errorf("failed to get vault seal status, %s", err)
				return
			}
			v.log.Info("Seal Status for  %v", podip, res)
			if res {
				v.log.Info("vault is sealed, trying to unseal")
				if svc == "vault-hash-0" {

					v.log.Info("Unsealing for first instance")
					podip, err := vc.GetPodIP(svc, v.conf.VaultSecretNameSpace)
					v.conf.LeaderPodIp = podip
					v.log.Info("Leader Ip", v.conf.LeaderPodIp)
					if err != nil {
						v.log.Errorf("failed to retrieve pod ip, %s", err)
						return
					}
					err = vc.Unseal()
					if err != nil {
						v.log.Errorf("failed to unseal vault, %s", err)
						return
					}

				} else {
					v.log.Info("Leader Pod Ip", v.conf.LeaderPodIp)
					leaderaddr, err := vc.LeaderAPIAddr(v.conf.LeaderPodIp)
					if err != nil {
						v.log.Errorf("failed to retrieve leader address, %s", err)
						return
					}
					v.log.Info("Leader Address", leaderaddr)
					podip, err := vc.GetPodIP(svc, v.conf.VaultSecretNameSpace)
					v.log.Infof("Unsealing for  %v instance", podip)
					if err != nil {
						v.log.Errorf("failed to retrieve pod ip, %s", err)
						return
					}

					err = vc.JoinRaftCluster(podip, leaderaddr)
					if err != nil {
						v.log.Errorf("Failed to join the HA cluster: %v\n", err)
						return

					}
					err = vc.Unseal()
					if err != nil {
						v.log.Errorf("failed to unseal vault, %s", err)
						return
					}

				}

			}

		}

	}
}
