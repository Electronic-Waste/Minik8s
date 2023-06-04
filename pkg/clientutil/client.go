package clientutil

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	digest "github.com/opencontainers/go-digest"
	"golang.org/x/sys/unix"
)

// this function is only useful in the unix system
func IsSocketAccessible(s string) error {
	abs, err := filepath.Abs(s)
	if err != nil {
		return err
	}
	// set AT_EACCESS to allow running nerdctl as a setuid binary
	return unix.Faccessat(-1, abs, unix.R_OK|unix.W_OK, unix.AT_EACCESS)
	return nil
}

func NewClient(ctx context.Context, namespace, address string, opts ...containerd.ClientOpt) (*containerd.Client, context.Context, context.CancelFunc, error) {

	ctx = namespaces.WithNamespace(ctx, namespace)

	address = strings.TrimPrefix(address, "unix://")
	const dockerContainerdaddress = "/var/run/docker/containerd/containerd.sock"
	if err := IsSocketAccessible(address); err != nil {
		if IsSocketAccessible(dockerContainerdaddress) == nil {
			err = fmt.Errorf("cannot access containerd socket %q (hint: try running with `--address %s` to connect to Docker-managed containerd): %w", address, dockerContainerdaddress, err)
		} else {
			err = fmt.Errorf("cannot access containerd socket %q: %w", address, err)
		}
		return nil, nil, nil, err
	}
	client, err := containerd.New(address, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	return client, ctx, cancel, nil
}

func getAddrHash(addr string) (string, error) {
	const addrHashLen = 8

	if runtime.GOOS != "windows" {
		addr = strings.TrimPrefix(addr, "unix://")

		var err error
		addr, err = filepath.EvalSymlinks(addr)
		if err != nil {
			return "", err
		}
	}

	d := digest.SHA256.FromString(addr)
	h := d.Encoded()[0:addrHashLen]
	return h, nil
}

// DataStore returns a string like "/var/lib/nerdctl/1935db59".
// "1935db9" is from `$(echo -n "/run/containerd/containerd.sock" | sha256sum | cut -c1-8)`
// on Windows it will return "%PROGRAMFILES%/nerdctl/1935db59"
func DataStore(dataRoot, address string) (string, error) {
	if err := os.MkdirAll(dataRoot, 0700); err != nil {
		return "", err
	}
	addrHash, err := getAddrHash(address)
	if err != nil {
		return "", err
	}
	dataStore := filepath.Join(dataRoot, addrHash)
	if err := os.MkdirAll(dataStore, 0700); err != nil {
		return "", err
	}
	return dataStore, nil
}
