package parser

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/harrisonturton/submission-control/ci/types"
	"io"
)

// ParseConfig takes raw JSON data and unmarshals it into
// a types.TestConfig object
func ParseConfig(data io.Reader) (types.TestConfig, error) {
	config := types.TestConfig{}
	err := json.NewDecoder(data).Decode(&config)
	if err != nil {
		return config, fmt.Errorf("could not parse config: %s", err)
	}
	if config.Version == nil {
		version := "1"
		config.Version = &version
	}
	if config.Env == nil {
		return config, errors.New("could not parse config: missing environment field")
	}
	if config.Env.Image == nil {
		return config, errors.New("could not parse config: missing environment.image field")
	}
	if config.Env.Vars == nil {
		data := make(map[string]string)
		config.Env.Vars = &data
	}
	return config, nil
}

// SerializeConfig will convert a TestConfig instance into
// a series of bytes.
func SerializeConfig(config types.TestConfig) ([]byte, error) {
	gob.Register(types.TestConfig{})
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(config)
	return buf.Bytes(), err
}

// DeserializeConfig will convert bytes into a TestConfig
// instance.
func DeserializeConfig(data []byte) (types.TestConfig, error) {
	gob.Register(types.TestConfig{})
	var config types.TestConfig
	reader := bytes.NewReader(data)
	decoder := gob.NewDecoder(reader)
	err := decoder.Decode(&config)
	return config, err
}
