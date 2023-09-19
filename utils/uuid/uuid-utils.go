package uuid

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid/v5"
	"github.com/juju/zaputil/zapctx"
	"github.com/pkg/errors"
)

func String2UUID(ctx context.Context, value string) (uuid.UUID, error) {
	logger := zapctx.Logger(ctx).Sugar()
	logger.Debugf("value: %s", value)

	result, err := uuid.FromString(value)
	if err != nil {
		msg := fmt.Errorf("%w: '%s'", ErrParsingUUID, value)

		logger.Error(errors.Wrapf(err, msg.Error()))

		return uuid.Nil, msg
	}

	return result, nil
}
