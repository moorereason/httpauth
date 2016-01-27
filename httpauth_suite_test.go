package httpauth_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestHttpauth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Httpauth Suite")
}
