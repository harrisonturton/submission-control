package types

import "time"

// TestJob is the object that holds metadata
// about the test job we recieve.
type TestJob struct {
	Timestamp  time.Time
	TestConfig TestConfig
}

// TestConfig is the type that holds the default
// test configurations, and the test commands for
// each test.
type TestConfig struct {
	Version *string      `json:"version"`
	Env     *Environment `json:"environment"`
}

// Environment is the configuration for the testing
// environment in general, within the container
type Environment struct {
	Image *string            `json:"image"`
	Vars  *map[string]string `json:"vars"`
}

// Compare determines if the two configurations
// are the same.
func (a TestConfig) Compare(b TestConfig) bool {
	if a.Version != b.Version || *a.Version != *b.Version {
		return false
	}
	if (a.Env == nil) != (b.Env == nil) {
		return false
	}
	if a.Env == nil && b.Env == nil {
		return true
	}
	return a.Env.Compare(*b.Env)
}

// Compare two environment configs
func (a Environment) Compare(b Environment) bool {
	if *a.Image != *b.Image {
		return false
	}
	for keyA, valA := range *a.Vars {
		valB, ok := (*b.Vars)[keyA]
		if !ok {
			return false
		}
		if valA != valB {
			return false
		}
	}
	return true
}
