package time

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog"
)

var ErrParsingTime = errors.New("error parsing string to Time")

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
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msgf("value: %v", value)

	if value == nil {
		errParsing := fmt.Errorf("%w: '%v'", ErrParsingTime, value)
		logger.Error().Err(errParsing).Msg("")

		return nil, errParsing
	}

	result, err := time.ParseInLocation(layout, *value, GlobalTime)
	if err != nil {
		errParsing := fmt.Errorf("%w: '%s'", ErrParsingTime, *value)
		logger.Error().Err(err).Msg(errParsing.Error())

		return nil, errParsing
	}

	return &result, nil
}
