// +build unit

package types

import (
	"strings"
	"testing"
)

const (
	noImage = `{"environment":{"vars":{"k":"v"}}}`  // Fail
	noVars  = `{"environment":{"image":"haskell"}}` // Pass
	noEnv   = `{"version":"1"}`                     // Fail
)

// TestExpectedFailure will test if bad TestJobs are
// always rejected.
func TestExpectedFailure(t *testing.T) {
	expectFailure := func(err error, testCase string) {
		if err == nil {
			t.Errorf("Failed to reject bad data: %s", testCase)
		}
	}
	// Test missing env field
	config := TestConfig{}
	err := config.UnmarshalJSON(strings.NewReader(noEnv))
	expectFailure(err, noEnv)
	// Test missing env.image field
	config = TestConfig{}
	err = config.UnmarshalJSON(strings.NewReader(noImage))
	expectFailure(err, noImage)
}

// TestUnexpectedFailure will test that good tests always
// pass.
func TestUnexpectedFailure(t *testing.T) {
	expectAccept := func(err error, testCase string) {
		if err != nil {
			t.Errorf("Failed to parse good data: %s", testCase)
		}
	}
	config := TestConfig{}
	err := config.UnmarshalJSON(strings.NewReader(noVars))
	expectAccept(err, noVars)
}
