package client

import (
	"context"
	//"encoding/base64"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (vc *VaultClient) IsVaultSealed() (bool, error) {
	status, err := vc.c.Sys().SealStatus()
	if err != nil {
		return false, err
	}
	return status.Sealed, nil
}

func (vc *VaultClient) Unseal(podip string) error {
	address := fmt.Sprintf("http://%s:8200", podip)
	err := vc.c.SetAddress(address)
	if err != nil {
		vc.log.Errorf("Error while setting address")
	}
	vc.log.Debug("Address", address)
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

func (vc *VaultClient) UnsealVaultInstance(podip string, unsealKey []string) error {
	// Create a Vault API client
	vc.log.Debug("Checking Unseal status for vault Instance")
	address := fmt.Sprintf("http://%s:8200", podip)
	err := vc.c.SetAddress(address)
	if err != nil {
		vc.log.Errorf("Error while setting address")
	}
	vc.log.Debug("Address", address)

	for _, key := range unsealKey {
		unsealResponse, err := vc.c.Sys().Unseal(key)
		if err != nil {
			return errors.WithMessage(err, "error while unsealing")
		}
		if unsealResponse.Sealed {
			vc.log.Debug("Vault is still sealed after unsealing attempt")
		}
	}

	// Check if Vault is sealed and unseal if necessary

	// Vault is sealed; unseal it
	// unsealResponse, err := vc.c.Sys().Unseal(unsealKey)
	// if err != nil {
	// 	return err
	// }

	// if unsealResponse.Sealed {
	// 	vc.log.Debug("Vault is still sealed after unsealing attempt")
	// }

	return nil
}

func (vc *VaultClient) GetVaultSecretValuesforMultiInstance() (string, []string, error) {
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
			//	decodedValue, err := base64.StdEncoding.DecodeString(val)
			if err != nil {
				return "", nil, errors.WithMessage(err, "error decoding value")
			}

			unsealKeys = append(unsealKeys, val)
			vc.log.Debug("Unseal Keys", unsealKeys)
			continue
		}
		if strings.EqualFold(key, vc.conf.VaultSecretTokenKeyName) {
			//		decodedValue, err := base64.StdEncoding.DecodeString(val)
			if err != nil {
				return "", nil, errors.WithMessage(err, "error decoding root token")
			}
			rootToken = val
			vc.log.Debug("Root Token Key", rootToken)
		}
	}
	return rootToken, unsealKeys, nil
}

func (vc *VaultClient) IsVaultSealedForAllInstances(svc string) (bool, error) {
	address := fmt.Sprintf("http://%s:8200", svc)
	err := vc.c.SetAddress(address)
	vc.log.Debug("Address for checking vault status", address)
	if err != nil {
		vc.log.Errorf("Error while setting address")
	}
	status, err := vc.c.Sys().SealStatus()
	if err != nil {
		return false, err
	}
	return status.Sealed, nil
}

func (vc *VaultClient) GetPodIP(podName, namespace string) (string, error) {
	k8s, err := NewK8SClient(vc.log)
	if err != nil {
		return "", errors.WithMessage(err, "error initializing k8s client")
	}

	// Get the pod's IP address
	pod, err := k8s.client.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	vc.log.Debug("Pod ip", pod.Status.PodIP)
	return pod.Status.PodIP, nil
}
