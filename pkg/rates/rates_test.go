package rates

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCurrencyRates(t *testing.T) {
	t.Run(
		"valid conversion from USD to EUR",
		func(t *testing.T) {
			t.Parallel()
			param := &RequestParameter{From: "USD", To: "EUR"}
			expected := 0.85

			rate, err := CurrencyRates(param)
			assert.NoError(
				t,
				err,
			)
			assert.Equal(
				t,
				expected,
				rate,
			)
		},
	)

	t.Run(
		"unsupported currency",
		func(t *testing.T) {
			t.Parallel()
			param := &RequestParameter{From: "XYZ", To: "USD"}

			_, err := CurrencyRates(param)
			assert.Error(
				t,
				err,
			)
		},
	)

	t.Run(
		"conversion rate not found",
		func(t *testing.T) {
			t.Parallel()
			param := &RequestParameter{From: "USD", To: "XYZ"}

			_, err := CurrencyRates(param)
			assert.Error(
				t,
				err,
			)
		},
	)
}
