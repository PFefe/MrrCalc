package currencies

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertCurrency(t *testing.T) {
	t.Run(
		"valid conversion from USD to EUR",
		func(t *testing.T) {
			params := &Currency{
				From:   "USD",
				To:     "EUR",
				Amount: 100.0,
			}
			expected := 85.0 // 100 * 0.85

			result, err := ConvertCurrency(params)
			assert.NoError(
				t,
				err,
			)
			assert.Equal(
				t,
				expected,
				result,
			)
		},
	)

	t.Run(
		"unsupported currency",
		func(t *testing.T) {
			params := &Currency{
				From:   "XYZ",
				To:     "USD",
				Amount: 100.0,
			}

			_, err := ConvertCurrency(params)
			assert.Error(
				t,
				err,
			)
		},
	)

	t.Run(
		"conversion rate not found",
		func(t *testing.T) {
			params := &Currency{
				From:   "USD",
				To:     "XYZ",
				Amount: 100.0,
			}

			_, err := ConvertCurrency(params)
			assert.Error(
				t,
				err,
			)
		},
	)
}
