package utils_test

import (
	"context"
	"testing"

	"github.com/itbasis/go-service/utils"
	testUtils "github.com/itbasis/go-test-utils"
	"github.com/stretchr/testify/assert"
)

func TestString2UUID_Success(t *testing.T) {
	tests := []struct {
		value string
	}{
		{"5ae42556-9272-4052-9985-04aba5b7d5ed"},
		{"e4dc94d9-7ca6-4d13-9663-de4952a1f31b"},
	}
	for _, test := range tests {
		t.Run(
			test.value, func(t *testing.T) {
				result, err := utils.String2UUID(testUtils.TestLoggerWithContext(context.Background()), test.value)
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, test.value, result.String())
			},
		)
	}
}

func TestString2UUID_Fail(t *testing.T) {
	testData := []struct {
		value     string
		expectErr string
	}{
		{value: "", expectErr: "error parsing string to UUID: ''"},
		{value: ".", expectErr: "error parsing string to UUID: '.'"},
		{value: "_", expectErr: "error parsing string to UUID: '_'"},
		{value: "b3edafe2-c4c0-4cf9-8a8b-8c350ed5d27", expectErr: "error parsing string to UUID: 'b3edafe2-c4c0-4cf9-8a8b-8c350ed5d27'"},
	}
	for _, test := range testData {
		t.Run(
			test.value, func(t *testing.T) {
				uid, err := utils.String2UUID(testUtils.TestLoggerWithContext(context.Background()), test.value)
				assert.True(t, uid.IsNil())
				assert.EqualError(t, err, test.expectErr)
			},
		)
	}
}
