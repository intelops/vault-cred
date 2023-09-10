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
	servicename := []string{"capten-dev-vault-0", "capten-dev-vault-1", "capten-dev-vault-2"}
	// res, err := vc.IsVaultSealed()
	// if err != nil {
	// 	v.log.Errorf("failed to get vault seal status, %s", err)
	// 	return
	// }

	if v.conf.HAEnabled {
		v.log.Infof("HA ENABLED", v.conf.HAEnabled)
		// res, err := vc.IsVaultSealed()
		// if err != nil {
		// 	v.log.Errorf("failed to get vault seal status, %s", err)
		// 	return
		// }
		for _, svc := range servicename {

			res, err := vc.IsVaultSealedForAllInstances(svc)
			if err != nil {
				v.log.Errorf("failed to get vault seal status, %s", err)
				return
			}
			if res {
				if svc == "capten-dev-vault-0" {
					vc.Unseal()
				} else {
					_, unsealKeys, err := vc.GetVaultSecretValuesforMultiInstance()
					if err != nil {
						v.log.Errorf("Failed to fetch the credential: %v\n", err)
						return
					}
					key := unsealKeys[0]
					vc.UnsealVaultInstance(svc, key)
					err = vc.JoinRaftCluster()
					if err != nil {
						v.log.Errorf("Failed to join the HA cluster: %v\n", err)
						return

					}

				}

			}

			//res, err = vc.IsVaultSealed()
			// Perform the unseal operation on the Vault instance within the pod using the podIP
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

	// if res {
	// 	v.log.Info("vault is sealed, trying to unseal")

	// 	err := vc.Unseal()
	// 	if err != nil {
	// 		v.log.Errorf("failed to unseal vault, %s", err)
	// 		return
	// 	}
	// 	v.log.Info("vault unsealed executed")

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
