package k8s

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestK8sCreateSecret(t *testing.T) {

	fakeClient := fake.NewSimpleClientset()
	k := &K8sClient{
		client:    fakeClient,
		namespace: "test",
	}
	t.Run("test CreateSecret returns expected results", func(t *testing.T) {
		err := k.CreateSecret("test", mockSecretData(), false)
		assert.Nil(t, err)
	})
	t.Run("test CreateSecret fails with alreadyExists", func(t *testing.T) {
		err := k.CreateSecret("test", mockSecretData(), false)
		assert.NotNil(t, err)
		assert.True(t, kerr.IsAlreadyExists(err))
	})

}

func mockSecretData() map[string]string {
	var secret = make(map[string]string)
	secret["foo"] = "bar"
	secret["secret"] = "squirrel"
	return secret
}

func mockSecret(secret map[string]string) *v1.Secret {
	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:            "test",
			Namespace:       "test",
			OwnerReferences: []metav1.OwnerReference{},
		},
		StringData: secret,
	}
}
