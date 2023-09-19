package files_test

import (
	"context"
	"os"
	"path/filepath"

	"github.com/itbasis/go-service/v2/utils/files"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable(
	"SHA256 checksum", func(fileName, wantHash string, wantSize int) {
		file, err := os.Open(filepath.Clean(fileName))
		Ω(err).Should(Succeed())

		size, hash, err := files.Sha256Hash(context.Background(), file)
		Ω(err).Should(Succeed())
		Ω(hash).To(Equal(wantHash))
		Ω(size).To(Equal(uint64(wantSize)))
	},
	Entry(nil, "file-utils.go", "c31994d1b870e88a6e4f3c0b2ae6034ef2b22843ed3504847e4c6410d3bb9a6f", 866),
	Entry(nil, "../uuid/uuid-utils.go", "dfa0e70015284219719c1da0d8cb1c1512136e9e1f01d66500727b0753d2067f", 488),
)
