package files

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	logUtils "github.com/itbasis/go-log-utils"
	testUtils "github.com/itbasis/go-test-utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFileUtils(t *testing.T) {
	RegisterFailHandler(Fail)

	testUtils.TestLogger.Output(GinkgoWriter)
	logUtils.ConfigureDefaultContextLogger(false)

	RunSpecs(t, "File Utils")
}

var _ = DescribeTable(
	"SHA256 checksum", func(fileName, wantHash string, wantSize int) {
		file, err := os.Open(filepath.Clean(fileName))
		Ω(err).Should(Succeed())

		size, hash, err := Sha256Hash(context.Background(), file)
		Ω(err).Should(Succeed())
		Ω(hash).To(Equal(wantHash))
		Ω(size).To(Equal(uint64(wantSize)))
	},
	Entry(nil, "files.go", "3cde1598374bc811ba3311f6d78a763fc1df4e528f067bbe2ab74cdbd90469a2", 820),
	Entry(nil, "../uuid/uuid-utils.go", "e7c38db11fee2e6b771641f936d9b65ce95946e65b2adc2803215935b345e3af", 521),
)
