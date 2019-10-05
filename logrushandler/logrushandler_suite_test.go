package logrushandler_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLogrushandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Logrushandler Suite")
}
