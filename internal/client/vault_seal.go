package client

import (
	"context"
	"fmt"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"

	"path/filepath"
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

	keys, err := readUnsealKeysFromPath(vc.conf.VaultUnSealKeyPath)
	if err != nil {
		return err
	}

	for _, key := range keys {
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

	stringData := make(map[string]string)
	for i, value := range unsealKeys {
		key := fmt.Sprintf("key%d", i+1)
		stringData[key] = value
	}

	key := "roottoken"
	stringData[key] = rootToken

	secData := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      vc.conf.VaultSecretName,
			Namespace: vc.conf.VaultSecretNameSpace,
		},
		StringData: stringData,
	}

	k8s, err := NewK8SClient(vc.log)
	if err != nil {
		return errors.WithMessage(err, "error initializing k8s client")
	}
	err = k8s.CreateOrUpdateSecret(context.Background(), vc.conf.VaultSecretNameSpace, secData)
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
	initRes, err := vc.c.Sys().Init(res)
	if err != nil {
		return nil, "", err
	}

	unsealKeys = append(unsealKeys, initRes.Keys...)
	rootToken := initRes.RootToken
	return unsealKeys, rootToken, err
}

func readUnsealKeysFromPath(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(files))
	for _, file := range files {
		if !file.IsDir() {
			keyPath := filepath.Join(path, file.Name())
			key, err := readFileContent(keyPath)
			if err != nil {
				return nil, errors.WithMessagef(err, "error in reading unseal key file %s", keyPath)
			}
			keys = append(keys, key)
		}
	}
	return keys, nil
}
