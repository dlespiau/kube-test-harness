package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/dlespiau/kube-test-harness"
	"github.com/dlespiau/kube-test-harness/logger"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func run(path string, args ...string) {
	cmd := exec.Command(path, args...)

	var wg sync.WaitGroup

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal("minikube:", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal("minikube:", err)
	}

	wg.Add(2)
	go func() {
		io.Copy(os.Stdout, stdout)
		wg.Done()
	}()
	go func() {
		io.Copy(os.Stderr, stderr)
		wg.Done()
	}()

	if err := cmd.Run(); err != nil {
		log.Fatal("minikube:", err)
	}
	wg.Wait()
}

func startMinikube() {
	fmt.Println("=== Starting minikube")
	run("minikube", "start")
	fmt.Println("=== Setting up kubernetes context")
	run("minikube", "update-context")
}

func main() {
	startMinikube()

	h := harness.New(harness.Options{
		LogLevel: logger.Debug,
	})
	if err := h.Setup(); err != nil {
		log.Fatal(err)
	}
	kube := h.NewTest(nil)

	kube.WaitForNodesReady(1, 3*time.Minute)

	for _, node := range kube.ListNodes(metav1.ListOptions{}).Items {
		fmt.Println("node:", node.Name, node.Status.Addresses)
	}
}
