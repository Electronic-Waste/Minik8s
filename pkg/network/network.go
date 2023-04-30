package network

import (
	"context"
	"errors"
	"fmt"
	"github.com/containerd/cgroups"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/containers"
	"github.com/containerd/containerd/defaults"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/containerd/containerd/pkg/netns"
	gocni "github.com/containerd/go-cni"
	"github.com/opencontainers/go-digest"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/sirupsen/logrus"
	"io/fs"
	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/resolvconf"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func withCustomResolvConf(src string) func(context.Context, oci.Client, *containers.Container, *oci.Spec) error {
	return func(_ context.Context, _ oci.Client, _ *containers.Container, s *oci.Spec) error {
		s.Mounts = append(s.Mounts, specs.Mount{
			Destination: "/etc/resolv.conf",
			Type:        "bind",
			Source:      src,
			Options:     []string{"bind", "rprivate"}, // writable
		})
		return nil
	}
}

func withCustomEtcHostname(src string) func(context.Context, oci.Client, *containers.Container, *oci.Spec) error {
	return func(_ context.Context, _ oci.Client, _ *containers.Container, s *oci.Spec) error {
		s.Mounts = append(s.Mounts, specs.Mount{
			Destination: "/etc/hostname",
			Type:        "bind",
			Source:      "/etc/hostname",
			Options:     []string{"bind"}, // writable
		})
		return nil
	}
}

func withCustomHosts(src string) func(context.Context, oci.Client, *containers.Container, *oci.Spec) error {
	return func(_ context.Context, _ oci.Client, _ *containers.Container, s *oci.Spec) error {
		s.Mounts = append(s.Mounts, specs.Mount{
			Destination: "/etc/hosts",
			Type:        "bind",
			Source:      "/etc/hosts",
			Options:     []string{"bind"}, // writable
		})
		return nil
	}
}

type GlobalCommandOptions struct {
	Debug            bool     `toml:"debug"`
	DebugFull        bool     `toml:"debug_full"`
	Address          string   `toml:"address"`
	Namespace        string   `toml:"namespace"`
	Snapshotter      string   `toml:"snapshotter"`
	CNIPath          string   `toml:"cni_path"`
	CNINetConfPath   string   `toml:"cni_netconfpath"`
	DataRoot         string   `toml:"data_root"`
	CgroupManager    string   `toml:"cgroup_manager"`
	InsecureRegistry bool     `toml:"insecure_registry"`
	HostsDir         []string `toml:"hosts_dir"`
	Experimental     bool     `toml:"experimental"`
	HostGatewayIP    string   `toml:"host_gateway_ip"`
}

func HostGatewayIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func IsSystemdAvailable() bool {
	fi, err := os.Lstat("/run/systemd/system")
	if err != nil {
		return false
	}
	return fi.IsDir()
}

// CgroupManager defaults to:
// - "systemd"  on v2 (rootful & rootless)
// - "cgroupfs" on v1 rootful
// - "none"     on v1 rootless
func CgroupManager() string {
	if cgroups.Mode() == cgroups.Unified && IsSystemdAvailable() {
		return "systemd"
	}
	return "cgroupfs"
}

// New creates a default Config object statically,
// without interpolating CLI flags, env vars, and toml.
func New() *GlobalCommandOptions {
	return &GlobalCommandOptions{
		Debug:            false,
		DebugFull:        false,
		Address:          defaults.DefaultAddress,
		Namespace:        namespaces.Default,
		Snapshotter:      containerd.DefaultSnapshotter,
		CNIPath:          gocni.DefaultCNIDir,
		CNINetConfPath:   gocni.DefaultNetDir,
		DataRoot:         "/var/lib/minik8s",
		CgroupManager:    CgroupManager(),
		InsecureRegistry: false,
		HostsDir:         []string{"/etc/containerd/certs.d", "/etc/docker/certs.d"},
		Experimental:     true,
		HostGatewayIP:    HostGatewayIP(),
	}
}

// types.NetworkOptionsManager implementation for CNI networking settings.
// This is a more specialized and OS-dependendant networking model so this
// struct provides different implementations on different platforms.
type CniNetworkManager struct {
	globalOptions GlobalCommandOptions
	netOpts       core.NetworkOptions
	netNs         *netns.NetNS
}

func DefaultNetOpt() core.NetworkOptions {
	netOpts := core.NetworkOptions{}
	netOpts.NetworkSlice = []string{"nat"}
	netOpts.MACAddress = ""
	netOpts.IPAddress = ""
	netOpts.Hostname = ""
	netOpts.DNSServers = nil
	netOpts.DNSSearchDomains = nil
	netOpts.DNSResolvConfOptions = []string{}
	netOpts.AddHost = nil
	netOpts.UTSNamespace = ""
	return netOpts
}

