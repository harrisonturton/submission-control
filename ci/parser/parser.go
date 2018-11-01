package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/harrisonturton/submission-control/ci/types"
	"io"
)

// ParseConfig takes raw JSON data and unmarshals it into
// a types.TestConfig object
func ParseConfig(data io.ReadCloser) (types.TestConfig, error) {
	defer data.Close()
	config := types.TestConfig{}
	err := json.NewDecoder(data).Decode(&config)
	if err != nil {
		return config, fmt.Errorf("could not parse config: %s", err)
	}
	if config.Version == nil {
		*config.Version = "1"
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
