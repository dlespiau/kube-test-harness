package harness

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (test *Test) createNamespace(name string) (*v1.Namespace, error) {
	test.Debugf("creating namespace %s", name)

	namespace, err := test.harness.kubeClient.CoreV1().Namespaces().Create(context.TODO(), &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create namespace with name %v: %w", name, err)
	}
	return namespace, nil
}

// CreateNamespace creates a new namespace.
func (test *Test) CreateNamespace(name string) {
	_, err := test.createNamespace(name)
	test.err(err)

	test.addNamespace(name)

	test.addFinalizer(func() error {
		if err := test.deleteNamespace(name); err != nil {
			return err
		}
		return nil
	})
}

func (test *Test) deleteNamespace(name string) error {
	test.Debugf("deleting namespace %s", name)

	test.removeNamespace(name)

	return test.harness.kubeClient.CoreV1().Namespaces().Delete(context.TODO(), name, metav1.DeleteOptions{})
}

// DeleteNamespace deletes a Namespace.
func (test *Test) DeleteNamespace(name string) {
	err := test.deleteNamespace(name)
	test.err(err)
}

// GetNamespace returns a Namespace object if it exists or error.
func (test *Test) GetNamespace(name string) (*v1.Namespace, error) {
	ns, err := test.harness.kubeClient.CoreV1().Namespaces().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return ns, nil
}
