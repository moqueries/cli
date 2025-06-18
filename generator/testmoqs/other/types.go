// Package other contains multiple test types for use in unit testing
package other

// Other is used for testing embedded interfaces
//nolint: iface // redundant for testing
type Other interface {
	Another
}

// Params encapsulates the parameters for use in various test types
type Params struct {
	SParam string
	BParam bool
}

// Results encapsulates the results for use in various test types
type Results struct {
	SResult string
	Err     error
}

// Another is used for testing embedded interfaces
//nolint: iface // redundant for testing
type Another interface {
	Other(Params) Results
}
