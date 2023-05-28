package network

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/containerd/cgroups"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/containers"
	"github.com/containerd/containerd/defaults"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/containerd/containerd/pkg/netns"
	gocni "github.com/containerd/go-cni"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/sirupsen/logrus"
	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/clientutil"
	"minik8s.io/pkg/constant"
	"minik8s.io/pkg/idutil/containerwalker"
	"minik8s.io/pkg/network/nettype"
	"minik8s.io/pkg/resolvconf"
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

// types.NetworkOptionsManager is an interface for reading/setting networking
// options for containers based on the provided command flags.
// TODO : we only support container-network and cni mode (maybe we can support much more mode in the future)
type NetworkOptionsManager interface {
	// TODO(wjl) : maybe we need to add a vertify method to check the network

	// Returns the set of NetworkingOptions which should be set as labels on the container.
	//
	// These options can potentially differ from the actual networking options
	// that the NetworkOptionsManager was initially instantiated with.
	// E.g: in container networking mode, the label will be normalized to an ID:
	// `--net=container:myContainer` => `--net=container:<ID of myContainer>`.
	// InternalNetworkingOptionLabels(context.Context) (core.NetworkOptions, error)

	// Returns a slice of `oci.SpecOpts` and `containerd.NewContainerOpts` which represent
	// the network specs which need to be applied to the container with the given ID.
	ContainerNetworkingOpts(context.Context, string) ([]oci.SpecOpts, []containerd.NewContainerOpts, error)
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
		Debug:          false,
		DebugFull:      false,
		Address:        defaults.DefaultAddress,
		Namespace:      namespaces.Default,
		Snapshotter:    containerd.DefaultSnapshotter,
		CNIPath:        gocni.DefaultCNIDir,
		CNINetConfPath: gocni.DefaultNetDir,
		// this config is different in the nerdctl source code config
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

// TODO : add the logic to determine use which mode
func ConstructNetworkManager(options GlobalCommandOptions, networkOptions core.NetworkOptions) NetworkOptionsManager {
	netype, _ := nettype.Detect(networkOptions.NetworkSlice)
	var manager NetworkOptionsManager
	switch netype {
	case nettype.Container:
		manager = &containerNetworkManager{
			globalOptions: options,
			netOpts:       networkOptions,
		}
	case nettype.CNI:
		manager = &CniNetworkManager{
			globalOptions: options,
			netOpts:       networkOptions,
			netNs:         nil,
		}
	}

	return manager
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

	dataStore, err := clientutil.DataStore(m.globalOptions.DataRoot, m.globalOptions.Address)
	if err != nil {
		return nil, nil, err
	}

	stateDir, err := ContainerStateDirPath(m.globalOptions, dataStore, containerID)
	if err != nil {
		return nil, nil, err
	}

	// generate a new conf path to hold for the config message of dns
	resolvConfPath := filepath.Join(stateDir, "resolv.conf")
	if err := m.buildResolvConf(resolvConfPath); err != nil {
		return nil, nil, err
	}

	opts = append(opts, withCustomResolvConf(resolvConfPath), withCustomHosts(""), withCustomEtcHostname(""))

	return opts, cOpts, nil
}

// types.NetworkOptionsManager implementation for container networking settings.
type containerNetworkManager struct {
	globalOptions GlobalCommandOptions
	netOpts       core.NetworkOptions
}

// Returns the relevant paths of the `hostname`, `resolv.conf`, and `hosts` files
// in the datastore of the container with the given ID.
func (m *containerNetworkManager) getContainerNetworkFilePaths(containerID string) (string, string, string, error) {
	// TODO : fix the hard in the future
	dataStore, err := clientutil.DataStore(constant.NerdctlDataRoot, m.globalOptions.Address)
	if err != nil {
		return "", "", "", err
	}
	conStateDir, err := ContainerStateDirPath(m.globalOptions, dataStore, containerID)
	if err != nil {
		return "", "", "", err
	}

	hostnamePath := filepath.Join(conStateDir, "hostname")
	resolvConfPath := filepath.Join(conStateDir, "resolv.conf")
	// etcHostsPath := hostsstore.HostsPath(dataStore, m.globalOptions.Namespace, containerID)
	// TODO : not config etcHostsPath here (may cause bug here)
	etcHostsPath := ""

	return hostnamePath, resolvConfPath, etcHostsPath, nil
}

// ContainerNetNSPath returns the netns path of a container.
func ContainerNetNSPath(ctx context.Context, c containerd.Container) (string, error) {
	task, err := c.Task(ctx, nil)
	if err != nil {
		return "", err
	}
	status, err := task.Status(ctx)
	if err != nil {
		return "", err
	}
	if status.Status != containerd.Running {
		return "", fmt.Errorf("invalid target container: %s, should be running", c.ID())
	}
	return fmt.Sprintf("/proc/%d/ns/net", task.Pid()), nil
}

func (m *containerNetworkManager) ContainerNetworkingOpts(ctx context.Context, _ string) ([]oci.SpecOpts, []containerd.NewContainerOpts, error) {
	opts := []oci.SpecOpts{}
	cOpts := []containerd.NewContainerOpts{}
	fmt.Printf("can reach here and netSlice is %s\n", &m.netOpts.NetworkSlice[0])
	container, err := m.getNetworkingContainerForArgument(ctx, m.netOpts.NetworkSlice[0])
	if err != nil {
		return nil, nil, err
	}
	containerID := container.ID()
	fmt.Printf("get the container id is %d\n", containerID)
	s, err := container.Spec(ctx)
	if err != nil {
		return nil, nil, err
	}
	hostname := s.Hostname

	netNSPath, err := ContainerNetNSPath(ctx, container)
	if err != nil {
		return nil, nil, err
	}

	hostnamePath, resolvConfPath, _, err := m.getContainerNetworkFilePaths(containerID)
	if err != nil {
		return nil, nil, err
	}

	opts = append(opts,
		oci.WithLinuxNamespace(specs.LinuxNamespace{
			Type: specs.NetworkNamespace,
			Path: netNSPath,
		}),
		withCustomResolvConf(resolvConfPath),
		withCustomHosts(""),
		oci.WithHostname(hostname),
		withCustomEtcHostname(hostnamePath),
	)

	return opts, cOpts, nil
}

// Searches for and returns the networking container for the given network argument.
func (m *containerNetworkManager) getNetworkingContainerForArgument(ctx context.Context, containerNetArg string) (containerd.Container, error) {
	netItems := strings.Split(containerNetArg, ":")
	if len(netItems) < 2 {
		return nil, fmt.Errorf("container networking argument format must be 'container:<id|name>', got: %q", containerNetArg)
	}
	containerName := netItems[1]

	// Namespace is "default" (default value) and Address is "/run/containerd/containerd.sock" (default value)
	client, ctxt, cancel, err := clientutil.NewClient(ctx, m.globalOptions.Namespace, m.globalOptions.Address)
	if err != nil {
		return nil, err
	}
	defer cancel()

	var foundContainer containerd.Container
	walker := &containerwalker.ContainerWalker{
		Client: client,
		OnFound: func(ctx context.Context, found containerwalker.Found) error {
			if found.MatchCount > 1 {
				return fmt.Errorf("container networking: multiple containers found with prefix: %s", containerName)
			}
			foundContainer = found.Container
			return nil
		},
	}
	n, err := walker.Walk(ctxt, containerName)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, fmt.Errorf("container networking: could not find container: %s", containerName)
	}

	return foundContainer, nil
}
