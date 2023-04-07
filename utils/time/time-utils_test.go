package time_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	serviceTimeUtils "github.com/itbasis/go-service/utils/time"
	testUtils "github.com/itbasis/go-test-utils"
	"github.com/stretchr/testify/assert"
)

var testData = []struct {
	string string
	time   time.Time
}{
	{string: "2022-08-02T06:48:38Z", time: time.Date(2022, time.August, 2, 6, 48, 38, 0, time.UTC)},
	{string: "1993-12-23T10:24:41Z", time: time.Date(1993, time.December, 23, 10, 24, 41, 0, time.UTC)},
}

func TestTime2String(t *testing.T) {
	for _, test := range testData {
		t.Run(
			test.string, func(t *testing.T) {
				assert.Equal(t, test.string, *serviceTimeUtils.Time2String(test.time))
			},
		)
	}
}

func TestString2Time_Success(t *testing.T) {
	for _, test := range testData {
		t.Run(
			test.string, func(t *testing.T) {
				actual, err := serviceTimeUtils.String2Time(testUtils.TestLoggerWithContext(context.Background()), &test.string)
				assert.Nil(t, err)
				assert.Equal(t, test.time, *actual)
			},
		)
	}
}

func TestString2Time_Fail_Nil(t *testing.T) {
	actual, err := serviceTimeUtils.String2Time(testUtils.TestLoggerWithContext(context.Background()), nil)
	assert.Nil(t, actual)
	assert.EqualError(t, err, "error parsing string to Time: '<nil>'")
}

func TestString2Time_Fail(t *testing.T) {
	tests := []struct {
		value     string
		expectErr string
	}{
		{value: "", expectErr: "error parsing string to Time: ''"},
		{value: "t", expectErr: "error parsing string to Time: 't'"},
		{value: "2022", expectErr: "error parsing string to Time: '2022'"},
		{value: "20221228", expectErr: "error parsing string to Time: '20221228'"},
		{value: "2022-12-28", expectErr: "error parsing string to Time: '2022-12-28'"},
		{value: "2022-12-28 02:01:00", expectErr: "error parsing string to Time: '2022-12-28 02:01:00'"},
	}

	for i, test := range tests {
		t.Run(
			fmt.Sprintf("%d: %s", i, test.value), func(t *testing.T) {
				actual, err := serviceTimeUtils.String2Time(testUtils.TestLoggerWithContext(context.Background()), &test.value)
				assert.Nil(t, actual)
				assert.EqualError(t, err, test.expectErr)
			},
		)
	}
}
