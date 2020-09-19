package testing

//go:generate moqueries --destination moq_testing.go --export MoqT

// MoqT is that interface defining standard library *testing.T methods used by
// Moqueries
type MoqT interface {
	Fatalf(format string, args ...interface{})
}
