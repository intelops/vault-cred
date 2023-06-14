package client

import (
	"context"
	"strings"

	"github.com/intelops/go-common/logging"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type K8SClient struct {
	client *kubernetes.Clientset
	log    logging.Logger
}

func NewK8SClient(log logging.Logger) (*K8SClient, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &K8SClient{client: clientset, log: log}, nil
}
func (k *K8SClient) CreateOrUpdateSecret(ctx context.Context, secretName, namespace string, data map[string]string) error {
	secData := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
		},
		StringData: data,
	}

	_, err := k.client.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil && err.Error() == "secrets \""+secretName+"\" not found" {
		createdSecret, err := k.client.CoreV1().Secrets(namespace).Create(context.TODO(), secData, metav1.CreateOptions{})
		if err != nil {
			return errors.WithMessage(err, "error in creating vault secret")
		}
		k.log.Infof("Secret %s created in namespace %s", createdSecret.Name, createdSecret.Namespace)
	} else if (err != nil) && err.Error() != "secrets \""+secretName+"\" not found" {
		return errors.WithMessage(err, "error in getting vault secret")
	} else {
		updatedsecret, err := k.client.CoreV1().Secrets(namespace).Update(context.TODO(), secData, metav1.UpdateOptions{})
		if err != nil {
			return errors.WithMessage(err, "error in creating vault secret")
		}
		k.log.Infof("Secret %s updated in namespace %s", updatedsecret.Name, updatedsecret.Namespace)
	}

	return nil
}

func (k *K8SClient) GetSecret(ctx context.Context, secretName, namespace string) (map[string]string, error) {
	secData, err := k.client.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.WithMessage(err, "error in creating vault secret")
	}

	secretMap := make(map[string]string)
	for key, value := range secData.Data {
		val := string(value)
		secretMap[key] = val
	}

	k.log.Debugf("Secret %s fetched from namespace %s", secretName, namespace)
	return secretMap, nil

}

func (k *K8SClient) GetConfigMapsHasPrefix(ctx context.Context, prefix string) (map[string]map[string]string, error) {
	configMaps := []corev1.ConfigMap{}
	namespaces, err := k.client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, errors.WithMessage(err, "failed to list namespaces: %v")
	}

	for _, ns := range namespaces.Items {
		cmList, err := k.client.CoreV1().ConfigMaps(ns.Name).List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, errors.WithMessagef(err, "failed to list ConfigMaps in namespace %s", ns.Name)
		}
		configMaps = append(configMaps, cmList.Items...)
	}

	allConfigMapData := map[string]map[string]string{}
	for _, cm := range configMaps {
		if strings.HasPrefix(cm.Name, "vault-policy") {
			cmKey := cm.Namespace + ":" + cm.Name
			allConfigMapData[cmKey] = cm.Data
		}
	}
	return allConfigMapData, nil
}
