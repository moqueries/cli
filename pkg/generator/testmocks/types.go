package testmocks

//go:generate moqueries --destination moq_usualfn_test.go UsualFn
//go:generate moqueries --destination exported/moq_usualfn.go --export UsualFn

// UsualFn is a typical function type
type UsualFn func(sParam string, bParam bool) (sResult string, err error)

//go:generate moqueries --destination moq_nonamesfn_test.go NoNamesFn
//go:generate moqueries --destination exported/moq_nonamesfn.go --export NoNamesFn

// NoNamesFn is a typical function type
type NoNamesFn func(sParam string, bParam bool) (sResult string, err error)

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

//go:generate moqueries --destination moq_usual_test.go Usual
//go:generate moqueries --destination exported/moq_usual.go --export Usual

// Usual combines all of the above function types into an interface
type Usual interface {
	Usual(sParam string, bParam bool) (sResult string, err error)
	NoNames(string, bool) (string, error)
	NoResults(sParam string, bParam bool)
	NoParams() (sResult string, err error)
	Nothing()
	Variadic(other bool, args ...string) (sResult string, err error)
	RepeatedIds(sParam1, sParam2 string, bParam bool) (sResult1, sResult2 string, err error)
}
