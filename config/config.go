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
}

type VaultEnv struct {
	Address               string        `envconfig:"VAULT_ADDR" required:"true"`
	CACert                string        `envconfig:"VAULT_CACERT" required:"false"`
	ReadTimeout           time.Duration `envconfig:"VAULT_READ_TIMEOUT" default:"60s"`
	MaxRetries            int           `envconfig:"VAULT_MAX_RETRIES" default:"5"`
	VaultTokenForRequests bool          `envconfig:"VAULT_TOKEN_FOR_REQUESTS" default:"false"`
	VaultTokenPath        string        `envconfig:"VAULT_TOKEN_PATH"`
	VaultUnSealKeyPath    string        `envconfig:"VAULT_UNSEAL_PATH"`
	VaultUnSealKeyNames   []string      `envconfig:"VAULT_UNSEAL_KEY_NAMES" default:"unsealkey1,unsealkey2,unsealkey3"`
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
