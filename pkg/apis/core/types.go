package core

import (
	"errors"
	"net"
	"path/filepath"
	"strconv"
	"strings"

	gocni "github.com/containerd/go-cni"
	"minik8s.io/pkg/apis/meta"

	"github.com/opencontainers/runtime-spec/specs-go"
)

// Protocol defines network protocols supported for things like container ports.
// +enum
type Protocol string

const (
	// ProtocolTCP is the TCP protocol.
	ProtocolTCP Protocol = "TCP"
	// ProtocolUDP is the UDP protocol.
	ProtocolUDP Protocol = "UDP"
	// ProtocolSCTP is the SCTP protocol.
	ProtocolSCTP Protocol = "SCTP"
)

// PodPhase is a label for the condition of a pod at the current time.
type PodPhase string

// These are the valid statuses of pods.
const (
	// PodPending means the pod has been accepted by the system, but one or more of the containers
	// has not been started. This includes time before being bound to a node, as well as time spent
	// pulling images onto the host.
	PodPending PodPhase = "Pending"
	// PodRunning means the pod has been bound to a node and all of the containers have been started.
	// At least one container is still running or is in the process of being restarted.
	PodRunning PodPhase = "Running"
	// PodSucceeded means that all containers in the pod have voluntarily terminated
	// with a container exit code of 0, and the system is not going to restart any of these containers.
	PodSucceeded PodPhase = "Succeeded"
	// PodFailed means that all containers in the pod have terminated, and at least one container has
	// terminated in a failure (exited with a non-zero exit code or was stopped by the system).
	PodFailed PodPhase = "Failed"
	// PodUnknown means that for some reason the state of the pod could not be obtained, typically due
	// to an error in communicating with the host of the pod.
	// Deprecated in v1.21: It isn't being set since 2015 (74da3b14b0c0f658b3bb8d2def5094686d0e9095)
	PodUnknown PodPhase = "Unknown"
)

type ContainerPort struct {
	// If specified, this must be an IANA_SVC_NAME and unique within the pod. Each
	// named port in a pod must have a unique name. Name for the port that can be
	// referred to by services.
	// +optional
	Name string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
	// Number of port to expose on the host.
	// If specified, this must be a valid port number, 0 < x < 65536.
	// If HostNetwork is specified, this must match ContainerPort.
	// Most containers do not need this.
	// +optional
	HostPort int32 `json:"hostPort,omitempty" protobuf:"varint,2,opt,name=hostPort"`
	// Number of port to expose on the pod's IP address.
	// This must be a valid port number, 0 < x < 65536.
	ContainerPort int32 `json:"containerPort" protobuf:"varint,3,opt,name=containerPort"`
	// Protocol for port. Must be UDP, TCP, or SCTP.
	// Defaults to "TCP".
	// +optional
	// +default="TCP"
	Protocol Protocol `json:"protocol,omitempty" protobuf:"bytes,4,opt,name=protocol,casttype=Protocol"`
	// What host IP to bind the external port to.
	// +optional
	HostIP string `json:"hostIP,omitempty" protobuf:"bytes,5,opt,name=hostIP"`
}

type Mount struct {
	// host file path
	SourcePath string

	// container file path
	DestinationPath string
}

// Container represents a single container that is expected to be run on the host.
type Container struct {
	// Required: This must be a DNS_LABEL.  Each container in a pod must
	// have a unique name.
	Name string
	// Required.
	Image string
	// Optional: The container image's entrypoint is used if this is not provided; cannot be updated.
	// Variable references $(VAR_NAME) are expanded using the container's environment.  If a variable
	// cannot be resolved, the reference in the input string will be unchanged.  Double $$ are reduced
	// to a single $, which allows for escaping the $(VAR_NAME) syntax: i.e. "$$(VAR_NAME)" will
	// produce the string literal "$(VAR_NAME)".  Escaped references will never be expanded, regardless
	// of whether the variable exists or not.
	// +optional
	Command []string
	// Optional: The container image's cmd is used if this is not provided; cannot be updated.
	// Variable references $(VAR_NAME) are expanded using the container's environment.  If a variable
	// cannot be resolved, the reference in the input string will be unchanged.  Double $$ are reduced
	// to a single $, which allows for escaping the $(VAR_NAME) syntax: i.e. "$$(VAR_NAME)" will
	// produce the string literal "$(VAR_NAME)".  Escaped references will never be expanded, regardless
	// of whether the variable exists or not.
	// +optional
	Args []string
	// Optional: Defaults to the container runtime's default working directory.
	// +optional
	WorkingDir string
	// List of ports to expose from the container. Not specifying a port here
	// DOES NOT prevent that port from being exposed. Any port which is
	// listening on the default "0.0.0.0" address inside a container will be
	// accessible from the network.
	// Modifying this array with strategic merge patch may corrupt the data.
	// For more information See https://github.com/kubernetes/kubernetes/issues/108255.
	// Cannot be updated.
	// +optional
	// +patchMergeKey=containerPort
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=containerPort
	// +listMapKey=protocol
	// -p/--publish=127.0.0.1:80:8080/tcp ... but in nervctl version : only 127.0.0.1:80:8080/tcp
	Ports []ContainerPort `json:"ports,omitempty" patchStrategy:"merge" patchMergeKey:"containerPort" protobuf:"bytes,6,rep,name=ports"`

	Mounts []Mount
	// TODO(wjl) : add functional function step by step(such as volume and network and so on .......)
}

type PodStatus struct {
	// +optional
	Phase PodPhase
}

