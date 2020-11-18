package moq

// ExpectationMode determines the behavior of a mock when a method is invoked
// with no matching expectations
type ExpectationMode int

const (
	// Strict mode causes a mock to validate each method invocation
	Strict ExpectationMode = iota
	// Nice mode will return zero values for any unexpected invocation
	Nice
)

// SequenceMode is used in conjunction with the generated seq and noSeq methods
// when checking call sequences
type SequenceMode int

const (
	// SeqDefaultOff indicates that call sequences will not be reserved for any
	// calls but individual calls can turn on sequences
	SeqDefaultOff SequenceMode = iota
	// SeqDefaultOn indicates that call sequences will be reserved for all
	// calls but individual calls can turn off sequences
	SeqDefaultOn
)

// MockConfig is passed to add generated mocks to configure the mock
type MockConfig struct {
	Expectation ExpectationMode
	Sequence    SequenceMode
}
