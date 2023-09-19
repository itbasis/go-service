package time

import (
	"context"
	"fmt"
	"time"

	"github.com/juju/zaputil/zapctx"
	"github.com/pkg/errors"
)

// TODO Добавить таймзону в тесты и парсинг
const timeLayout = time.RFC3339

//goland:noinspection GoNameStartsWithPackageName
func Time2String(value time.Time) *string {
	s := value.Format(timeLayout)

	return &s
}

func String2Time(ctx context.Context, value *string) (*time.Time, error) {
	return String2TimeWithCustomLayout(ctx, value, timeLayout)
}

func String2TimeWithCustomLayout(ctx context.Context, value *string, layout string) (*time.Time, error) {
	logger := zapctx.Logger(ctx).Sugar()
	logger.Debugf("value: %v", value)

	if value == nil {
		errParsing := fmt.Errorf("%w: '%v'", ErrParsingTime, value)
		logger.Error(errParsing)

		return nil, errParsing
	}

	if *value == "" {
		errParsing := fmt.Errorf("%w: ''", ErrParsingTime)
		logger.Error(errParsing)

		return nil, errParsing
	}

	result, err := time.ParseInLocation(layout, *value, GlobalTime)
	if err != nil {
		errParsing := fmt.Errorf("%w: '%s'", ErrParsingTime, *value)
		logger.Error(errors.Wrap(err, errParsing.Error()))

		return nil, errParsing
	}

	return &result, nil
}
