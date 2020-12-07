package moq

import "sync/atomic"

// Scene stores a collection of moqs so that they can work together
type Scene struct {
	T               T
	moqs            []Moq
	nextRecorderSeq uint32
	nextMockSeq     uint32
}

//go:generate moqueries --destination moq_moq_test.go Moq

// Moq is implemented by all moqs so that they can integrate with a scene
type Moq interface {
	Reset()
	AssertExpectationsMet()
}

//go:generate moqueries --destination moq_testing.go --export T

// T is the interface defining standard library *testing.T methods used by
// Moqueries
type T interface {
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

// NewScene creates a new empty scene with no moqs
func NewScene(t T) *Scene {
	return &Scene{
		T:               t,
		nextRecorderSeq: 0,
		nextMockSeq:     0,
	}
}

// AddMoq adds a moq to a scene
func (s *Scene) AddMoq(m Moq) {
	s.moqs = append(s.moqs, m)
}

// Reset resets the state of all moqs in the scene so that they can be used in
// another test
func (s *Scene) Reset() {
	for _, m := range s.moqs {
		m.Reset()
	}
}

// AssertExpectationsMet asserts that all expectations for all moqs in the
// scene are met
func (s *Scene) AssertExpectationsMet() {
	for _, m := range s.moqs {
		m.AssertExpectationsMet()
	}
}

// NextRecorderSequence returns the next sequence value for a recorder when
// recording expectations
func (s *Scene) NextRecorderSequence() uint32 {
	return atomic.AddUint32(&s.nextRecorderSeq, 1)
}

// NextMockSequence returns the next sequence value when a call is being made
// to a mock
func (s *Scene) NextMockSequence() uint32 {
	return atomic.AddUint32(&s.nextMockSeq, 1)
}
