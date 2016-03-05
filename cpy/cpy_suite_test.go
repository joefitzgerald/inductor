package cpy_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCpy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cpy Suite")
}
