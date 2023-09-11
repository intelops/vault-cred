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
	// addresses := []string{
	// 	v.conf.Address,
	// 	v.conf.Address2,
	// 	v.conf.Adddress3,
	// }
	servicename := []string{"capten-dev-vault-0", "capten-dev-vault-1", "capten-dev-vault-2"}

	// var vc *client.VaultClient
	// var vaultClients []*client.VaultClient
	// for _, address := range addresses {
	// 	conf := config.VaultEnv{
	// 		Address:     address,
	// 		ReadTimeout: 30,
	// 		MaxRetries:  3,

	// 		// Set other configuration options as needed
	// 	}
	// 	v.log.Debug("Address Configuration", conf)

	// 	vc, err := client.NewVaultClient(v.log, v.conf)

	// 	if err != nil {
	// 		v.log.Errorf("%s", err)
	// 		return
	// 	}

	// 	vaultClients = append(vaultClients, vc)
	// }
	// v.log.Debug("Vault Clients", vaultClients)

	if v.conf.HAEnabled {

		v.log.Infof("HA ENABLED", v.conf.HAEnabled)

		for _, svc := range servicename {
			// switch svc {
			// case "capten-dev-vault-0":
			// 	vc = vaultClients[0]
			// 	v.log.Debug("Vault Client:", vc)
			// case "capten-dev-vault-1":
			// 	vc = vaultClients[1]
			// 	v.log.Debug("Vault Client:", vc)
			// case "capten-dev-vault-2":
			// 	vc = vaultClients[2]
			// 	v.log.Debug("Vault Client:", vc)
			// default:
			// 	// Handle the case where the service name doesn't match any of the instances
			// }
			res, err := vc.IsVaultSealedForAllInstances(svc)
			if err != nil {
				v.log.Errorf("failed to get vault seal status, %s", err)
				return
			}
			if res {
				v.log.Info("vault is sealed, trying to unseal")
				if svc == "capten-dev-vault-0" {
					_, unsealKeys, err := vc.GetVaultSecretValuesforMultiInstance()
					v.log.Debug("Unseal Keys",unsealKeys)
					if err != nil {
						v.log.Errorf("Failed to fetch the credential: %v\n", err)
						return
					}
				//	err := vc.Unseal()
					err = vc.UnsealVaultInstance(svc, unsealKeys)
					if err != nil {
						v.log.Errorf("failed to unseal vault, %s", err)
						return
					}

				} else {
					_, unsealKeys, err := vc.GetVaultSecretValuesforMultiInstance()
					v.log.Debug("Unseal Keys", unsealKeys)
					if err != nil {
						v.log.Errorf("Failed to fetch the credential: %v\n", err)
						return
					}
					err = vc.UnsealVaultInstance(svc, unsealKeys)
					//err = vc.Unseal()
					if err != nil {
						v.log.Errorf("failed to unseal vault, %s", err)
						return
					}
					//	v.log.Info("vault unsealed executed")
					err = vc.JoinRaftCluster()
					if err != nil {
						v.log.Errorf("Failed to join the HA cluster: %v\n", err)
						return

					}
				}

				// err := vc.Unseal()
				// if err != nil {
				// 	v.log.Errorf("failed to unseal vault, %s", err)
				// 	return
				// }

				// 	res, err := vc.IsVaultSealed()
				// 	if res {

				// 		err := vc.Unseal()
				// 		if err != nil {
				// 			v.log.Errorf("failed to unseal vault, %s", err)
				// 			return
				// 		}

				// 	}
				// 	if err != nil {
				// 		v.log.Errorf("failed to get vault seal status, %s", err)
				// 		return
				// 	}
				// 	v.log.Infof("vault sealed status: %v", res)

				// } else {
				// 	v.log.Debug("vault is in unsealed status")
				// }
			}

		}
		for _, svc := range servicename {
			res, err := vc.IsVaultSealedForAllInstances(svc)

			v.log.Debug("Seal Status of %v :%v", svc, res)
			if err != nil {
				v.log.Errorf("failed to get vault seal status, %s", err)
				return
			}
			v.log.Infof("vault sealed status: %v", res)
		}
	}
}
