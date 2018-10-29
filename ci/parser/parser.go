package parser

import (
	"encoding/json"
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
	return config, nil
}
