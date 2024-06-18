// Package currencies pkg/currencies/currencies.go
package currencies

import (
	"fmt"

	rates "MrrCalc/pkg/rates"
)

type Currency struct {
	From   string
	To     string
	Amount float64
}

func ConvertCurrency(params *Currency) (float64, error) {
	param := &rates.RequestParameter{
		From: params.From,
		To:   params.To,
	}
	rate, err := rates.CurrencyRates(param)
	if err != nil {
		return 0, fmt.Errorf(
			"error getting rate: %w",
			err,
		)
	}
	return params.Amount * rate, nil
}
