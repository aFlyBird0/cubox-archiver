package cubox

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCubox(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cubox Suite")
}
