package job
import (
	"github.com/intelops/go-common/logging"
	"github.com/intelops/vault-cred/config"
	"github.com/intelops/vault-cred/internal/client"
)

type ServicePolicyWatcher struct {
	log       logging.Logger
	conf      config.VaultEnv
	frequency string
}

func NewServicePolicyWatcher(log logging.Logger, frequency string) (*ServicePolicyWatcher, error) {
	conf, err := config.GetVaultEnv()
	if err != nil {
		return nil, err
	}
	return &ServicePolicyWatcher{
		log:       log,
		frequency: frequency,
		conf:      conf,
	}, nil
}

func (v *ServicePolicyWatcher) CronSpec() string {
	return v.frequency
}

func (v *ServicePolicyWatcher) Run() {
	_, err := client.NewVaultClientForVaultToken(v.conf)
	if err != nil {
		v.log.Errorf("%s", err)
	}

//	Svc iam creates vault-policy-iam config map
//	Svc capten creates vault-policy-capen
	// get config maps with prefix "vault-policy-<>"
	// for each config map
	// read SA name, credentail-acess : [] {"credential-type", "access type"}
	// "serviceAccount" : "iam-sa"
	// "credentailAccess" : "system-cred:read, client-cert:admin, xyz: read"

	 // prepare role "iam-sa-role" and policy "iam-sa-policy"  with path
	 // paths are based on credential types
	 // secret/data/system-cred/* with read,  secret/data/client-cert/* with read, write, list, delete
	 // secret/data/xyz/* read
	 // create policy
	 
	 // is policy exists, is any differnet from vault-policy data
	 // update delta
}