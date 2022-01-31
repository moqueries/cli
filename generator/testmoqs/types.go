package testmoqs

import "io"

// NB: Keep in sync with testmoq_test.go TestGenerating

//go:generate moqueries --destination moq_usualfn_test.go UsualFn
//go:generate moqueries --destination exported/moq_usualfn.go --export UsualFn

// UsualFn is a typical function type
type UsualFn func(sParam string, bParam bool) (sResult string, err error)

//go:generate moqueries --destination moq_nonamesfn_test.go NoNamesFn
//go:generate moqueries --destination exported/moq_nonamesfn.go --export NoNamesFn

// NoNamesFn is a typical function type
type NoNamesFn func(string, bool) (string, error)

//go:generate moqueries --destination moq_noresultsfn_test.go NoResultsFn
//go:generate moqueries --destination exported/moq_noresultsfn.go --export NoResultsFn

// NoResultsFn is a function with no return values
type NoResultsFn func(sParam string, bParam bool)

//go:generate moqueries --destination moq_noparamsfn_test.go NoParamsFn
//go:generate moqueries --destination exported/moq_noparamsfn.go --export NoParamsFn

// NoParamsFn is a function with no parameters
type NoParamsFn func() (sResult string, err error)

//go:generate moqueries --destination moq_nothingfn_test.go NothingFn
//go:generate moqueries --destination exported/moq_nothingfn.go --export NothingFn

// NothingFn is a function with no parameters and no return values
type NothingFn func()

//go:generate moqueries --destination moq_variadicfn_test.go VariadicFn
//go:generate moqueries --destination exported/moq_variadicfn.go --export VariadicFn

// VariadicFn is a function with a variable number of arguments
type VariadicFn func(other bool, args ...string) (sResult string, err error)

//go:generate moqueries --destination moq_repeatedidsfn_test.go RepeatedIdsFn
//go:generate moqueries --destination exported/moq_repeatedidsfn.go --export RepeatedIdsFn

// RepeatedIdsFn is a function with multiple arguments of the same type
type RepeatedIdsFn func(sParam1, sParam2 string, bParam bool) (sResult1, sResult2 string, err error)

//go:generate moqueries --destination moq_timesfn_test.go TimesFn
//go:generate moqueries --destination exported/moq_timesfn.go --export TimesFn

// TimesFn takes a parameter called times which should generate a valid moq
type TimesFn func(times string, bParam bool) (sResult string, err error)

//go:generate moqueries --destination moq_difficultparamnamesfn_test.go DifficultParamNamesFn
//go:generate moqueries --destination exported/moq_difficultparamnamesfn.go --export DifficultParamNamesFn

// DifficultParamNamesFn has parameters with names that have been problematic
type DifficultParamNamesFn func(m, r bool, sequence string, param, params int, result, results float32)

//go:generate moqueries --destination moq_difficultresultnamesfn_test.go DifficultResultNamesFn
//go:generate moqueries --destination exported/moq_difficultresultnamesfn.go --export DifficultResultNamesFn

// DifficultResultNamesFn has parameters with names that have been problematic
type DifficultResultNamesFn func() (m, r string, sequence error, param, params int, result, results float32)

// PassByReferenceParams encapsulates the parameters for passing by reference
// tests
type PassByReferenceParams struct {
	SParam string
	BParam bool
}

//go:generate moqueries --destination moq_passbyreferencefn_test.go PassByReferenceFn
//go:generate moqueries --destination exported/moq_passbyreferencefn.go --export PassByReferenceFn

// PassByReferenceFn tests passing parameters by reference
type PassByReferenceFn func(p *PassByReferenceParams) (sResult string, err error)

//go:generate moqueries --destination moq_interfaceparamfn_test.go InterfaceParamFn
//go:generate moqueries --destination exported/moq_interfaceparamfn.go --export InterfaceParamFn

// InterfaceParamWriter is used for testing functions that take an interface as
// a parameter
type InterfaceParamWriter struct {
	SParam string
	BParam bool
}

func (w *InterfaceParamWriter) Write([]byte) (int, error) {
	return 0, nil
}

// InterfaceParamFn tests passing interface parameters
type InterfaceParamFn func(w io.Writer) (sResult string, err error)

//go:generate moqueries --destination moq_interfaceresultfn_test.go InterfaceResultFn
//go:generate moqueries --destination exported/moq_interfaceresultfn.go --export InterfaceResultFn

// InterfaceResultReader is used for testing functions that return an interface
type InterfaceResultReader struct {
	SResult string
	Err     error
}

func (r *InterfaceResultReader) Read(p []byte) (n int, err error) {
	return 0, nil
}

// InterfaceResultFn tests returning interface results
type InterfaceResultFn func(sParam string, bParam bool) (r io.Reader)

//go:generate moqueries --destination moq_usual_test.go Usual
//go:generate moqueries --destination exported/moq_usual.go --export Usual

// Usual combines all the above function types into an interface
type Usual interface {
	Usual(sParam string, bParam bool) (sResult string, err error)
	NoNames(string, bool) (string, error)
	NoResults(sParam string, bParam bool)
	NoParams() (sResult string, err error)
	Nothing()
	Variadic(other bool, args ...string) (sResult string, err error)
	RepeatedIds(sParam1, sParam2 string, bParam bool) (sResult1, sResult2 string, err error)
	Times(sParam string, times bool) (sResult string, err error)
	DifficultParamNames(m, r bool, sequence string, param, params int, result, results float32)
	DifficultResultNames() (m, r string, sequence error, param, params int, result, results float32)
	PassByReference(p *PassByReferenceParams) (sResult string, err error)
	InterfaceParam(w io.Writer) (sResult string, err error)
	InterfaceResult(sParam string, bParam bool) (r io.Reader)
}
