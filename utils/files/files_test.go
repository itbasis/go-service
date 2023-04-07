package files

import (
	"context"
	"os"
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
		file, err := os.Open(fileName)
		Ω(err).Should(Succeed())

		size, hash, err := Sha256Hash(context.Background(), file)
		Ω(err).Should(Succeed())
		Ω(hash).To(Equal(wantHash))
		Ω(size).To(Equal(uint64(wantSize)))
	},
	Entry(nil, "files.go", "3cde1598374bc811ba3311f6d78a763fc1df4e528f067bbe2ab74cdbd90469a2", 820),
	Entry(nil, "../uuid.go", "8ab15d5d07664cfa6d647248acb831085dd3ee35e32402a176601115ab923311", 481),
)
