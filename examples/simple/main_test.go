package simple

import (
	"os"
	"testing"

	"github.com/dlespiau/kube-test-harness"
)

var kube *harness.Harness

func TestMain(m *testing.M) {
	kube = harness.New(harness.Options{})
	os.Exit(kube.Run(m))
}
