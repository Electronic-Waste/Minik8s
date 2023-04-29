package core

import (
	"errors"
	"net"
	"strconv"
	"strings"
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

	// !!!add functional function step by step(such as volume and network and so on .......)
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
