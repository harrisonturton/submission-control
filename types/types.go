package types

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
