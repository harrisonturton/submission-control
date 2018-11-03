package types

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"io"
	"time"
)

// TestJob holds metadata around a TestConfig. It is
// the datastructure that is passed around internally.
type TestJob struct {
	Config    TestConfig
	Timestamp time.Time
}

// TestConfig is the configuration of the testing
// environment. It gives the information on what
// images to use, and what tests to run.
type TestConfig struct {
	Version string  `json:"version"`
	Env     TestEnv `json:"environment"`
}

// TestEnv is the environment within the Docker
// containers, and the image the container is
// built from.
type TestEnv struct {
	Image string            `json:"image"`
	Vars  map[string]string `json:"vars"`
}

// TestResult is the results from running a
// container.
type TestResult struct {
	Stdout    string    `json:"stdout"`
	Stderr    string    `json:"stderr"`
	Timestamp time.Time `json:"timestamp"`
}

var defaultConfig = TestConfig{
	Version: "1",
	Env: TestEnv{
		Image: "",
		Vars:  map[string]string{},
	},
}

// UnmarshalJSON will populate the TestJob using JSON
// data. It will return errors on malformed input and
// missing fields.
func (config *TestConfig) UnmarshalJSON(data io.Reader) error {
	*config = defaultConfig
	err := json.NewDecoder(data).Decode(config)
	if err != nil {
		return errors.New("could not decode json: " + err.Error())
	}
	if config.Env.Image == "" {
		return errors.New("field environment.image cannot be blank")
	}
	return nil
}

// MarshalJSON will convert a TestResult into JSON.
func (result *TestResult) MarshalJSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(result)
}

// Serialize will convert a TestJob into bytes
func (job *TestJob) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(job)
	return buf.Bytes(), err
}

// Deserialize will populate the TestJob fields with
// data from the serialization.
func (job *TestJob) Deserialize(data []byte) error {
	reader := bytes.NewReader(data)
	decoder := gob.NewDecoder(reader)
	return decoder.Decode(job)
}

// Serialize will convert a TestResult into bytes
func (result *TestResult) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(result)
	return buf.Bytes(), err
}

// Deserialize will populate the TestResult fields with
// data from the serialization.
func (result *TestResult) Deserialize(data []byte) error {
	reader := bytes.NewReader(data)
	decoder := gob.NewDecoder(reader)
	return decoder.Decode(result)
}
