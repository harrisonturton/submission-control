package types

import "testing"

var compareTests = []struct {
	inputA   TestConfig
	inputB   TestConfig
	expected bool
}{
	{
		inputA: TestConfig{
			Version: "1",
			Env:     nil,
		},
		inputB: TestConfig{
			Version: "1",
			Env:     nil,
		},
		expected: true,
	},
	{
		inputA: TestConfig{
			Version: "1",
			Env:     nil,
		},
		inputB: TestConfig{
			Version: "1",
			Env:     &Environment{},
		},
		expected: false,
	},
	{
		inputA: TestConfig{
			Version: "1",
			Env: &Environment{
				Image: "haskell",
				Vars:  map[string]string{},
			},
		},
		inputB: TestConfig{
			Version: "1",
			Env: &Environment{
				Image: "haskell",
				Vars:  map[string]string{},
			},
		},
		expected: true,
	},
	{
		inputA: TestConfig{
			Version: "1",
			Env: &Environment{
				Image: "haskell",
				Vars: map[string]string{
					"user":     "ubuntu",
					"duration": "10s",
				},
			},
		},
		inputB: TestConfig{
			Version: "1",
			Env: &Environment{
				Image: "haskell",
				Vars: map[string]string{
					"user":     "ubuntu",
					"duration": "10s",
				},
			},
		},
		expected: true,
	},
	{
		inputA: TestConfig{
			Version: "1",
			Env: &Environment{
				Image: "haskell",
				Vars: map[string]string{
					"user":     "ubuntu",
					"duration": "10s",
				},
			},
		},
		inputB: TestConfig{
			Version: "1",
			Env: &Environment{
				Image: "haskell",
				Vars: map[string]string{
					"user": "ubuntu",
				},
			},
		},
		expected: false,
	},
}

func TestConfigCompare(t *testing.T) {
	for i, testCase := range compareTests {
		a := testCase.inputA
		b := testCase.inputB
		expected := testCase.expected
		if a.Compare(b) != expected {
			t.Errorf("Failed table test %d", i)
		}
	}
}
