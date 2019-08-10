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

func (c *K8sClient) SetNamespace(ns string) {
	if len(ns) > 0 {
		c.namespace = ns
	}
}

func (c *K8sClient) GetNamespace() string {
	return c.namespace
}

func (c *K8sClient) CreateSecret(secretname string, secrets map[string]string, tls bool) error {

	if len(secrets) == 0 {
		return fmt.Errorf(fmt.Sprintf("k8s.CreateSecret: no secrets provided."))
	}
	var stype v1.SecretType = "Opaque"
	if tls {
		stype = "kubernetes.io/tls"
	}
	_, err := c.client.CoreV1().Secrets(c.namespace).Create(&v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: secretname,
		},
		Type: stype,
		Data: secretStringToBytes(secrets),
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *K8sClient) UpdateSecret(secretname string, secrets map[string]string) error {

	if len(secrets) == 0 {
		return fmt.Errorf(fmt.Sprintf("k8s.UpdateSecret: no secrets provided."))
	}
	_, err := c.client.CoreV1().Secrets(c.namespace).Update(&v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: secretname,
		},
		StringData: secrets,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *K8sClient) GetSecret(secretname string) (map[string]string, error) {

	secret, err := c.client.CoreV1().Secrets(c.namespace).Get(secretname, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return secretDataToString(secret), nil
}

func secretDataToString(secret *v1.Secret) map[string]string {
	results := make(map[string]string)
	for k, v := range secret.Data {
		results[k] = string(v)
	}
	return results
}

func secretStringToBytes(secret map[string]string) map[string][]byte {
	results := make(map[string][]byte)
	for k, v := range secret {
		results[k] = []byte(v)
	}
	return results
}
