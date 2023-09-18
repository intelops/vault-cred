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
	//	vc, err := client.NewVaultClient(v.log, v.conf)
	// if err != nil {
	// 	v.log.Errorf("%s", err)
	// 	return
	// }
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
	//servicename := []string{"vault-hash-0", "vault-hash-1", "vault-hash-2"}

	var vc *client.VaultClient
	var leaderpodip string
	var vaultClients []*client.VaultClient
	for _, address := range addresses {
		conf := v.conf // Make a copy of the existing configuration
		conf.Address = address
		// conf := config.VaultEnv{
		// 	Address:     address,
		// //	ReadTimeout: 30,
		// //	MaxRetries:  3,
		// }
		v.log.Debug("Address Configuration", conf)

		vc, err := client.NewVaultClient(v.log, conf)

		if err != nil {
			v.log.Errorf("%s", err)
			return
		}

		vaultClients = append(vaultClients, vc)
	}
	v.log.Debug("Vault Clients", vaultClients)

	if v.conf.HAEnabled {

		v.log.Infof("HA ENABLED", v.conf.HAEnabled)

		for _, svc := range podname {
			switch svc {
			case "vault-hash-0":
				vc = vaultClients[0]
				v.log.Debug("Vault Client", vc)

			case "vault-hash-1":
				vc = vaultClients[1]
				v.log.Debug("Vault Client", vc)
			case "vault-hash-2":
				vc = vaultClients[2]
				v.log.Debug("Vault Client", vc)
			default:
				// Handle the case where the service name doesn't match any of the instances
			}
			podip, err := vc.GetPodIP(svc, v.conf.VaultSecretNameSpace)
			if err != nil {
				v.log.Errorf("failed to retrieve pod ip, %s", err)
				return
			}
			v.log.Info("POD IP", podip)
			res, err := vc.IsVaultSealedForAllInstances(podip)
			if err != nil {
				v.log.Errorf("failed to get vault seal status, %s", err)
				return
			}
			v.log.Info("Seal Status", res)
			if res {
				v.log.Info("vault is sealed, trying to unseal")
				if svc == "vault-hash-0" {

					v.log.Info("Unsealing for first instance")
					// _, unsealKeys, err := vc.GetVaultSecretValuesforMultiInstance()
					// if err != nil {
					// 	v.log.Errorf("Failed to fetch the credential: %v\n", err)
					// 	return
					// }
					//err = vc.UnsealVaultInstance(podip,unsealKeys)
					err := vc.Unseal()
					if err != nil {
						v.log.Errorf("failed to unseal vault, %s", err)
						return
					}
					podip, err := vc.GetPodIP(svc, "default")
					v.log.Info("Unsealing for second % vinstance", podip)
					if err != nil {
						v.log.Errorf("failed to retrieve pod ip, %s", err)
						return
					}
					leaderpodip = podip
					v.log.Info("Leader Ip", leaderpodip)

				} else {
					leaderaddr, err := vc.LeaderAPIAddr(leaderpodip)
					if err != nil {
						v.log.Errorf("failed to retrieve leader address, %s", err)
						return
					}
					v.log.Info("Leader Address", leaderaddr)
					podip, err := vc.GetPodIP(svc, "default")
					v.log.Info("Unsealing for second % vinstance", podip)
					if err != nil {
						v.log.Errorf("failed to retrieve pod ip, %s", err)
						return
					}
					v.log.Debug("POD IP", podip)
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
					// _, unsealKeys, err := vc.GetVaultSecretValuesforMultiInstance()
					// v.log.Debug("Unseal Keys", unsealKeys)
					// if err != nil {
					// 	v.log.Errorf("Failed to fetch the credential: %v\n", err)
					// 	return
					// }

					// err = vc.UnsealVaultInstance(podip, unsealKeys)

					// if err != nil {
					// 	v.log.Errorf("failed to unseal vault, %s", err)
					// 	return
					// }

				}

			}

		}
		for _, svc := range podname {
			podip, _ := vc.GetPodIP(svc, "default")
			res, err := vc.IsVaultSealedForAllInstances(podip)

			v.log.Debug("Seal Status of %v :%v", svc, res)
			if err != nil {
				v.log.Errorf("failed to get vault seal status, %s", err)
				return
			}
			v.log.Infof("vault sealed status: %v", res)
		}
	}
}
