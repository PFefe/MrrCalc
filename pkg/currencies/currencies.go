// Package currencies pkg/currencies/currencies.go
package currencies

import (
	"MrrCalc/pkg/rates"
	"fmt"
)

type Currency struct {
	From   string
	To     string
	Amount float64
}

type Converter struct {
	provider rates.ExchangeRateProvider
}

func (c *Converter) Convert(params *Currency) (float64, error) {
	param := &rates.RequestParameter{
		From: params.From,
		To:   params.To,
	}
	rate, err := c.provider.GetRate(param)
	if err != nil {
		return 0, fmt.Errorf(
			"error getting rate: %w",
			err,
		)
	}
	return params.Amount * rate, nil
}

func NewConverter(provider rates.ExchangeRateProvider) *Converter {
	return &Converter{
		provider: provider,
	}
}
