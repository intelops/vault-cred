package client

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/intelops/go-common/logging"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)
const (
	labelSelector = "app.kubernetes.io/name=vault"
)
type K8SClient struct {
	client *kubernetes.Clientset
	log    logging.Logger
}

type ConfigMapData struct {
	Name            string
	Namespace       string
	Data            map[string]string
	LastUpdatedTime time.Time
}

type SecretData struct {
	Data            map[string]string
	LastUpdatedTime time.Time
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

func (k *K8SClient) GetClusterConfig() (*rest.Config, error) {
	return rest.InClusterConfig()
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
	if err != nil && k8serrors.IsNotFound(err) {
		createdSecret, err := k.client.CoreV1().Secrets(namespace).Create(context.TODO(), secData, metav1.CreateOptions{})
		if err != nil {
			return errors.WithMessage(err, "error in creating vault secret")
		}
		k.log.Infof("Secret %s created in namespace %s", createdSecret.Name, createdSecret.Namespace)
		return nil
	}

	updatedsecret, err := k.client.CoreV1().Secrets(namespace).Update(context.TODO(), secData, metav1.UpdateOptions{})
	if err != nil {
		return errors.WithMessage(err, "error in creating vault secret")
	}
	k.log.Infof("Secret %s updated in namespace %s", updatedsecret.Name, updatedsecret.Namespace)
	return nil
}

func (k *K8SClient) GetSecret(ctx context.Context, secretName, namespace string) (*SecretData, error) {
	secData, err := k.client.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, errors.New("secret not found")
		}
		return nil, errors.WithMessage(err, "error in creating vault secret")
	}
	lastUpdatedTime, err := time.Parse(time.RFC3339, secData.ObjectMeta.CreationTimestamp.Format(time.RFC3339))
	if err != nil {
		return nil, errors.New("secret date is not valid")
	}

	secretMap := make(map[string]string)
	for key, value := range secData.Data {
		val := string(value)
		secretMap[key] = val
	}

	k.log.Debugf("Secret %s fetched from namespace %s", secretName, namespace)
	return &SecretData{Data: secretMap, LastUpdatedTime: lastUpdatedTime}, nil
}

func (k *K8SClient) GetConfigMapsHasPrefix(ctx context.Context, prefix string) ([]ConfigMapData, error) {
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

	allConfigMapData := []ConfigMapData{}
	for _, cm := range configMaps {
		if strings.HasPrefix(cm.Name, prefix) {
			lastUpdatedTime, err := time.Parse(time.RFC3339, cm.ObjectMeta.CreationTimestamp.Format(time.RFC3339))
			if err != nil {
				k.log.Debugf("Configmap %s doesn't has vaild time, skipping", cm.Name, cm.Namespace)
				continue
			}
			configMapData := ConfigMapData{
				Name:            cm.Name,
				Namespace:       cm.Namespace,
				Data:            cm.Data,
				LastUpdatedTime: lastUpdatedTime,
			}
			allConfigMapData = append(allConfigMapData, configMapData)
		}
	}
	return allConfigMapData, nil
}

func (k *K8SClient) GetVaultPodInstances(ctx context.Context) ([]string, error) {
	var podnames []string


	pods, err := k.client.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil,errors.WithMessage(err, "error while retrieving the pods ")
	}


	for _, pod := range pods.Items {
		podnames = append(podnames, pod.Name)
		

	}
	sort.Strings(podnames)
	
	return podnames,nil
}
