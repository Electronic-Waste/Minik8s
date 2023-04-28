package core

import "strings"

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
