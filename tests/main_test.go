package tests

import (
	"flag"
	"os"
	"testing"

	"github.com/dlespiau/kube-test-harness"
	"github.com/dlespiau/kube-test-harness/logger"
)

var kube *harness.Harness

func TestMain(m *testing.M) {
	kubeconfig := flag.String("k8s.kubeconfig", "", "kube config path, e.g. $HOME/.kube/config")
	noCleanup := flag.Bool("k8s.no-cleanup", false, "should test cleanup after themselves")
	verbose := flag.Bool("k8s.log.verbose", false, "turn on more verbose logging")
	interactive := flag.Bool("k8s.log.interactive", false, "print log messages as they happen")

	flag.Parse()

	options := harness.Options{
		Kubeconfig:        *kubeconfig,
		ManifestDirectory: "manifests",
		NoCleanup:         *noCleanup,
	}
	if *verbose {
		options.LogLevel = logger.Debug
	}
	if *interactive {
		options.Logger = &logger.PrintfLogger{}
	}

	kube = harness.New(options)
	os.Exit(kube.Run(m))
}
