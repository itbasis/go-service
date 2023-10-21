package uuid_test

import (
	"context"

	"github.com/gofrs/uuid/v5"
	itbasisServiceUtilsUUID "github.com/itbasis/go-service/v2/utils/uuid"
	. "github.com/onsi/ginkgo/v2" //nolint:revive
	. "github.com/onsi/gomega"    //nolint:revive
)

var _ = Describe(
	"String2UUID", func() {
		DescribeTable(
			"Success", func(value string) {
				want, err := uuid.FromString(value)
				Ω(err).Should(Succeed())

				Ω(itbasisServiceUtilsUUID.String2UUID(context.Background(), value)).To(Equal(want))
			},
			Entry(nil, "5ae42556-9272-4052-9985-04aba5b7d5ed"),
			Entry(nil, "e4dc94d9-7ca6-4d13-9663-de4952a1f31b"),
		)

		DescribeTable(
			"Fail", func(value, wantErr string) {
				actual, err := itbasisServiceUtilsUUID.String2UUID(context.Background(), value)
				Ω(actual).To(Equal(uuid.Nil))
				Ω(err).Should(MatchError(wantErr))
			},
			Entry(nil, "", "error parsing string to UUID: ''"),
			Entry(nil, ".", "error parsing string to UUID: '.'"),
			Entry(nil, "_", "error parsing string to UUID: '_'"),
			Entry(nil, "b3edafe2-c4c0-4cf9-8a8b-8c350ed5d27", "error parsing string to UUID: 'b3edafe2-c4c0-4cf9-8a8b-8c350ed5d27'"),
		)
	},
)
