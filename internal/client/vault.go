package client

import (
	"context"
	"os"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/vault/api"
	vaultauth "github.com/hashicorp/vault/api/auth/kubernetes"
	"github.com/intelops/vault-cred/config"
	"github.com/pkg/errors"
)

type VaultClient struct {
	c *api.Client
}

func NewVaultClientForServiceAccount(ctx context.Context, conf config.VaultEnv, vaultRole, saToken string) (c *VaultClient, err error) {
	if conf.VaultTokenForRequests {
		return NewVaultClientForVaultToken(conf)
	}

	vc, err := newVaultClient(conf)
	if err != nil {
		return nil, err
	}

	err = configureAuthToken(ctx, vc.c, vaultRole, saToken)
	if err != nil {
		return nil, err
	}
	return vc, nil
}

func NewVaultClientForVaultToken(conf config.VaultEnv) (*VaultClient, error) {
	vc, err := newVaultClient(conf)
	if err != nil {
		return nil, err
	}
	if conf.VaultTokenPath != "" {
		token, err := readFileContent(conf.VaultTokenPath)
		if err != nil {
			return nil, errors.WithMessage(err, "error in reading token file")
		}
		vc.c.SetToken(token)
		return vc, nil
	}
	return nil, errors.New("vault token path not found")
}

func newVaultClient(conf config.VaultEnv) (*VaultClient, error) {
	cfg, err := prepareVaultConfig(conf)
	if err != nil {
		return nil, err
	}

	c, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &VaultClient{
		c: c,
	}, nil
}

func prepareVaultConfig(conf config.VaultEnv) (cfg *api.Config, err error) {
	cfg = api.DefaultConfig()
	cfg.Address = conf.Address
	cfg.Timeout = conf.ReadTimeout
	cfg.Backoff = retryablehttp.DefaultBackoff
	cfg.MaxRetries = conf.MaxRetries
	if conf.CACert != "" {
		tlsConfig := api.TLSConfig{CACert: conf.CACert}
		err = cfg.ConfigureTLS(&tlsConfig)
	}
	return
}

func configureAuthToken(ctx context.Context, vc *api.Client, vaultRole, saToken string) (err error) {
	k8sAuth, err := vaultauth.NewKubernetesAuth(
		vaultRole,
		vaultauth.WithServiceAccountToken(saToken),
	)
	if err != nil {
		return errors.WithMessagef(err, "error in initializing Kubernetes auth method")
	}

	authInfo, err := vc.Auth().Login(ctx, k8sAuth)
	if err != nil {
		return errors.WithMessagef(err, "error in login with Kubernetes auth")
	}
	if authInfo == nil {
		return errors.New("no auth info was returned after login")
	}
	return nil
}

func readFileContent(path string) (s string, err error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return
	}
	s = strings.TrimSpace(string(b))
	return
}

func (vc *VaultClient) GetCredential(ctx context.Context, mountPath, secretPath string) (cred map[string]string, err error) {
	secretValByPath, err := vc.c.KVv2(mountPath).Get(context.Background(), secretPath)
	if err != nil {
		err = errors.WithMessagef(err, "error in reading certificate data from %s", secretPath)
		return
	}

	if secretValByPath == nil {
		err = errors.WithMessagef(err, "crdentaial not found at %s", secretPath)
		return
	}
	if secretValByPath.Data == nil {
		err = errors.WithMessagef(err, "crdentaial data is corrupted for %s", secretPath)
		return
	}
	cred = map[string]string{}
	for key, val := range secretValByPath.Data {
		cred[key] = val.(string)
	}
	return
}

func (vc *VaultClient) PutCredential(ctx context.Context, mountPath, secretPath string, cred map[string]string) (err error) {
	credData := map[string]interface{}{}
	for key, val := range cred {
		credData[key] = val
	}
	_, err = vc.c.KVv2(mountPath).Put(ctx, secretPath, credData)
	if err != nil {
		err = errors.WithMessagef(err, "error in putting credentail at %s", secretPath)
	}
	return
}

func (vc *VaultClient) DeleteCredential(ctx context.Context, mountPath, secretPath string) (err error) {
	err = vc.c.KVv2(mountPath).Delete(ctx, secretPath)
	if err != nil {
		err = errors.WithMessagef(err, "error in deleting credentail at %s", secretPath)
	}
	return
}
