package client

import (
	"context"
	"encoding/base64"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/vault/api"
	vaultauth "github.com/hashicorp/vault/api/auth/kubernetes"
	"github.com/intelops/go-common/logging"
	"github.com/intelops/vault-cred/config"
	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
)

const (
	vaultRoleKey    string = "vault-role"
	serviceTokenKey string = "service-token"
)

type VaultClient struct {
	c    *api.Client
	conf config.VaultEnv
	log  logging.Logger
}

func NewVaultClientForServiceAccount(ctx context.Context, log logging.Logger, conf config.VaultEnv) (*VaultClient, error) {
	if conf.VaultTokenForRequests {
		return NewVaultClientForVaultToken(log, conf)
	}

	vc, err := NewVaultClient(log, conf)
	if err != nil {
		return nil, err
	}

	err = vc.configureAuthToken(ctx)
	if err != nil {
		return nil, err
	}
	return vc, nil
}

func NewVaultClientForVaultToken(log logging.Logger, conf config.VaultEnv) (*VaultClient, error) {
	vc, err := NewVaultClient(log, conf)
	if err != nil {
		return nil, err
	}

	k8s, err := NewK8SClient(vc.log)
	if err != nil {
		return nil, errors.WithMessage(err, "error initializing k8s client")
	}
	vaultSec, err := k8s.GetSecret(context.Background(), vc.conf.VaultSecretName, vc.conf.VaultSecretNameSpace)
	if err != nil {
		return nil, errors.WithMessage(err, "error fetching vault secret")
	}

	rootToken := vaultSec[vc.conf.VaultSecretTokenKeyName]
	if len(rootToken) == 0 {
		return nil, errors.New("vault root token not found")
	}
	vc.c.SetToken(rootToken)
	return vc, nil
}

func NewVaultClient(log logging.Logger, conf config.VaultEnv) (*VaultClient, error) {
	cfg, err := prepareVaultConfig(conf)
	if err != nil {
		return nil, err
	}

	c, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &VaultClient{
		c:    c,
		conf: conf,
		log:  log,
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

func (vc *VaultClient) configureAuthToken(ctx context.Context) (err error) {
	metadata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.WithMessagef(err, "vault auth context is missing")
	}
	roleData := metadata[vaultRoleKey]
	tokenData := metadata[serviceTokenKey]
	if !(len(roleData) == 1 && len(tokenData) == 1) {
		return errors.WithMessagef(err, "vault auth context is missing")
	}

	serviceToken, err := base64.StdEncoding.DecodeString(tokenData[0])
	if !ok {
		return errors.WithMessagef(err, "vault auth context decoding error")
	}

	k8sAuth, err := vaultauth.NewKubernetesAuth(
		roleData[0],
		vaultauth.WithServiceAccountToken(string(serviceToken)),
	)
	if err != nil {
		return errors.WithMessagef(err, "error in initializing Kubernetes auth method")
	}

	authInfo, err := vc.c.Auth().Login(ctx, k8sAuth)
	if err != nil {
		return errors.WithMessagef(err, "error in login with Kubernetes auth")
	}
	if authInfo == nil {
		return errors.New("no auth info was returned after login")
	}
	return nil
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
