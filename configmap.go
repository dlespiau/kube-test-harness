package harness

import (
	"fmt"

	"github.com/pkg/errors"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func (test *Test) createConfigMap(namespace string, cm *v1.ConfigMap) error {
	cm.Namespace = namespace
	if _, err := test.harness.kubeClient.CoreV1().ConfigMaps(namespace).Create(cm); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to create ConfigMap %v ", cm.Name))
	}
	return nil
}

// CreateConfigMap creates a ConfigMap in the given namespace.
func (test *Test) CreateConfigMap(namespace string, cm *v1.ConfigMap) {
	err := test.createConfigMap(namespace, cm)
	test.err(err)
}

func (test *Test) loadConfigMap(manifestPath string) (*v1.ConfigMap, error) {
	manifest, err := test.harness.openManifest(manifestPath)
	if err != nil {
		return nil, err
	}
	dep := v1.ConfigMap{}
	if err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&dep); err != nil {
		return nil, errors.Wrapf(err, "failed to decode ConfigMap %s", manifestPath)
	}

	return &dep, nil
}

// LoadConfigMap loads a ConfigMap from a YAML manifest. The path to the
// manifest is relative to Harness.ManifestsDirectory.
func (test *Test) LoadConfigMap(manifestPath string) *v1.ConfigMap {
	dep, err := test.loadConfigMap(manifestPath)
	test.err(err)
	return dep
}

func (test *Test) createConfigMapFromFile(namespace string, manifestPath string) (*v1.ConfigMap, error) {
	s, err := test.loadConfigMap(manifestPath)
	if err != nil {
		return nil, err
	}
	err = test.createConfigMap(test.Namespace, s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// CreateConfigMapFromFile creates a ConfigMap from a manifest file in the given namespace.
func (test *Test) CreateConfigMapFromFile(namespace string, manifestPath string) *v1.ConfigMap {
	d, err := test.createConfigMapFromFile(namespace, manifestPath)
	test.err(err)
	return d
}

func (test *Test) deleteConfigMap(ConfigMap *v1.ConfigMap) error {
	if err := test.harness.kubeClient.CoreV1().ConfigMaps(ConfigMap.Namespace).Delete(ConfigMap.Name, nil); err != nil {
		return errors.Wrap(err, fmt.Sprintf("deleting ConfigMap %v failed", ConfigMap.Name))
	}
	return nil
}

// DeleteConfigMap deletes a ConfigMap.
func (test *Test) DeleteConfigMap(ConfigMap *v1.ConfigMap) {
	err := test.deleteConfigMap(ConfigMap)
	test.err(err)
}
