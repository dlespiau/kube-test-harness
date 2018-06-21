package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadDeleteDeployment(t *testing.T) {
	test := kube.NewTest(t).Setup()
	defer test.Close()

	dep := test.LoadDeployment("nginx-deployment.yaml")
	test.CreateDeployment(test.Namespace, dep)
	test.WaitForDeploymentReady(dep, 30*time.Second)

	// We should have the same number of pods in the cluster than the ones defined in the manifest
	pods := test.ListPodsFromDeployment(dep).Items
	assert.Equal(t, int(*dep.Spec.Replicas), len(pods))
	for _, pod := range pods {
		ready, err := test.PodReady(pod)
		assert.NoError(t, err)
		assert.True(t, ready)
	}

	test.DeleteDeployment(dep)
	test.WaitForDeploymentDeleted(dep, 30*time.Second)
}