type HostPathVolumeSource struct {
	Path string
}

type VolumeSource struct {
	// only support host map at this time
	HostPath *HostPathVolumeSource

	// TODO : try to add emptyDir type of volume
}

type Volume struct {
	// each volume in the pod must have a unique name
	Name string

	VolumeSource
}

type PodSpec struct {
	Volumes []Volume

	// not consider the init Container

	// List of containers belonging to the pod.
	Containers []Container
}

// ensure a variable which can identify a Pod
type Pod struct {
	meta.ObjectMeta

	Kind string

	Spec PodSpec

	Status PodStatus
}

func (c *Container) String() string {
	str := ""
	str += "Name " + c.Name + "\n" +
		"Image " + c.Image + "\n" +
		"WorkingDir" + c.WorkingDir + "\n"
	str += "Command " + strings.Join(c.Command, " ") + "\n"
	str += "Args " + strings.Join(c.Args, " ") + "\n"
	return str
}

func ConstructPort(str string) (ContainerPort, error) {
	// default protocol
	var protocol Protocol = "TCP"
	// 127.0.0.1:80:8080/tcp
	// split to get protocol first
	strSlice := strings.Split(str, "/")
	if len(strSlice) > 1 {
		if len(strSlice) > 2 {
			return ContainerPort{}, errors.New("wrong format of port")
		}
		protocol = Protocol(strSlice[1])
	}

	// parse ip and hostport and containerport

	leftSlice := strings.Split(strSlice[0], ":")
	switch len(leftSlice) {
	case 1:
		{
			// we see it as a random port case
			return ContainerPort{}, errors.New("we have not develop random port assignment\n")
		}
	case 2:
		{
			// format : 80:8080
			return ContainerPort{
				HostPort: func(i int, _ error) int32 {
					return int32(i)
				}(strconv.Atoi(leftSlice[0])),
				ContainerPort: func(i int, _ error) int32 {
					return int32(i)
				}(strconv.Atoi(leftSlice[1])),
				Protocol: protocol,
				HostIP:   "0.0.0.0",
			}, nil
		}
	case 3:
		{
			if err := net.ParseIP(leftSlice[0]); err != nil {
				return ContainerPort{}, errors.New("ip format error")
			}
			// format : 127.0.0.1:80:8080
			return ContainerPort{
				HostIP: leftSlice[0],
				HostPort: func(i int, _ error) int32 {
					return int32(i)
				}(strconv.Atoi(leftSlice[1])),
				ContainerPort: func(i int, _ error) int32 {
					return int32(i)
				}(strconv.Atoi(leftSlice[2])),
				Protocol: protocol,
			}, nil
		}
	default:
		{
			return ContainerPort{}, errors.New("more num")
		}
	}
}

func ConstructPorts(str string) ([]ContainerPort, error) {
	// format : 127.0.0.1:80:8080/tcp,127.0.0.1:1000:8000/tcp
	// split by ','
	port_slice := strings.Split(str, ",")
	res := []ContainerPort{}
	for _, port_str := range port_slice {
		p, err := ConstructPort(port_str)
		if err != nil {
			return []ContainerPort{}, err
		}
		res = append(res, p)
	}
	return res, nil
}

func ConstructMount(str string) (Mount, error) {
	get_slice := strings.Split(str, ":")
	switch len(get_slice) {
	case 1:
		{
			if !filepath.IsAbs(get_slice[0]) {
				return Mount{}, errors.New("only support abs file path")
			}
			return Mount{
				DestinationPath: get_slice[0],
			}, nil
		}
	case 2:
		{
			if !filepath.IsAbs(get_slice[0]) || !filepath.IsAbs(get_slice[1]) {
				return Mount{}, errors.New("only support abs file path")
			}
			return Mount{
				SourcePath:      get_slice[0],
				DestinationPath: get_slice[1],
			}, nil
		}
	default:
		{
			return Mount{}, errors.New("error format of mount")
		}

	}
}

func ConstructMounts(str string) ([]Mount, error) {
	// we design the format of mount is mount1|mount2|mount3
	res := []Mount{}
	mount_strs := strings.Split(str, "|")
	for _, mount_str := range mount_strs {
		m, err := ConstructMount(mount_str)
		if err != nil {
			return []Mount{}, err
		}
		res = append(res, m)
	}
	return res, nil
}

func ConvertMount(mount Mount) specs.Mount {
	return specs.Mount{
		Source:      mount.SourcePath,
		Destination: mount.DestinationPath,
		Type:        "bind",
		Options:     []string{"bind"},
	}
}

func ConvertMounts(mounts []Mount) []specs.Mount {
	res := []specs.Mount{}
	for _, mount := range mounts {
		res = append(res, ConvertMount(mount))
	}
	return res
}

// Struct defining networking-related options.
type NetworkOptions struct {
	// --net/--network=<net name> ...
	NetworkSlice []string

	// --mac-address=<MAC>
	MACAddress string

	// --ip=<container static IP>
	IPAddress string

	// -h/--hostname=<container Hostname>
	Hostname string

	// --dns=<DNS host> ...
	DNSServers []string

	// --dns-opt/--dns-option=<resolv.conf line> ...
	DNSResolvConfOptions []string

	// --dns-search=<domain name> ...
	DNSSearchDomains []string

	// --add-host=<host:IP> ...
	AddHost []string

	// --uts=<Unix Time Sharing namespace>
	UTSNamespace string

	// -p/--publish=127.0.0.1:80:8080/tcp ...
	PortMappings []gocni.PortMapping
}