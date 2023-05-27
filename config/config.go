package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	Host string `envconfig:"HOST" default:"0.0.0.0"`
	Port int    `envconfig:"PORT" default:"9098"`
}

type VaultEnv struct {
	Address               string        `envconfig:"VAULT_ADDR" required:"true"`
	CACert                string        `envconfig:"VAULT_CACERT" required:"true"`
	ReadTimeout           time.Duration `envconfig:"VAULT_READ_TIMEOUT" default:"60s"`
	MaxRetries            int           `envconfig:"VAULT_MAX_RETRIES" default:"5"`
	VaultTokenForRequests bool          `envconfig:"VAULT_TOKEN_FOR_REQUESTS" default:"false"`
	VaultTokenPath        string        `envconfig:"VAULT_TOKEN_PATH"`
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
