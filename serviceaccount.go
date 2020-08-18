package harness

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func (test *Test) createServiceAccount(namespace string, serviceAccount *v1.ServiceAccount) error {
	test.Debugf("creating serviceaccount %s", serviceAccount.Name)

	serviceAccount.Namespace = namespace
	if _, err := test.harness.kubeClient.CoreV1().ServiceAccounts(namespace).Create(context.TODO(), serviceAccount, metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("failed to create ServiceAccount %s: %w", serviceAccount.Name, err)
	}
	return nil
}

// CreateServiceAccount creates a service account in the given namespace.
func (test *Test) CreateServiceAccount(namespace string, serviceAccount *v1.ServiceAccount) {
	err := test.createServiceAccount(namespace, serviceAccount)
	test.err(err)
}

func (test *Test) loadServiceAccount(manifestPath string) (*v1.ServiceAccount, error) {
	manifest, err := test.harness.openManifest(manifestPath)
	if err != nil {
		return nil, err
	}
	dep := v1.ServiceAccount{}
	if err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&dep); err != nil {
		return nil, fmt.Errorf("failed to decode ServiceAccount %s: %w", manifestPath, err)
	}

	return &dep, nil
}

// LoadServiceAccount loads a service account from a YAML manifest. The path to the
// manifest is relative to Harness.ManifestDirectory.
func (test *Test) LoadServiceAccount(manifestPath string) *v1.ServiceAccount {
	sa, err := test.loadServiceAccount(manifestPath)
	test.err(err)
	return sa
}

func (test *Test) createServiceAccountFromFile(namespace string, manifestPath string) (*v1.ServiceAccount, error) {
	sa, err := test.loadServiceAccount(manifestPath)
	if err != nil {
		return nil, err
	}
	err = test.createServiceAccount(namespace, sa)
	if err != nil {
		return nil, err
	}
	return sa, nil
}

// CreateServiceAccountFromFile creates a service account from a manifest file in the given namespace.
func (test *Test) CreateServiceAccountFromFile(namespace string, manifestPath string) *v1.ServiceAccount {
	sa, err := test.createServiceAccountFromFile(namespace, manifestPath)
	test.err(err)
	return sa
}

func (test *Test) deleteServiceAccount(serviceAccount *v1.ServiceAccount) error {
	test.Debugf("deleting serviceaccount %s ", serviceAccount.Name)

	if err := test.harness.kubeClient.CoreV1().ServiceAccounts(serviceAccount.Namespace).Delete(context.TODO(), serviceAccount.Name, metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("deleting ServiceAccount %s failed: %w", serviceAccount.Name, err)
	}
	return nil
}

// DeleteServiceAccount deletes a ServiceAccount.
func (test *Test) DeleteServiceAccount(serviceAccount *v1.ServiceAccount) {
	err := test.deleteServiceAccount(serviceAccount)
	test.err(err)
}

// GetServiceAccount returns a ServiceAccount object if it exists or error.
func (test *Test) GetServiceAccount(namespace, name string) (*v1.ServiceAccount, error) {
	return test.harness.kubeClient.CoreV1().ServiceAccounts(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (test *Test) waitForServiceAccountReady(serviceAccount *v1.ServiceAccount) error {
	test.Debugf("waiting for serviceaccount %s to be ready", serviceAccount.Name)

	err := wait.Poll(time.Second, time.Minute*5, func() (bool, error) {
		_, err := test.GetServiceAccount(serviceAccount.Namespace, serviceAccount.Name)
		if err != nil {
			return false, err
		}
		return true, nil
	})
	return err
}

// WaitForServiceAccountReady waits until ConfigMap is created, otherwise times out.
func (test *Test) WaitForServiceAccountReady(serviceAccount *v1.ServiceAccount) {
	test.err(test.waitForServiceAccountReady(serviceAccount))
}
