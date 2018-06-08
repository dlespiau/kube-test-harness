package harness

import (
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dlespiau/balance/e2e/harness/logger"

	"golang.org/x/sync/errgroup"
)

type finalizer func() error

// Test is a single test running in a kubernetes cluster.
type Test struct {
	// ID is a unique identifier for the test, defined from the test function name.
	ID string
	// Namespace is name of the namespace automatically crated by Setup for the
	// test to run in.
	Namespace string

	nextObjectID uint64
	harness      *Harness
	t            *testing.T
	logger       logger.Logger
	inError      bool
	cleanUpFns   []finalizer
}

// NewTest creates a new test. Call Close() to free kubernetes resources
// allocated during the test.
func (h *Harness) NewTest(t *testing.T) *Test {
	// TestCtx is used among others for namespace names where '/' is forbidden
	prefix := strings.TrimPrefix(
		strings.Replace(
			t.Name(),
			"/",
			"-",
			-1,
		),
		"Test",
	)

	id := toSnake(prefix) + "-" + strconv.FormatInt(time.Now().Unix(), 10)
	test := &Test{
		ID:      id,
		harness: h,
		t:       t,
		logger:  h.options.Logger.ForTest(t),
	}
	test.Namespace = test.getObjID("ns")

	test.Infof("using API server %s", h.apiServer)

	return test
}

// getObjID returns an unique ID that can be used to name kubernetes objects. We
// also encode the object type in the name.
func (t *Test) getObjID(objectType string) string {
	id := atomic.AddUint64(&t.nextObjectID, 1)
	return t.ID + "-" + objectType + "-" + fmt.Sprintf("%d", id)
}

// Setup setups the test to be run in the Test.Namespace temporary namespace.
func (t *Test) Setup() *Test {
	t.CreateNamespace(t.Namespace)
	return t
}

// Close frees all kubernetes resources allocated during the test.
func (t *Test) Close() {
	// We're being called while panicking, don't cleanup!
	if r := recover(); r != nil {
		// XXX: Display more information about the test namespace and events.
		panic(r)
	}
	if t.t.Failed() || t.inError {
		// XXX: Display more information about the test namespace and events.
		return
	}

	if t.harness.options.NoCleanup {
		return
	}

	var eg errgroup.Group

	for i := len(t.cleanUpFns) - 1; i >= 0; i-- {
		eg.Go(t.cleanUpFns[i])
	}

	if err := eg.Wait(); err != nil {
		t.t.Fatal(err)
	}
}

func (t *Test) err(err error) {
	if err != nil {
		t.inError = true
		t.t.Fatal(err)
	}
}

func (t *Test) addFinalizer(fn finalizer) {
	t.cleanUpFns = append(t.cleanUpFns, fn)
}

// Debug prints a debug message.
func (t *Test) Debug(msg string) {
	t.t.Helper()
	t.logger.Logf(logger.Debug, msg)
}

// Debugf prints a debug message with a format string.
func (t *Test) Debugf(f string, args ...interface{}) {
	t.t.Helper()
	t.logger.Logf(logger.Debug, f, args...)
}

// Info prints an informational message.
func (t *Test) Info(msg string) {
	t.t.Helper()
	t.logger.Log(logger.Info, msg)
}

// Infof prints a informational message with a format string.
func (t *Test) Infof(f string, args ...interface{}) {
	t.t.Helper()
	t.logger.Logf(logger.Info, f, args...)
}
