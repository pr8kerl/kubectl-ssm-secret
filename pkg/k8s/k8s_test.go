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

func TestK8sUpdateSecret(t *testing.T) {

	fakeClient := fake.NewSimpleClientset()
	k := &K8sClient{
		client:    fakeClient,
		namespace: "test",
	}
	t.Run("test UpdateSecret fails when secrets not exists", func(t *testing.T) {
		err := k.UpdateSecret("test", mockSecretData())
		assert.NotNil(t, err)
		assert.True(t, kerr.IsNotFound(err))
	})
	t.Run("test UpdateSecret succeeds", func(t *testing.T) {
		err := k.CreateSecret("test", mockSecretData(), false)
		assert.Nil(t, err)
		err = k.UpdateSecret("test", mockSecretData())
		assert.Nil(t, err)
	})

}

func TestK8sGetSecret(t *testing.T) {

	fakeClient := fake.NewSimpleClientset()
	k := &K8sClient{
		client:    fakeClient,
		namespace: "test",
	}
	wanted := mockSecretData()
	t.Run("test GetSecret returns expected results", func(t *testing.T) {
		err := k.CreateSecret("test", wanted, false)
		assert.Nil(t, err)
		secret, err := k.GetSecret("test")
		assert.Nil(t, err)
		assert.Equal(t, wanted, secret)
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
