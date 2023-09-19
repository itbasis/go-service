package files

import (
	"bufio"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"

	"github.com/juju/zaputil/zapctx"
)

const defaultBufferSize = 64 * 1024 // 64Kb

func Sha256Hash(ctx context.Context, reader io.Reader) (uint64, string, error) {
	logger := zapctx.Logger(ctx).Sugar()
	bfr := bufio.NewReader(reader)

	hash := sha256.New()
	buf := make([]byte, defaultBufferSize)

	totalSize := 0

	for {
		//nolint:varnamelen
		n, err := bfr.Read(buf)

		if err == nil {
			logger.Debugf("read bytes: %d", n)

			totalSize += n
			hash.Write(buf[:n])

			continue
		}

		if errors.Is(err, io.EOF) {
			s := fmt.Sprintf("%x", hash.Sum(nil))

			logger.Debugf("reading %d bytes. sha256: %s", totalSize, s)

			return uint64(totalSize), s, nil
		}

		err = fmt.Errorf("buffer read error for SHA256 calculation: %w", err)
		logger.Error(err)

		return 0, "", err
	}
}
