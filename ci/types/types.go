package types

// TestConfig is the type that holds the default
// test configurations, and the test commands for
// each test.
type TestConfig struct {
	Defaults TestCase
	Tests    []TestCase
}

// TestCase is the configuration for a single test.
type TestCase struct {
	Command    string
	Expect     Expectation
	Weight     float32
	StopOnFail bool
}

// EnvConfig is the configuration for the testing
// environment inside the Docker container.
type EnvConfig struct {
	Image string
	Vars  map[string]string
}

// Expectation is the expected behaviour of the
// outcome of a test.
type Expectation struct {
	ExitCode    int
	Duration    string
	Contains    string
	ContainsAll []string
	AtLeastOne  []string
	Not         *Expectation
	And         *Expectation
	Or          *Expectation
}
