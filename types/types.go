package types

import "fmt"

// Arguments required to create a service
type ServiceCreateSpec struct {
	BaseImage string
	Replicas  uint64
	Commands  []string
}

// Arguments required to kill a service
type ServiceRemoveSpec struct {
	ServiceID string
}

// Arguments required to scale a service
type ServiceScaleSpec struct {
	ServiceID string
	Replicas  uint64
}

func (spec *ServiceCreateSpec) String() string {
	result := ""
	result += fmt.Sprintf("\n  Base Image: %s", spec.BaseImage)
	result += fmt.Sprintf("\n  Replicas: %d", spec.Replicas)
	result += "\n  Commands: "
	for _, command := range spec.Commands {
		result += command + " "
	}
	return result + "\n"
}

func (spec *ServiceRemoveSpec) String() string {
	return fmt.Sprintf("\n  Service ID: %s\n", spec.ServiceID)
}

func (spec *ServiceScaleSpec) String() string {
	result := ""
	result += fmt.Sprintf("\n  Service ID: %s", spec.ServiceID)
	result += fmt.Sprintf("\n  Replicas: %d\n", spec.Replicas)
	return result
}
