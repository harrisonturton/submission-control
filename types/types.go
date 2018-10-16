package types

// An RPC call to scale a service
type ScaleArgs struct {
	ServiceID string
	Replicas  uint64
}
