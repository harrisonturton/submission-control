// +build unit

package parser

import (
	"github.com/harrisonturton/submission-control/ci/types"
	"reflect"
	"strings"
	"testing"
)

const (
	noImage = `{"environment":{"vars":{"k":"v"}}}`  // Fail
	noVars  = `{"environment":{"image":"haskell"}}` // Pass
	noEnv   = `{"version":"1"}`                     // Fail
)

func TestFailures(t *testing.T) {
	testMissingField := func(err error, testCase string) {
		if err == nil {
			t.Errorf("Failed to reject bad data: %s", testCase)
		}
	}
	// Test missing env field
	reader := strings.NewReader(noEnv)
	_, err := ParseConfig(reader)
	testMissingField(err, noEnv)
	// Test missing env.image field
	reader = strings.NewReader(noImage)
	_, err = ParseConfig(reader)
	testMissingField(err, noImage)
}

func TestPassing(t *testing.T) {
	testPresentField := func(err error, testCase string) {
		if err != nil {
			t.Errorf("Failed to parse good data: %s", testCase)
		}
	}
	// Test missing vars
	reader := strings.NewReader(noVars)
	_, err := ParseConfig(reader)
	testPresentField(err, noImage)
}

func TestSerialization(t *testing.T) {
	version := "1"
	image := "haskell"
	config := types.TestConfig{
		Version: &version,
		Env: &types.Environment{
			Image: &image,
		},
	}
	data, err := SerializeConfig(config)
	if err != nil {
		t.Fatalf("Failed to serialize config: %s", err)
	}

	result, err := DeserializeConfig(data)
	if err != nil {
		t.Fatalf("Failed to serialize config: %s", err)
	}
	if !reflect.DeepEqual(config, result) {
		t.Errorf("Serialization changes the config - bad!")
	}
}
