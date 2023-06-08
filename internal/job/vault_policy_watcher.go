package job

import (
	"github.com/intelops/go-common/logging"
	"github.com/intelops/vault-cred/config"
	"github.com/intelops/vault-cred/internal/client"
)

type VaultPolicyWatcher struct {
	log       logging.Logger
	conf      config.VaultEnv
	frequency string
}

func NewVaultPolicyWatcher(log logging.Logger, frequency string) (*VaultPolicyWatcher, error) {
	conf, err := config.GetVaultEnv()
	if err != nil {
		return nil, err
	}
	return &VaultPolicyWatcher{
		log:       log,
		frequency: frequency,
		conf:      conf,
	}, nil
}

func (v *VaultPolicyWatcher) CronSpec() string {
	return v.frequency
}

func (v *VaultPolicyWatcher) Run() {
	v.log.Info("started vault policy watcher")
	_, err := client.NewVaultClientForVaultToken(v.log, v.conf)
	if err != nil {
		v.log.Errorf("%s", err)
	}

	// Svc iam creates vault-policy-iam
	// Svc capten creates vault-policy-capen
	// get config maps with prefix "vault-policy-<>"
	// for each config map
	// read SA name, credentail-acess : [] {"credential-type", "access type"}
	//"serviceAccount" : "iam-sa"
	//"credentailAccess" : "system-cred:read, client-cert:admin, xyz: read"

	// prepare role "iam-sa-role" and policy "iam-sa-policy"  with path
	// paths are based on credential types
	// secret/data/system-cred/* with read,  secret/data/client-cert/* with read, write, list, delete
	// secret/data/xyz/* read
	// create policy

	// is policy exists, is any differnet from vault-policy data
	// update delta
}
