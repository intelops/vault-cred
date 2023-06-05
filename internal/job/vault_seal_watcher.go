package job

import (
	"log"

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
	_, err := client.NewVaultClientForVaultToken(v.conf)
	if err != nil {
		v.log.Errorf("%s", err)
	}
	res,_:=client.VaultClient{}.IsVaultSealed()
	  // get vault status
       // if status unsealed, then return
	   if res{
           err:=client.VaultClient{}.Unseal()
		   if err!=nil {
			log.Fatal("Error while unsealing",err)
		   }
	   }

      // if status sealed, then unseal
}
