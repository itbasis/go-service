package utils

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofrs/uuid/v5"
	"github.com/rs/zerolog"
)

var ErrParsingUUID = errors.New("error parsing string to UUID")

func String2UUID(ctx context.Context, value string) (uuid.UUID, error) {
	logger := zerolog.Ctx(ctx)

	result, err := uuid.FromString(value)
	if err != nil {
		msg := fmt.Errorf("%w: '%s'", ErrParsingUUID, value)

		logger.Error().Err(err).Msg(msg.Error())

		return uuid.Nil, msg
	}

	return result, nil
}