func ConstructNetworkManager(options GlobalCommandOptions, networkOptions core.NetworkOptions) CniNetworkManager {
	return CniNetworkManager{
		options,
		networkOptions,
		nil,
	}
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

// Returns the path to the Nerdctl-managed state directory for the container with the given ID.
func ContainerStateDirPath(globalOptions GlobalCommandOptions, dataStore, id string) (string, error) {
	// may need a step to check the correctness of
	dataStore = filepath.Join(dataStore, "containers")
	if err := os.MkdirAll(dataStore, 0700); err != nil {
		return "", err
	}
	dataStore = filepath.Join(dataStore, globalOptions.Namespace)
	if err := os.MkdirAll(dataStore, 0700); err != nil {
		return "", err
	}
	dataStore = filepath.Join(dataStore, id)
	if err := os.MkdirAll(dataStore, 0700); err != nil {
		return "", err
	}
	return dataStore, nil
}

func (m *CniNetworkManager) buildResolvConf(resolvConfPath string) error {
	var err error
	slirp4Dns := []string{}

	var (
		nameServers   = m.netOpts.DNSServers
		searchDomains = m.netOpts.DNSSearchDomains
		dnsOptions    = m.netOpts.DNSResolvConfOptions
	)

	// Use host defaults if any DNS settings are missing:
	if len(nameServers) == 0 || len(searchDomains) == 0 || len(dnsOptions) == 0 {
		fmt.Println("len is zero")
		conf, err := resolvconf.Get()
		if err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				return err
			}
			// if resolvConf file does't exist, using default resolvers
			conf = &resolvconf.File{}
			logrus.WithError(err).Debugf("resolvConf file doesn't exist on host")
		}
		conf, err = resolvconf.FilterResolvDNS(conf.Content, true)
		if err != nil {
			return err
		}
		if len(nameServers) == 0 {
			nameServers = resolvconf.GetNameservers(conf.Content, resolvconf.IPv4)
		}
		if len(searchDomains) == 0 {
			searchDomains = resolvconf.GetSearchDomains(conf.Content)
		}
		if len(dnsOptions) == 0 {
			dnsOptions = resolvconf.GetOptions(conf.Content)
		}
	}

	_, err = resolvconf.Build(resolvConfPath, append(slirp4Dns, nameServers...), searchDomains, dnsOptions)
	return err
}

func (m *CniNetworkManager) ContainerNetworkingOpts(_ context.Context, containerID string) ([]oci.SpecOpts, []containerd.NewContainerOpts, error) {
	opts := []oci.SpecOpts{}
	cOpts := []containerd.NewContainerOpts{}

	dataStore, err := DataStore(m.globalOptions.DataRoot, m.globalOptions.Address)
	if err != nil {
		return nil, nil, err
	}

	stateDir, err := ContainerStateDirPath(m.globalOptions, dataStore, containerID)
	if err != nil {
		return nil, nil, err
	}

	resolvConfPath := filepath.Join(stateDir, "resolv.conf")
	if err := m.buildResolvConf(resolvConfPath); err != nil {
		return nil, nil, err
	}

	// the content of /etc/hosts is created in OCI Hook
	//etcHostsPath, err := hostsstore.AllocHostsFile(dataStore, m.globalOptions.Namespace, containerID)
	//if err != nil {
	//	return nil, nil, err
	//}
	//opts = append(opts, withCustomResolvConf(resolvConfPath), withCustomHosts(etcHostsPath))
	opts = append(opts, withCustomResolvConf(resolvConfPath), withCustomHosts(""), withCustomEtcHostname(""))
	//if m.netOpts.UTSNamespace != UtsNamespaceHost {
	//	// If no hostname is set, default to first 12 characters of the container ID.
	//	hostname := m.netOpts.Hostname
	//	if hostname == "" {
	//		hostname = containerID
	//		if len(hostname) > 12 {
	//			hostname = hostname[0:12]
	//		}
	//	}
	//	m.netOpts.Hostname = hostname
	//
	//	hostnameOpts, err := writeEtcHostnameForContainer(m.globalOptions, m.netOpts.Hostname, containerID)
	//	if err != nil {
	//		return nil, nil, err
	//	}
	//	if hostnameOpts != nil {
	//		opts = append(opts, hostnameOpts...)
	//	}
	//}

	return opts, cOpts, nil
}
