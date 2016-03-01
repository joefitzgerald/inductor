package tpl_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTpl(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tpl Suite")
}
