package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	Host                     string `envconfig:"HOST" default:"0.0.0.0"`
	Port                     int    `envconfig:"PORT" default:"9098"`
	VaultSealWatchInterval   string `envconfig:"VAULT_SEAL_WATCH_INTERVAL"`
	VaultPolicyWatchInterval string `envconfig:"VAULT_POLICY_WATCH_INTERVAL"`
	VaultCredSyncInterval    string `envconfig:"VAULT_CRED_SYNC_INTERVAL"`

}

type VaultEnv struct {
	HAEnabled                  bool          `envconfig:"HA_ENABLED" default:"true"`
	Address                    string        `envconfig:"VAULT_ADDR" required:"true"`
	NodeAddresses              []string      `envconfig:"VAULT_NODE_ADDRESSES" required:"true"`
	CACert                     string        `envconfig:"VAULT_CACERT" required:"false"`
	ReadTimeout                time.Duration `envconfig:"VAULT_READ_TIMEOUT" default:"60s"`
	MaxRetries                 int           `envconfig:"VAULT_MAX_RETRIES" default:"5"`
	VaultTokenForRequests      bool          `envconfig:"VAULT_TOKEN_FOR_REQUESTS" default:"true"`
	VaultSecretName            string        `envconfig:"VAULT_SECRET_NAME" default:"vault-server"`
	VaultSecretNameSpace       string        `envconfig:"POD_NAMESPACE" required:"true"`
	VaultSecretTokenKeyName    string        `envconfig:"VAULT_SECRET_TOKEN_KEY_NAME" default:"root-token"`
	VaultSecretUnSealKeyPrefix string        `envconfig:"VAULT_SECRET_UNSEAL_KEY_PREFIX" default:"unsealkey"`
	VaultToken                 string        `envconfig:"VAULT_TOKEN"`
	VaultCredSyncSecretName    string        `envconfig:"VAULT_CRED_SYNC_SECRET_NAME" default:"vault-cred-sync-data"`
}

func FetchConfiguration() (Configuration, error) {
	cfg := Configuration{}
	err := envconfig.Process("", &cfg)
	return cfg, err
}

func GetVaultEnv() (VaultEnv, error) {
	cfg := VaultEnv{}
	err := envconfig.Process("", &cfg)
	return cfg, err
}


