package time_test

import (
	"context"
	"fmt"
	"time"

	serviceTimeUtils "github.com/itbasis/go-service/v2/utils/time"
	. "github.com/onsi/ginkgo/v2" //nolint:revive
	. "github.com/onsi/gomega"    //nolint:revive
)

var testData = []struct {
	string string
	time   time.Time
}{
	{string: "2022-08-02T06:48:38Z", time: time.Date(2022, time.August, 2, 6, 48, 38, 0, time.UTC)},
	{string: "1993-12-23T10:24:41Z", time: time.Date(1993, time.December, 23, 10, 24, 41, 0, time.UTC)},
}

var _ = Describe(
	"Time2String", func() {
		for idx, test := range testData {
			It(
				fmt.Sprintf("#%d: %s", idx, test.string), func() {
					Ω(serviceTimeUtils.Time2String(test.time)).To(HaveValue(Equal(test.string)))
				},
			)
		}
	},
)

var _ = Describe(
	"String2Time", func() {
		for idx, test := range testData {
			It(
				fmt.Sprintf("#%d: %s", idx, test.string), func() {
					//nolint:gosec
					Ω(serviceTimeUtils.String2Time(context.Background(), &test.string)).To(HaveValue(BeEquivalentTo(test.time)))
				},
			)
		}

		It(
			"nil", func() {
				actual, err := serviceTimeUtils.String2Time(context.Background(), nil)
				Ω(actual).To(BeNil())
				Ω(err).Should(MatchError("error parsing string to Time: '<nil>'"))
			},
		)

		DescribeTable(
			"Fail", func(value, wantErr string) {
				actual, err := serviceTimeUtils.String2Time(context.Background(), &value)
				Ω(actual).To(BeNil())
				Ω(err).Should(MatchError(wantErr))
			},
			Entry(nil, "", "error parsing string to Time: ''"),
			Entry(nil, "t", "error parsing string to Time: 't'"),
			Entry(nil, "2022", "error parsing string to Time: '2022'"),
			Entry(nil, "20221228", "error parsing string to Time: '20221228'"),
			Entry(nil, "2022-12-28", "error parsing string to Time: '2022-12-28'"),
			Entry(nil, "2022-12-28 02:01:00", "error parsing string to Time: '2022-12-28 02:01:00'"),
		)
	},
)
