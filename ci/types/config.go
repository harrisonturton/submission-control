package types

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"
)

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

// UnmarshalJSON will populate the TestJob fields by
// reading JSON data.
func (config *TestConfig) UnmarshalJSON(data io.Reader) error {
	err := json.NewDecoder(data).Decode(config)
	if err != nil {
		return fmt.Errorf("could not parse config: %s", err)
	}
	if config.Version == nil {
		version := "1"
		config.Version = &version
	}
	if config.Env == nil {
		return errors.New("could not parse config: missing environment field")
	}
	if config.Env.Image == nil {
		return errors.New("could not parse config: missing environment.image field")
	}
	if config.Env.Vars == nil {
		vars := make(map[string]string)
		config.Env.Vars = &vars
	}
	return nil
}

// Serialize will convert a TestConfig into bytes
func (config *TestConfig) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(config)
	return buf.Bytes(), err
}

// Deserialize will populate the TestConfig fields with
// data from the serialization.
func (config *TestConfig) Deserialize(data []byte) error {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	return decoder.Decode(config)
}