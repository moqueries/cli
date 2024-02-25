package other

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

type Another interface {
	Other(Params) Results
}
