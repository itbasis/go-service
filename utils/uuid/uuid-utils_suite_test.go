package uuid_test

import (
	"testing"

	itbasisTestUtils "github.com/itbasis/go-test-utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUuidUtils(t *testing.T) {
	RegisterFailHandler(Fail)
	itbasisTestUtils.ConfigureTestLoggerForGinkgo()
	RunSpecs(t, "UUID Utils")
}
