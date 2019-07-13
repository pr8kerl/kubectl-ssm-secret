package k8s

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sClient struct {
	client    kubernetes.Interface
	namespace string
}

type K8sConfig struct {
	rest      *rest.Config
	namespace string
}

func NewK8sClient(client kubernetes.Interface, ns string) *K8sClient {
	return &K8sClient{
		client:    client,
		namespace: ns,
	}
}

func NewK8sClientFromConfig(config *K8sConfig) (*K8sClient, error) {
	// create the clientset
	clientset, err := kubernetes.NewForConfig(config.rest)
	if err != nil {
		return nil, err
	}
	return NewK8sClient(clientset, config.namespace), nil
}

func NewK8sConfig() (*K8sConfig, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, err
	}
	ns, _, err := kubeConfig.Namespace()
	if err != nil {
		return nil, err
	}

	return &K8sConfig{
		rest:      config,
		namespace: ns,
	}, nil
}

func (c *K8sClient) CreateSecret(secretname string, secrets map[string]string) error {

	if len(secrets) == 0 {
		return fmt.Errorf(fmt.Sprintf("k8s.CreateSecret: no secrets provided."))
	}
	var stype v1.SecretType = "Opaque"
	_, err := c.client.CoreV1().Secrets(c.namespace).Create(&v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: secretname,
		},
		Type:       stype,
		StringData: secrets,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *K8sClient) CreateTlsSecret(secretname string, secrets map[string]string) error {

	if len(secrets) == 0 {
		return fmt.Errorf(fmt.Sprintf("k8s.CreateTlsSecret: no secrets provided."))
	}
	var stype v1.SecretType = "kubernetes.io/tls"
	_, err := c.client.CoreV1().Secrets(c.namespace).Create(&v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: secretname,
		},
		Type:       stype,
		StringData: secrets,
	})
	if err != nil {
		return err
	}
	return nil
}
