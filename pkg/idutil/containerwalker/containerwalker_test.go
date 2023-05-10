package containerwalker

// test file
import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"testing"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/defaults"
	"github.com/containerd/containerd/namespaces"
	"minik8s.io/pkg/constant"
)

func TestNewWalker(t *testing.T) {
	// construct a new pause container
	client, err := containerd.New(defaults.DefaultAddress)
	if err != nil {
		t.Error(err)
	}

	// try to create the SandBox
	// use cmd to build a pause container
	// run cmd : nerdctl run -d  --name fake_k8s_pod_pause   registry.aliyuncs.com/google_containers/pause:3.9
	cmd := exec.Command("nerdctl", "run", "-d", "--name", "test", constant.SandBox_Image)
	fmt.Println("finish the init of cmd")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		t.Error(err)
	}

	// try to get the SandBox by container walker
	var foundContainer containerd.Container
	walker := &ContainerWalker{
		Client: client,
		OnFound: func(ctx context.Context, found Found) error {
			if found.MatchCount > 1 {
				return fmt.Errorf("container networking: multiple containers found with prefix: %s", "test")
			}
			foundContainer = found.Container
			return nil
		},
	}

	// !!! : need to specify the namespace of finding container
	ctx := namespaces.WithNamespace(context.Background(), "default")
	n, err := walker.Walk(ctx, "test")
	if err != nil {
		t.Error(err)
	}
	if n == 0 {
		t.Error(err)
	}

	fmt.Println("pass test")
	fmt.Println(foundContainer)
}
