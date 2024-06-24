package currencies

import (
	"MrrCalc/pkg/rates"
	"fmt"
	"testing"
)

type mockRateProvider struct{}

func (m *mockRateProvider) GetRate(param *rates.RequestParameter) (float64, error) {
	if param.From == "EUR" && param.To == "GBP" {
		return 0.85, nil
	}
	return 0, fmt.Errorf(
		"unsupported currency pair: %s to %s",
		param.From,
		param.To,
	)
}

func TestConverter_Convert(t *testing.T) {
	provider := &mockRateProvider{}
	converter := NewConverter(provider)

	tests := []struct {
		from   string
		to     string
		amount float64
		want   float64
	}{
		{"EUR", "GBP", 100, 85},
		{"EUR", "GBP", 200, 170},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf(
				"%s to %s",
				tt.from,
				tt.to,
			),
			func(t *testing.T) {
				params := &Currency{
					From:   tt.from,
					To:     tt.to,
					Amount: tt.amount,
				}
				got, err := converter.Convert(params)
				if err != nil {
					t.Fatalf(
						"unexpected error: %v",
						err,
					)
				}
				if got != tt.want {
					t.Errorf(
						"got %f, want %f",
						got,
						tt.want,
					)
				}
			},
		)
	}
}
