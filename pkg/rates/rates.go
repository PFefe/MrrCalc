// Package rates pkg/rates/rates.go
package rates

import "fmt"

type RequestParameter struct {
	From string
	To   string
}

func CurrencyRates(param *RequestParameter) (float64, error) {
	rates := map[string]map[string]float64{
		"USD": {
			"EUR": 0.85,
			"JPY": 110.0,
			"GBP": 0.72,
			"TRY": 32.50,
			"AED": 3.67,
			"USD": 1.00,
		},
		"EUR": {
			"USD": 1.18,
			"JPY": 129.0,
			"GBP": 0.85,
			"TRY": 38.0,
			"AED": 4.35,
			"EUR": 1.00,
		},
		"GBP": {
			"USD": 1.38,
			"EUR": 1.17,
			"JPY": 142.0,
			"TRY": 42.0,
			"GBP": 1.00,
		},
		"JPY": {
			"USD": 0.0091,
			"EUR": 0.0078,
			"GBP": 0.0070,
			"TRY": 0.32,
			"AED": 0.037,
			"JPY": 1.00,
		},
		"TRY": {
			"USD": 0.031,
			"EUR": 0.026,
			"GBP": 0.024,
			"JPY": 3.12,
			"AED": 0.27,
			"TRY": 1.00,
		},
		"AED": {
			"USD": 0.27,
			"EUR": 0.23,
			"GBP": 0.20,
			"JPY": 27.0,
			"TRY": 7.0,
			"AED": 1.00,
		},
	}

	if _, ok := rates[param.From]; !ok {
		return 0, fmt.Errorf(
			"currency not supported: %s",
			param.From,
		)
	}

	if rate, ok := rates[param.From][param.To]; ok {
		return rate, nil
	}

	return 0, fmt.Errorf(
		"conversion rate not found from %s to %s",
		param.From,
		param.To,
	)
}
