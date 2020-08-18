package harness

import (
	"context"
	"fmt"
	"time"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func (test *Test) createClusterRoleBinding(crb *rbacv1.ClusterRoleBinding) error {
	test.Debugf("creating cluster role binding %s", crb.Name)

	if _, err := test.harness.kubeClient.RbacV1().ClusterRoleBindings().Create(context.TODO(), crb, metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("failed to create cluster role binding %s: %w", crb.Name, err)
	}
	return nil
}

// CreateClusterRoleBinding creates a cluster role binding.
func (test *Test) CreateClusterRoleBinding(crb *rbacv1.ClusterRoleBinding) {
	err := test.createClusterRoleBinding(crb)
	test.err(err)
}

func (test *Test) loadClusterRoleBinding(manifestPath string) (*rbacv1.ClusterRoleBinding, error) {
	manifest, err := test.harness.openManifest(manifestPath)
	if err != nil {
		return nil, err
	}
	dep := rbacv1.ClusterRoleBinding{}
	if err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&dep); err != nil {
		return nil, fmt.Errorf("failed to decode cluster role binding %s: %w", manifestPath, err)
	}
	return &dep, nil
}

// LoadClusterRoleBinding loads a cluster role binding from a YAML manifest. The path to the
// manifest is relative to Harness.ManifestDirectory.
func (test *Test) LoadClusterRoleBinding(manifestPath string) *rbacv1.ClusterRoleBinding {
	crb, err := test.loadClusterRoleBinding(manifestPath)
	test.err(err)
	return crb
}

func (test *Test) createClusterRoleBindingFromFile(manifestPath string) (*rbacv1.ClusterRoleBinding, error) {
	crb, err := test.loadClusterRoleBinding(manifestPath)
	if err != nil {
		return nil, err
	}
	err = test.createClusterRoleBinding(crb)
	if err != nil {
		return nil, err
	}
	return crb, nil
}

// CreateClusterRoleBindingFromFile creates a cluster role binding from a manifest file.
func (test *Test) CreateClusterRoleBindingFromFile(manifestPath string) *rbacv1.ClusterRoleBinding {
	crb, err := test.createClusterRoleBindingFromFile(manifestPath)
	test.err(err)
	return crb
}

func (test *Test) deleteClusterRoleBinding(crb *rbacv1.ClusterRoleBinding) error {
	test.Debugf("deleting cluster role binding %s", crb.Name)

	if err := test.harness.kubeClient.RbacV1().ClusterRoleBindings().Delete(context.TODO(), crb.Name, metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("deleting cluster role binding %s failed: %w", crb.Name, err)
	}
	return nil
}

// DeleteClusterRoleBinding deletes a cluster role binding.
func (test *Test) DeleteClusterRoleBinding(crb *rbacv1.ClusterRoleBinding) {
	err := test.deleteClusterRoleBinding(crb)
	test.err(err)
}

// GetClusterRoleBinding returns a ClusterRoleBinding object if it exists or error.
func (test *Test) GetClusterRoleBinding(name string) (*rbacv1.ClusterRoleBinding, error) {
	crb, err := test.harness.kubeClient.RbacV1().ClusterRoleBindings().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return crb, nil
}

func (test *Test) waitForClusterRoleBindingReady(name string, timeout time.Duration) error {
	test.Debugf("waiting for cluster role binding %s to be ready", name)

	return wait.Poll(time.Second, timeout, func() (bool, error) {
		_, err := test.GetClusterRoleBinding(name)
		if err != nil {
			return false, err
		}
		return true, nil
	})
}

// WaitForClusterRoleBindingReady waits until ClusterRoleBinding is created, otherwise times out.
func (test *Test) WaitForClusterRoleBindingReady(crb *rbacv1.ClusterRole, timeout time.Duration) {
	err := test.waitForClusterRoleBindingReady(crb.Name, timeout)
	test.err(err)
}
