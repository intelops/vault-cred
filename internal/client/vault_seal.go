package client

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/intelops/vault-cred/config"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	//	"context"
	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"

	//"fmt"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func (v VaultClient) GenerateUnsealKeysFromVaultClient() ([]string, string, error) {
	res := &api.InitRequest{
		SecretThreshold: 2,
		SecretShares:    3,
	}
	unsealKeys := []string{}

	key, err := v.c.Sys().Init(res)

	if err != nil {
		log.Fatal("Error while initializing ", err)
	}
	for _, key := range key.Keys {
		log.Fatal("Key is ", key)

		unsealKeys = append(unsealKeys, key)
	}

	rootToken := key.RootToken

	return unsealKeys, rootToken, err
}

func Config() *kubernetes.Clientset {
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal("Error whilr loading kubeconfig", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("Error while creating clientset", err)
	}
	return clientset
}
func (v VaultClient) Storekeys(nameSpace string, SecretName string) []string {
	clientset := Config()
	var values []string
	//var root_token []string
	namespace := nameSpace
	secretName := SecretName

	unsealKeys, rootToken, err := v.GenerateUnsealKeysFromVaultClient()
	if err != nil {
		log.Fatalf("Error while generating unseal keys %v", err)
	}

	stringData := make(map[string]string)
	for i, value := range unsealKeys {
		key := fmt.Sprintf("key%d", i+1)
		stringData[key] = value
		values = append(values, value)
	}

	key := "roottoken"
	stringData[key] = rootToken

	values = append(values, rootToken)
	newSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
		},
		StringData: stringData,
	}
	createdSecret, err := clientset.CoreV1().Secrets(namespace).Create(context.TODO(), newSecret, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("Failed to create secret: %v\n", err)
	}

	fmt.Printf("Secret '%s' created in namespace '%s'\n", createdSecret.Name, createdSecret.Namespace)

	return values
}

//

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
				return nil, errors.WithMessage(err, "error in reading unseal key file: "+file.Name())
			}
			keys = append(keys, key)
		}
	}

	return keys, nil
}

func (v VaultClient) IsVaultSealed() (bool, error) {
	status, err := v.c.Sys().SealStatus()
	if err != nil {
		log.Fatal("Errror while checking seal status", err)
		return false, err
	}

	return status.Sealed, nil

}
func (v VaultClient) Unseal() error {

	//var res *api.SealStatusResponse

	status, err := v.c.Sys().SealStatus()
	if err != nil {
		log.Fatalf("Error while checking seal status %v", err)
	}
	if !status.Initialized {

		_ = v.Storekeys("default", "vault-server")
	}
	//keys, err := readFileContent(config.VaultEnv{}.VaultUnSealKeyPath)
	if err != nil {
		return errors.WithMessage(err, "error in reading token file")
	}
	keys, err := readUnsealKeysFromPath(config.VaultEnv{}.VaultUnSealKeyPath)
	if err != nil {
		return fmt.Errorf("error while retrieving keys %v", err)
	}
	for _, key := range keys {

		_, err := v.c.Sys().Unseal(key)

		if err != nil {
			//flag = false
			return fmt.Errorf("error while unsealing  %v", err)
		}

	}

	return nil

}
