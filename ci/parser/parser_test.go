package parser

import (
	"errors"
	"github.com/harrisonturton/submission-control/ci/types"
	"strings"
	"testing"
)

type tableTest struct {
	input    string
	expected types.TestConfig
}

func TestParse(t *testing.T) {
	compare := func(input string, expected types.TestConfig) error {
		reader := strings.NewReader(input)
		config, err := ParseConfig(reader)
		if err != nil {
			return err
		}
		if !expected.Compare(config) {
			return errors.New("actual and expected config is different")
		}
		return nil
	}

	testAll(t, configTests, compare)
}

func testAll(t *testing.T, table []tableTest, compare func(actual string, expected types.TestConfig) error) {
	for i, testCase := range table {
		if err := compare(testCase.input, testCase.expected); err != nil {
			t.Errorf("Failed on table test %d: %s", i, err)
		}
	}
}

var configTests = []tableTest{
	{
		input: `
{
	"version": "1"
}
`,
		expected: types.TestConfig{
			Version: "1",
			Env:     nil,
		},
	},
	{
		input: `
{
	"version": "1",
	"environment": {}
}
`,
		expected: types.TestConfig{
			Version: "1",
			Env:     &types.Environment{},
		},
	},
	{
		input: `
{
	"version": "1",
	"environment": {
		"image": "haskell",
		"vars": {
			"user": "Ubuntu",
			"length": "30"
		}
	}
}
`,
		expected: types.TestConfig{
			Version: "1",
			Env: &types.Environment{
				Image: "haskell",
				Vars: map[string]string{
					"user":   "Ubuntu",
					"length": "30",
				},
			},
		},
	},
}
