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

type VolumeMount struct {
	// This must match the Name of a Volume.
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`

	// Path within the container at which the volume should be mounted.  Must
	// not contain ':'.
	MountPath string `json:"mountPath" protobuf:"bytes,3,opt,name=mountPath" yaml:"mountPath"`
}

// ResourceName is the name identifying various resources in a ResourceList.
type ResourceName string

type Quantity string

// ResourceList is a set of (resource name, quantity) pairs.
type ResourceList map[ResourceName]Quantity

// Resource names must be not more than 63 characters, consisting of upper- or lower-case alphanumeric characters,
// with the -, _, and . characters allowed anywhere, except the first or last character.
// The default convention, matching that for annotations, is to use lower-case names, with dashes, rather than
// camel case, separating compound words.
// Fully-qualified resource typenames are constructed from a DNS-style subdomain, followed by a slash `/` and a name.
const (
	// CPU, in cores. (500m = .5 cores)
	ResourceCPU ResourceName = "cpu"
	// Memory, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	ResourceMemory ResourceName = "memory"
)

type ResourceRequirements struct {
	// Limits describes the maximum amount of compute resources allowed.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	// +optional
	Limits ResourceList `json:"limits,omitempty" protobuf:"bytes,1,rep,name=limits,casttype=ResourceList,castkey=ResourceName"`
	// Requests describes the minimum amount of compute resources required.
	// If Requests is omitted for a container, it defaults to Limits if that is explicitly specified,
	// otherwise to an implementation-defined value. Requests cannot exceed Limits.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	// +optional
	Requests ResourceList `json:"requests,omitempty" protobuf:"bytes,2,rep,name=requests,casttype=ResourceList,castkey=ResourceName"`
}

// Container represents a single container that is expected to be run on the host.
type Container struct {
	// Required: This must be a DNS_LABEL.  Each container in a pod must
	// have a unique name.
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// Required.
	Image string `json:"image,omitempty" protobuf:"bytes,2,opt,name=image"`
	// Optional: The container image's entrypoint is used if this is not provided; cannot be updated.
	// Variable references $(VAR_NAME) are expanded using the container's environment.  If a variable
	// cannot be resolved, the reference in the input string will be unchanged.  Double $$ are reduced
	// to a single $, which allows for escaping the $(VAR_NAME) syntax: i.e. "$$(VAR_NAME)" will
	// produce the string literal "$(VAR_NAME)".  Escaped references will never be expanded, regardless
	// of whether the variable exists or not.
	// +optional
	Command []string `json:"command,omitempty" protobuf:"bytes,3,rep,name=command"`
	// Optional: The container image's cmd is used if this is not provided; cannot be updated.
	// Variable references $(VAR_NAME) are expanded using the container's environment.  If a variable
	// cannot be resolved, the reference in the input string will be unchanged.  Double $$ are reduced
	// to a single $, which allows for escaping the $(VAR_NAME) syntax: i.e. "$$(VAR_NAME)" will
	// produce the string literal "$(VAR_NAME)".  Escaped references will never be expanded, regardless
	// of whether the variable exists or not.
	// +optional
	Args []string `json:"args,omitempty" protobuf:"bytes,4,rep,name=args"`
	// Optional: Defaults to the container runtime's default working directory.
	// +optional
	WorkingDir string `json:"workingDir,omitempty" protobuf:"bytes,5,opt,name=workingDir"`
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
	Ports []ContainerPort `json:"ports,omitempty" yaml:"ports"`
	// Pod volumes to mount into the container's filesystem.
	// Cannot be updated.
	// +optional
	// +patchMergeKey=mountPath
	// +patchStrategy=merge
	VolumeMounts []VolumeMount `json:"volumeMounts,omitempty" patchStrategy:"merge" patchMergeKey:"mountPath" protobuf:"bytes,9,rep,name=volumeMounts" yaml:"volumeMounts""`

	// Compute resource requirements.
	// +optional
	Resources ResourceRequirements `json:"resources,omitempty" protobuf:"bytes,8,opt,name=resources" yaml:"resources""`

	// ------------------------------- next are params that will not in the yaml config file ----------------------------------------//
	Mounts []Mount
	// TODO(wjl) : add functional function step by step(such as volume and network and so on .......)
}

type PodStatus struct {
	// +optional
	Phase PodPhase
}

