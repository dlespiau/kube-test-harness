package harness

import (
	"os"
	"path/filepath"

	"github.com/dlespiau/kube-harness/logger"

	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func homeDirectory() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

// kubeconfigPath returns the kubeconfig location.
func kubeconfigPath() string {
	if env := os.Getenv("KUBECONFIG"); env != "" {
		return env
	}

	home := homeDirectory()
	if home == "" {
		return ""
	}
	return filepath.Join(home, ".kube", "config")
}

// newClientConfig returns a configuration object that can be used to configure
// a client in order to contact an API server with.
func newClientConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig == "" {
		kubeconfig = kubeconfigPath()
	}

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&clientcmd.ConfigOverrides{},
	).ClientConfig()
}

// Options are configuration options for the test harness.
type Options struct {
	// Kubeconfig is the path to a kubeconfig file. If not given, Harness will
	// honour the KUBECONFIG environment variable and try to use
	// $HOME/.kube/config.
	Kubeconfig string
	// ManifestDirectory is the root directory where the test Kubernetes manifests
	// are located. If not given, defaults to the current working directory.
	ManifestDirectory string
	// NoCleanup controls if tests should cleanup after them.
	NoCleanup bool
	// Logger is the Logger used to dispay test logs. If not given, Harness will
	// use logger.TestLogger which uses the logging built in the testing package.
	// This logger will only display logs on error or when -v is given to go test.
	// Additionally it will only dump the logs once the test has finished running.
	//
	// For writing and debugging long running tests, it is useful to have logs
	// being printed on stdout as it happens. logger.PrintfLogger can be used when
	// such behavior is needed.
	Logger logger.Logger
	// LogLevel controls how verbose the test logs are. Currently only Debug and
	// Info are available. If not given, defaults to Info.
	LogLevel logger.LogLevel
}

// Harness is a test harness for running integration tests on a kubernetes cluster.
type Harness struct {
	options    Options
	kubeClient kubernetes.Interface
	apiServer  string
}

// New creates a new test harness.
func New(options Options) *Harness {
	return &Harness{
		options: options,
	}
}

// Setup initializes the test harness.
func (h *Harness) Setup() error {
	// Logging
	if h.options.Logger == nil {
		h.options.Logger = &logger.TestLogger{}
	}
	if h.options.LogLevel == 0 {
		h.options.LogLevel = logger.Info
	}
	h.options.Logger.SetLevel(h.options.LogLevel)

	// Directories
	if h.options.ManifestDirectory == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return errors.Wrap(err, "harness")
		}
		h.options.ManifestDirectory = cwd
	}

	// Kubernetes client
	config, err := newClientConfig(h.options.Kubeconfig)
	if err != nil {
		return err
	}
	h.kubeClient, err = kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	h.apiServer = config.Host

	return nil
}

// Close terminates a test harness and frees its resources.
func (h *Harness) Close() error {
	return nil
}

func (h *Harness) openManifest(manifest string) (*os.File, error) {
	path := filepath.Join(h.options.ManifestDirectory, manifest)
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "open manifest")
	}

	return f, nil
}
