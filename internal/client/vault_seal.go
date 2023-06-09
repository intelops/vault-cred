package client

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"

	"path/filepath"
)

func (vc *VaultClient) IsVaultSealed() (bool, error) {
	status, err := vc.C.Sys().SealStatus()
	if err != nil {
		return false, err
	}
	return status.Sealed, nil
}

func (vc *VaultClient) Unseal() error {
	status, err := vc.C.Sys().SealStatus()
	if err != nil {
		return err
	}

	if !status.Initialized {
		err = vc.initializeVaultSecret()
		if err != nil {
			return err
		}
	}

	var unsealKeys []string
	if vc.conf.VaultTokenPath != "" {
		unsealKeys, err = vc.readUnsealKeysFromPath()
		if err != nil {
			return err
		}
	} else {
		unsealKeys, err = vc.readUnsealKeysFromSecret()
		if err != nil {
			return err
		}
	}

	for _, key := range unsealKeys {
		_, err := vc.C.Sys().Unseal(key)
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

	stringData := make(map[string]string)
	for i, value := range unsealKeys {
		key := fmt.Sprintf("%s%d", vc.conf.VaultSecretUnSealKeyPrefix, i+1)
		stringData[key] = value
	}

	stringData[vc.conf.VaultSecretTokenKeyName] = rootToken
	k8s, err := NewK8SClient(vc.log)
	if err != nil {
		return errors.WithMessage(err, "error initializing k8s client")
	}
	err = k8s.CreateOrUpdateSecret(context.Background(), vc.conf.VaultSecretName, vc.conf.VaultSecretNameSpace, stringData)
	if err != nil {
		return errors.WithMessage(err, "error creating vault secret")
	}
	return nil
}

func (vc *VaultClient) generateUnsealKeys() ([]string, string, error) {
	res := &api.InitRequest{
		SecretThreshold: 2,
		SecretShares:    3,
	}

	unsealKeys := []string{}
	initRes, err := vc.C.Sys().Init(res)
	if err != nil {
		return nil, "", err
	}

	unsealKeys = append(unsealKeys, initRes.Keys...)
	rootToken := initRes.RootToken
	return unsealKeys, rootToken, err
}

func (vc *VaultClient) readUnsealKeysFromPath() ([]string, error) {
	files, err := ioutil.ReadDir(vc.conf.VaultUnSealKeyPath)
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(files))
	for _, file := range files {
		if !file.IsDir() {
			keyPath := filepath.Join(vc.conf.VaultUnSealKeyPath, file.Name())
			key, err := readFileContent(keyPath)
			if err != nil {
				return nil, errors.WithMessagef(err, "error in reading unseal key file %s", keyPath)
			}
			keys = append(keys, key)
		}
	}
	return keys, nil
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

	unsealKeys := []string{}
	for key, val := range vaultSec {
		if strings.HasPrefix(key, vc.conf.VaultSecretUnSealKeyPrefix) {
			unsealKeys = append(unsealKeys, val)
		}
	}
	return unsealKeys, nil
}
