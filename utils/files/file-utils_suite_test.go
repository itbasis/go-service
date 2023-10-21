package files_test

import (
	"testing"

	_ "github.com/itbasis/go-test-utils/v2"
	. "github.com/onsi/ginkgo/v2" //nolint:revive
	. "github.com/onsi/gomega"    //nolint:revive
)

func TestFilesUtils(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Files Utils")
}
