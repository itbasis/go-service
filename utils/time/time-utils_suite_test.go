package time_test

import (
	"testing"

	_ "github.com/itbasis/go-test-utils/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTimeUtils(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "TimeUtils")
}
