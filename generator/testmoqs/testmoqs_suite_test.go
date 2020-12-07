package testmoqs_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMoqs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TestMoqs Suite")
}
