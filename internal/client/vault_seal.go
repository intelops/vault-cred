package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
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

	if !status.Initialized {
		err = vc.initializeVaultSecret()
		if err != nil {
			return err
		}
	}

	unsealKeys, err := vc.readUnsealKeysFromSecret()
	if err != nil {
		return err
	}

	vc.log.Debugf("found %d vault unseal keys", len(unsealKeys))
	for _, key := range unsealKeys {
		_, err := vc.c.Sys().Unseal(key)
		if err != nil {
			return errors.WithMessage(err, "error while unsealing")
		}
	}
	return nil
}

func (vc *VaultClient) initializeVaultSecret() error {
	vc.log.Debug("intializing vault secret")
	unsealKeys, rootToken, err := vc.generateUnsealKeys()
	if err != nil {
		return errors.WithMessage(err, "error while generating unseal keys")
	}

	stringData := make(map[string]string)
	for i, value := range unsealKeys {
		key := fmt.Sprintf("%s%d", vc.conf.VaultSecretUnSealKeyPrefix, i+1)
		stringData[key] = value
	}

	stringData[vc.conf.VaultSecretTokenKeyName] = rootToken
	vc.log.Debugf("REMOVE THIS LOG --> vault secret data, %v", stringData)
	k8s, err := NewK8SClient(vc.log)
	if err != nil {
		return errors.WithMessage(err, "error initializing k8s client")
	}
	err = k8s.CreateOrUpdateSecret(context.Background(), vc.conf.VaultSecretName, vc.conf.VaultSecretNameSpace, stringData)
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

func (vc *VaultClient) readUnsealKeysFromSecret() ([]string, error) {
	k8s, err := NewK8SClient(vc.log)
	if err != nil {
		return nil, errors.WithMessage(err, "error initializing k8s client")
	}
	vaultSec, err := k8s.GetSecret(context.Background(), vc.conf.VaultSecretName, vc.conf.VaultSecretNameSpace)
	if err != nil {
		return nil, errors.WithMessage(err, "error creating vault secret")
	}

	vc.log.Debugf("found %d vault secret vaules", len(vaultSec))
	unsealKeys := []string{}
	for key, val := range vaultSec {
		vc.log.Debugf("REMOVE THIS LOG -->  check prefix %s for %s : %s", vc.conf.VaultSecretUnSealKeyPrefix, key, val)
		if strings.HasPrefix(key, vc.conf.VaultSecretUnSealKeyPrefix) {
			unsealKeys = append(unsealKeys, val)
		}
	}
	return unsealKeys, nil
}
