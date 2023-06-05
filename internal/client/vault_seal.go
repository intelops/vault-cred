package client

import (
	"context"
	"fmt"
	"log"

	//	"time"

	//"time"

	//	"time"

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

func (v *VaultClient) MountSecrets(unsealKeys []string, rootTokens []string) error {
	// Initialize and configure your Vault client

	// Authenticate with Vault using the root token
	if len(rootTokens) == 0 {
		return fmt.Errorf("root token not provided")
	}

	v.c.SetToken(rootTokens[0]) // Use the first root token in the list, adjust if necessary

	// Mount the unseal key and root token in the desired Vault path

	for _, key := range unsealKeys {
		_, err := v.c.Logical().Write("/etc/vault/unsealvault", map[string]interface{}{
			"value": key,
		})
		if err != nil {
			return fmt.Errorf("failed to mount unseal key in Vault: %v", err)
		}
	}

	for _, token := range rootTokens {
		_, err := v.c.Logical().Write("/etc/vault/root-token", map[string]interface{}{
			"value": token,
		})
		if err != nil {
			return fmt.Errorf("failed to mount root token in Vault: %v", err)
		}
	}

	return nil
}
func (v *VaultClient) RetrieveKeys(nameSpace string, SecretName string) ([]string, []string, error) {

	var values []string
	var rootToken []string
	clientset := Config()
	namespace := nameSpace // Namespace where you want to create the Secret
	secretName := SecretName

	secret2, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		fmt.Println("Error while getting secret", err)
	}

	for key, value := range secret2.Data {

		if key == "roottoken" {
			rootToken = append(rootToken, string(value))

			continue // Skip the last element
		}

		//fmt.Printf("Retrieved value for  %s: %s\n", key, value)
		keys := string(value)
		values = append(values, keys)
		//	fmt.Println("Key is ", keys)

	}

	if (secret2.Name != "") && (secret2.Namespace != "") {
		fmt.Printf("Secret '%s' found in namespace '%s'\n", secret2.Name, secret2.Namespace)
	} else {
		log.Fatal("Given Namespace and Secret Name not found")
	}
	// Use the secret as needed
	for _, key := range rootToken {
		fmt.Println("Root Token", key)
	}
	// Mount the unseal key and root token to the Vault path
	err = v.MountSecrets(values, rootToken) // Replace with the actual function to mount the secrets in Vault
	if err != nil {
		return nil, nil, fmt.Errorf("failed to mount secrets to Vault: %v", err)

	}
	unsealkeys, root_token, err := v.ReadSecrets()
	if err != nil {
		log.Fatalf("Error while reading secrets from vault  %v", err)
	}

	return unsealkeys, root_token, nil
}

func (v *VaultClient) ReadSecrets() ([]string, []string, error) {
	var unsealKeys []string
	var rootTokens []string

	// Read the unseal key from Vault
	unsealSecret, err := v.c.Logical().Read("/etc/vault/unsealvault")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read unseal key from Vault: %v", err)
	}

	if unsealSecret != nil {
		key, ok := unsealSecret.Data["value"]
		if !ok {
			return nil, nil, fmt.Errorf("unseal key not found in Vault response")
		}
		unsealKeys = append(unsealKeys, fmt.Sprintf("%s", key))
	}

	// Read the root token from Vault
	rootTokenSecret, err := v.c.Logical().Read("/etc/vault/root-token")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read root token from Vault: %v", err)
	}

	if rootTokenSecret != nil {
		token, ok := rootTokenSecret.Data["value"]
		if !ok {
			return nil, nil, fmt.Errorf("root token not found in Vault response")
		}
		rootTokens = append(rootTokens, fmt.Sprintf("%s", token))
	}

	return unsealKeys, rootTokens, nil
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
	keys, _, err := v.ReadSecrets()
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
