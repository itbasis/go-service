package files

import (
	"bufio"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"

	"github.com/rs/zerolog"
)

const defaultBufferSize = 64 * 1024 // 64Kb

func Sha256Hash(ctx context.Context, reader io.Reader) (uint64, string, error) {
	logger := zerolog.Ctx(ctx)
	bfr := bufio.NewReader(reader)

	hash := sha256.New()
	buf := make([]byte, defaultBufferSize)

	totalSize := 0

	for {
		n, err := bfr.Read(buf)

		if err == nil {
			logger.Debug().Msgf("read bytes: %d", n)

			totalSize += n
			hash.Write(buf[:n])

			continue
		}

		if errors.Is(err, io.EOF) {
			s := fmt.Sprintf("%x", hash.Sum(nil))

			logger.Debug().Msgf("reading %d bytes. sha256: %s", totalSize, s)

			return uint64(totalSize), s, nil
		}

		logger.Error().Err(err).Msgf("Buffer read error for SHA256 calculation")

		return 0, "", err
	}
}
