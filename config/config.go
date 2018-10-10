package config

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Dependencies []string `yaml:"dependencies"`
	BuildCommand string   `yaml:"build-command"`
	Tests        []string `yaml:"test-commands"`
}

// Attempt to read a config file at a specified path
func ReadConfig(path string) (*Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("Failed to read " + path + ".\nDoes it have the correct name, and is in the right directory?")
	}
	c := Config{}
	err = yaml.Unmarshal(file, &c)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal "+path)
	}
	return &c, nil
}
