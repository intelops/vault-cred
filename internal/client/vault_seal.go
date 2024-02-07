package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
)

func (vc *VaultClient) IsVaultSealed() (bool, error) {
	status, err := vc.c.Sys().SealStatus()
	if err != nil {
		return false, err
	}
	return status.Sealed, nil
}

func (vc *VaultClient) Unseal() error {

	status, err := vc.c.Sys().SealStatus()
	if err != nil {
		return err
	}

	if !status.Sealed {
		return nil
	}

	rootToken, unsealKeys, err := vc.getVaultSecretValues()
	if err != nil {
		return err
	}

	if !status.Initialized && len(rootToken) == 0 && len(unsealKeys) == 0 {
		vc.log.Debug("intializing vault secret")
		err = vc.initializeVaultSecret()
		if err != nil {
			return err
		}
	}

	vc.log.Debugf("found %d vault unseal keys and roottoken length %d", len(unsealKeys), len(rootToken))
	for _, key := range unsealKeys {
		_, err := vc.c.Sys().Unseal(key)
		if err != nil {
			return errors.WithMessage(err, "error while unsealing")
		}
	}
	return nil
}

func (vc *VaultClient) initializeVaultSecret() error {

	unsealKeys, rootToken, err := vc.generateUnsealKeys()

	if err != nil {
		return errors.WithMessage(err, "error while generating unseal keys")
	}

	stringData := make(map[string][]byte)
	for i, value := range unsealKeys {
		key := fmt.Sprintf("%s%d", vc.conf.VaultSecretUnSealKeyPrefix, i+1)
		stringData[key] = []byte(value)
	}

	stringData[vc.conf.VaultSecretTokenKeyName] = []byte(rootToken)
	k8s, err := NewK8SClient(vc.log)
	if err != nil {
		return errors.WithMessage(err, "error initializing k8s client")
	}
	err = k8s.CreateOrUpdateSecret(context.Background(), vc.conf.VaultSecretNameSpace, vc.conf.VaultSecretName, v1.SecretTypeOpaque, stringData, nil)

	if err != nil {
		return errors.WithMessage(err, "error creating vault secret")
	}
	vc.log.Debugf("vault secret updated with %d vaules", len(stringData))
	return nil
}

func (vc *VaultClient) generateUnsealKeys() ([]string, string, error) {
	res := &api.InitRequest{
		SecretThreshold: 2,
		SecretShares:    3,
	}

	unsealKeys := []string{}
	initRes, err := vc.c.Sys().Init(res)
	if err != nil {
		return nil, "", err
	}

	unsealKeys = append(unsealKeys, initRes.Keys...)
	rootToken := initRes.RootToken
	return unsealKeys, rootToken, err
}

func (vc *VaultClient) Leader() (string, error) {
	res, err := vc.c.Sys().Leader()
	if err != nil {
		return "", err
	}
	return res.LeaderAddress, nil
}

func (vc *VaultClient) getVaultSecretValues() (string, []string, error) {
	k8s, err := NewK8SClient(vc.log)
	if err != nil {
		return "", nil, errors.WithMessage(err, "error initializing k8s client")
	}

	vaultSec, err := k8s.GetSecret(context.Background(), vc.conf.VaultSecretName, vc.conf.VaultSecretNameSpace)
	if err != nil {
		if strings.Contains(err.Error(), "secret not found") {
			vc.log.Debugf("secret %d not found", vc.conf.VaultSecretName)
			return "", nil, nil
		}
		return "", nil, errors.WithMessage(err, "error fetching vault secret")
	}

	vc.log.Debugf("found %d vault secret values", len(vaultSec.Data))
	unsealKeys := []string{}
	var rootToken string
	for key, val := range vaultSec.Data {
		if strings.HasPrefix(key, vc.conf.VaultSecretUnSealKeyPrefix) {
			unsealKeys = append(unsealKeys, val)
			continue
		}
		if strings.EqualFold(key, vc.conf.VaultSecretTokenKeyName) {
			rootToken = val
		}
	}
	return rootToken, unsealKeys, nil
}
