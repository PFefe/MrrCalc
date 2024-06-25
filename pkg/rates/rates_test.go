// rates_test.go
package rates

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

/*type mockExchangeRateProvider struct{}

func (m *mockExchangeRateProvider) GetRate(param *RequestParameter) (float64, error) {
	// Mocking responses based on param.From and param.To
	switch param.From {
	case "USD":
		switch param.To {
		case "EUR":
			return 0.85, nil
		default:
			return 0, fmt.Errorf(
				"unsupported currency: %s",
				param.To,
			)
		}
	default:
		return 0, fmt.Errorf(
			"unsupported currency: %s",
			param.From,
		)
	}
}*/

func TestCurrencyRates(t *testing.T) {
	// mock := &mockExchangeRateProvider{}
	svc := APIService{}

	tests := []struct {
		name     string
		param    *RequestParameter
		expected float64
		wantErr  bool
	}{
		{
			name:     "valid conversion from USD to USD",
			param:    &RequestParameter{From: "USD", To: "USD"},
			expected: 1.00,
			wantErr:  false,
		},
		{
			name:     "unsupported currency",
			param:    &RequestParameter{From: "XYZ", To: "USD"},
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "conversion rate not found",
			param:    &RequestParameter{From: "USD", To: "XYZ"},
			expected: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				t.Parallel()

				rate, err := svc.GetRate(tt.param)

				if tt.wantErr {
					assert.Error(
						t,
						err,
					)
				} else {
					assert.NoError(
						t,
						err,
					)
					assert.Equal(
						t,
						tt.expected,
						rate,
					)
				}
			},
		)
	}
}