type Volume struct {
	// each volume in the pod must have a unique name
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`

	HostPath string `json:"hostPath,omitempty" yaml:"hostPath,omitempty"`
}

type PodSpec struct {
	Volumes []Volume `json:"volumes,omitempty"`

	// not consider the init Container

	// List of containers belonging to the pod.
	Containers []Container

	RunningNode Node
}

// ensure a variable which can identify a Pod
type Pod struct {
	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`

	meta.ObjectMeta `json:"metadata" yaml:"metadata" mapstructure:"metadata"`

	Spec PodSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	Status PodStatus
}

type Node struct {
	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`

	MetaData meta.ObjectMeta `json:"metadata" yaml:"metadata" mapstructure:"metadata"`

	Spec NodeSpec `json:"spec" yaml:"spec"`
}

type NodeSpec struct {
	MasterIp string `json:"masterIp" yaml:"masterIp"`

	NodeIp string `json:"nodeIp" yaml:"nodeIp"`
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

type Deployment struct {
	Metadata meta.ObjectMeta
	Spec     DeploymentSpec
	Status   DeploymentStatus
}

type DeploymentSpec struct {
	Replicas int
	Template Pod
	Selector string //must match .spec.template.metadata.labels
	//strategy	DeploymentStrategy
}

type DeploymentStatus struct {
	//ObservedGeneration int
	AvailableReplicas int
	//for later use
	//UpdatedReplicas int
	//ReadyReplicas   int
}

// Service is a named abstraction of software service (for example, mysql) consisting of local port
// that the proxy listens on, and the selector that determines which pods will answer
// requests sent through the proxy.
type Service struct {
	// Service's name (can be omitted)
	meta.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" mapstructure:"metadata"`

	// Service's kind is Service
	Kind string `json:"kind" yaml:"kind"`

	// Spec defines the behavior of a service.
	Spec ServiceSpec `json:"spec" yaml:"spec"`
}

// ServicePort represents the port on which the service is exposed
type ServicePort struct {
	// Optional if only one ServicePort is defined on this service: The
	// name of this port within the service.  This must be a DNS_LABEL.
	// All ports within a ServiceSpec must have unique names.  This maps to
	// the 'Name' field in EndpointPort objects.
	Name string `json:"name" yaml:"name"`

	// The IP protocol for this port.  Supports "TCP", "UDP", and "SCTP".
	Protocol Protocol `json:"protocol,omitempty" yaml:"protocol,omitempty"`

	// The port that will be exposed on the service.
	Port int32 `json:"port,omitempty" yaml:"port,omitempty"`

	// Optional: The target port on pods selected by this service.  If this
	// is a string, it will be looked up as a named port in the target
	// Pod's container ports.  If this is not specified, the value
	// of the 'port' field is used (an identity map).
	// This field is ignored for services with clusterIP=None, and should be
	// omitted or set equal to the 'port' field.
	TargetPort int32 `json:"targetPort,omitempty" yaml:"targetPort,omitempty"`
}

// ServiceSpec describes the attributes that a user creates on a service
type ServiceSpec struct {
	// Type determines how the Service is exposed. Defaults to ClusterIP. Valid
	// options are ExternalName, ClusterIP, NodePort, and LoadBalancer.
	Type string `json:"type,omitempty" yaml:"type,omitempty"`

	// Required: The list of ports that are exposed by this service.
	Ports []ServicePort `json:"ports,omitempty" yaml:"ports,omitempty"`

	// Route service traffic to pods with label keys and values matching this
	// selector. If empty or not present, the service is assumed to have an
	// external process managing its endpoints, which Kubernetes will not
	// modify. Only applies to types ClusterIP, NodePort, and LoadBalancer.
	// Ignored if type is ExternalName.
	Selector map[string]string `json:"selector,omitempty" yaml:"selector,omitempty"`

	// ClusterIP is the IP address of the service and is usually assigned
	// randomly by the master. If an address is specified manually and is not in
	// use by others, it will be allocated to the service
	ClusterIP string `json:"clusterIP,omitempty" yaml:"clusterIP,omitempty"`
}

// KubeproxyServiceParam is received by kuebproxy, which is used for creating service
type KubeproxyServiceParam struct {
	// Service's name
	ServiceName string `json:"serviceName,omitempty"`

	// ClusterIP is the IP address of the service and is usually assigned
	// randomly by the master. If an address is specified manually and is not in
	// use by others, it will be allocated to the service
	ClusterIP string `json:"clusterIP,omitempty"`

	// ServicePort represents the port on which the service is exposed
	ServicePorts []ServicePort `json:"servicePorts,omitempty"`

	// Pods' names
	PodNames []string `json:"podNames,omitempty"`

	// Pods' IPs
	PodIPs []string `json:"podIPs,omitempty"`
}

type ScheduleParam struct {
	RunPod   Pod
	NodeList []Node
}
