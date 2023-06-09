package client

import (
	"context"
	"path/filepath"

	"github.com/intelops/go-common/logging"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type K8SClient struct {
	c   *kubernetes.Clientset
	log logging.Logger
}

func NewK8SClient(log logging.Logger) (*K8SClient, error) {
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &K8SClient{c: clientset, log: log}, nil
}

func (k *K8SClient) CreateOrUpdateSecret(ctx context.Context, namespace string, data *corev1.Secret) error {
	createdSecret, err := k.c.CoreV1().Secrets(namespace).Create(context.TODO(), data, metav1.CreateOptions{})
	if err != nil {
		return errors.WithMessage(err, "error in creating vault secret")
	}

	k.log.Infof("Secret %s created in namespace %s", createdSecret.Name, createdSecret.Namespace)
	return nil
}
