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

func (test *Test) createClusterRole(cr *rbacv1.ClusterRole) error {
	if _, err := test.harness.kubeClient.RbacV1().ClusterRoles().Create(context.TODO(), cr, metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("failed to create cluster role %s: %w", cr.Name, err)
	}
	return nil
}

// CreateClusterRole creates a cluster role.
func (test *Test) CreateClusterRole(cr *rbacv1.ClusterRole) {
	err := test.createClusterRole(cr)
	test.err(err)
}

func (test *Test) loadClusterRole(manifestPath string) (*rbacv1.ClusterRole, error) {
	manifest, err := test.harness.openManifest(manifestPath)
	if err != nil {
		return nil, err
	}
	dep := rbacv1.ClusterRole{}
	if err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&dep); err != nil {
		return nil, fmt.Errorf("failed to decode cluster role %s: %w", manifestPath, err)
	}
	return &dep, nil
}

// LoadClusterRole loads a cluster role from a YAML manifest. The path to the
// manifest is relative to Harness.ManifestDirectory.
func (test *Test) LoadClusterRole(manifestPath string) *rbacv1.ClusterRole {
	cr, err := test.loadClusterRole(manifestPath)
	test.err(err)
	return cr
}

func (test *Test) createClusterRoleFromFile(manifestPath string) (*rbacv1.ClusterRole, error) {
	cr, err := test.loadClusterRole(manifestPath)
	if err != nil {
		return nil, err
	}
	err = test.createClusterRole(cr)
	if err != nil {
		return nil, err
	}
	return cr, nil
}

// CreateClusterRoleFromFile creates a cluster role from a manifest file.
func (test *Test) CreateClusterRoleFromFile(manifestPath string) *rbacv1.ClusterRole {
	cr, err := test.createClusterRoleFromFile(manifestPath)
	test.err(err)
	return cr
}

func (test *Test) deleteClusterRole(cr *rbacv1.ClusterRole) error {
	if err := test.harness.kubeClient.RbacV1().ClusterRoles().Delete(context.TODO(), cr.Name, metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("deleting cluster role %s failed: %w", cr.Name, err)
	}
	return nil
}

// DeleteClusterRole deletes a cluster role.
func (test *Test) DeleteClusterRole(cr *rbacv1.ClusterRole) {
	err := test.deleteClusterRole(cr)
	test.err(err)
}

// GetClusterRole returns a ClusterRole object if it exists or error.
func (test *Test) GetClusterRole(name string) (*rbacv1.ClusterRole, error) {
	cr, err := test.harness.kubeClient.RbacV1().ClusterRoles().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return cr, nil
}

func (test *Test) waitForClusterRoleReady(name string, timeout time.Duration) error {
	return wait.Poll(time.Second, timeout, func() (bool, error) {
		_, err := test.GetClusterRole(name)
		if err != nil {
			return false, err
		}
		return true, nil
	})
}

// WaitForClusterRoleReady waits until ClusterRole is created, otherwise times out.
func (test *Test) WaitForClusterRoleReady(cr *rbacv1.ClusterRole, timeout time.Duration) {
	err := test.waitForClusterRoleReady(cr.Name, timeout)
	test.err(err)
}
