// +build unit

package parser

import (
	"strings"
	"testing"
)

const (
	noImage = `{"environment":{"vars":{"k":"v"}}}`  // Fail
	noVars  = `{"environment":{"image":"haskell"}}` // Pass
	noEnv   = `{"version":"1"}`                     // Fail
)

var (
	str       = "string"
	kv        = make(map[string]string)
	stringPtr = &str
	mapPtr    = &kv
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
