package time

import (
	"errors"
	"time"
)

var GlobalTime = time.UTC

var ErrParsingTime = errors.New("error parsing string to Time")
