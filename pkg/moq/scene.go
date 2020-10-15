package moq

// Scene stores a collection of mocks so that they can work together
type Scene struct {
	MoqT  MoqT
	mocks []Mock
}

//go:generate moqueries --destination moq_mock_test.go Mock

// Mock is implemented by all mocks so that they can work in a scene
type Mock interface {
	Reset()
	AssertExpectationsMet()
}

//go:generate moqueries --destination moq_testing.go --export MoqT

// MoqT is that interface defining standard library *testing.T methods used by
// Moqueries
type MoqT interface {
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

// NewScene creates a new scene with no mocks
func NewScene(t MoqT) *Scene {
	return &Scene{
		MoqT: t,
	}
}

// AddMock adds a mock to a scene
func (s *Scene) AddMock(m Mock) {
	s.mocks = append(s.mocks, m)
}

// Reset resets the state of all mocks in the scene so that they can be used in
// another test
func (s *Scene) Reset() {
	for _, m := range s.mocks {
		m.Reset()
	}
}

// AssertExpectationsMet asserts that all expectations for all mock in the
// scene are met
func (s *Scene) AssertExpectationsMet() {
	for _, m := range s.mocks {
		m.AssertExpectationsMet()
	}
}
