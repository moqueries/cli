package config

// ExpectationMode determines the behavior of a mock when a method is invoked
// with no matching expectations
type ExpectationMode int

const (
	// Strict mode causes a mock to validate each method invocation
	Strict ExpectationMode = iota
	// Nice mode will return zero values for any unexpected invocation
	Nice
)

// MockConfig is passed to add generated mocks to configure the mock
type MockConfig struct {
	Expectation ExpectationMode
}
