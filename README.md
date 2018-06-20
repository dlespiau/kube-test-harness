[![API Reference](https://godoc.org/github.com/dlespiau/kube-harness?status.svg)](http://godoc.org/github.com/dlespiau/kube-harness)

# Kubernetes Test Harness

This package implements a test harness for running integration or end to end tests in a Kubernetes cluster.

Features:

- Integrate with the native Go [testing](https://golang.org/pkg/testing/) package.
- Create Kubernetes objects such as Deployments, Services, Secrets, ConfigMaps from either manifest file or the client-go API.
- Full access to the client-go API to manipulate Kubernetes objects.
- Wait for various readiness conditions.
- Each test runs in its own namespace, allowing them to run in parallel.
- Display a detailed error state to help the developer debug failure cases with the pod status, events and logs of failing pods.
- Automatic error checking, no need to check `err` at every line!

## Example

[`example/simple`](https://github.com/dlespiau/kube-harness/tree/master/examples/simple) has a self-contained example that shows how to write a test with `kube-harness`:

```go
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
        data, err := test.PodProxyGet(&pod, "80", "/").DoRaw()
        assert.NoError(t, err)
        assert.Equal(t, index, string(data))
    }
}
```

To run `kube-harness` tests a Kubernetes cluster is needed. It's then a `go test` invocation away:

```console
$ go test -v ./examples/simple/`
=== RUN   TestDeployNginx
--- PASS: TestDeployNginx (3.08s)
    test.go:63: using API server https://192.168.99.116:8443
    namespace.go:12: creating namespace deploy-nginx-1529445457-ns-1
    deployment.go:16: creating deployment nginx
    deployment.go:75: waiting for deployment nginx to be ready
    namespace.go:42: deleting namespace deploy-nginx-1529445457-ns-1
PASS
ok      github.com/dlespiau/kube-harness/examples/simple    3.090s
```

## Error State

When a test fails, `kube-harness` will display the state of the cluster to help the developer debug the problem. As an example, I changed the Nginx manifest in [`example/simple`](https://github.com/dlespiau/kube-harness/tree/master/examples/simple) to have an invalid image name. Running the test displayed clues about what the problem was:

```console
$ go test -v ./examples/simple/
=== RUN   TestDeployNginx

=== pods, namespace=kube-system

NAME                         READY   STATUS
kube-addon-manager-kubecon   1/1     Ready
kube-dns-86f6f55dd5-t8kcd    1/1     Ready
kubernetes-dashboard-5k2mn   1/1     Ready
storage-provisioner          1/1     Ready

=== pods, namespace=deploy-nginx-1529447347-ns-1

NAME                     READY   STATUS
nginx-75f7677558-n8rtr   0/1     ImagePullBackOff
nginx-75f7677558-x4wdv   0/1     ImagePullBackOff

=== logs, pod=nginx-75f7677558-n8rtr, container=nginx

container "nginx" in pod "nginx-75f7677558-n8rtr" is waiting to start: trying and failing to pull image

=== logs, pod=nginx-75f7677558-x4wdv, container=nginx

container "nginx" in pod "nginx-75f7677558-x4wdv" is waiting to start: trying and failing to pull image

--- FAIL: TestDeployNginx (30.11s)
    test.go:63: using API server https://192.168.99.116:8443
    namespace.go:12: creating namespace deploy-nginx-1529447347-ns-1
    deployment.go:16: creating deployment nginx
    deployment.go:75: waiting for deployment nginx to be ready
    test.go:191: timed out waiting for the condition
FAIL
FAIL    github.com/dlespiau/kube-harness/examples/simple    30.121s
```
