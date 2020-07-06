package simple

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const index = `<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
<style>
    body {
        width: 35em;
        margin: 0 auto;
        font-family: Tahoma, Verdana, Arial, sans-serif;
    }
</style>
</head>
<body>
<h1>Welcome to nginx!</h1>
<p>If you see this page, the nginx web server is successfully installed and
working. Further configuration is required.</p>

<p>For online documentation and support please refer to
<a href="http://nginx.org/">nginx.org</a>.<br/>
Commercial support is available at
<a href="http://nginx.com/">nginx.com</a>.</p>

<p><em>Thank you for using nginx.</em></p>
</body>
</html>
`

// TestDeployNginx deploys nginx on a kubernetes cluster and tests it's running
// correctly by checking the response of an HTTP GET request on /.
func TestDeployNginx(t *testing.T) {
	// Create a new test from the harness object. A test runs in one namespace, can
	// create Kubernetes objects and perform various checking operations.
	//
	// Always call Close when finished running the test. It will cleanup the
	// resources created in that test and take care of displaying the error state
	// if the test has failed.
	test := kube.NewTest(t).Setup()
	defer test.Close()

	// Create a Deployment from a manifest file.
	//
	// test.Namespace holds the namespace automatically created by the test harness
	// for this test to run in.
	d := test.CreateDeploymentFromFile(test.Namespace, "nginx-deployment.yaml")

	// Wait until the deployment is up and running, timeout of 30s.
	test.WaitForDeploymentReady(d, 30*time.Second)

	// For each pod of the Deployment, check we receive a sensible response to a
	// GET request on /.
	for _, pod := range test.ListPodsFromDeployment(d).Items {
		data, err := test.PodProxyGet(&pod, "80", "/").DoRaw(context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, index, string(data))
	}
}
