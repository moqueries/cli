package moq

// ExpectationMode determines the behavior of a moq when a method is invoked
// with no matching expectations
type ExpectationMode int

const (
	// Strict mode causes a moq to validate each method invocation
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

// Config is passed to a moq to provide configuration for the moq
type Config struct {
	Expectation ExpectationMode
	Sequence    SequenceMode
}

// ParamIndexing values determine how parameters are indexed in a moq
type ParamIndexing int

const (
	// ParamIndexByValue indicates that a specific parameter of a specific
	// function will be indexed by it value alone. The parameter value will be
	// copied into a parameter key so simple equality will determine if an
	// expectation matches an actual call. The exact same instance must be
	// supplied to both the expectation call and the actual call.
	ParamIndexByValue ParamIndexing = iota
	// ParamIndexByHash indicates that a specific parameter of a specific
	// function will be indexed by a deep hash value. A deep hash library is
	// used to uniquely identify the parameter's value which includes the
	// values of any parameter subtypes. The exact same instance will only
	// match an expectation to an actual call if the internal state of the
	// instance hasn't changed.
	ParamIndexByHash
)
