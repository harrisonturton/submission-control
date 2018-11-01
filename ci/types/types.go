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
