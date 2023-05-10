package nettype

import (
	"fmt"
	"strings"
)

type Type int

const (
	Invalid Type = iota
	CNI
	Container
)

var netTypeToName = map[interface{}]string{
	Invalid:   "invalid",
	CNI:       "cni",
	Container: "container",
}

func Detect(names []string) (Type, error) {
	var res Type

	for _, name := range names {
		var tmp Type

		// In case of using --network=container:<container> to share the network namespace
		networkName := strings.SplitN(name, ":", 2)[0]
		switch networkName {
		case "container":
			tmp = Container
		default:
			tmp = CNI
		}
		if res != Invalid && res != tmp {
			return Invalid, fmt.Errorf("mixed network types: %v and %v", netTypeToName[res], netTypeToName[tmp])
		}
		res = tmp
	}

	// defaults to CNI
	if res == Invalid {
		res = CNI
	}

	return res, nil
}
