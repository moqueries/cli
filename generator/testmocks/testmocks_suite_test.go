package testmocks_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMocks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TestMocks Suite")
}
