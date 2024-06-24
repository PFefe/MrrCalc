// Package rates pkg/rates/rates.go
package rates

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type RequestParameter struct {
	From string
	To   string
}

type ExchangeRateProvider interface {
	GetRate(param *RequestParameter) (float64, error)
}

type APIService struct{}

type apiResponse struct {
	Rates map[string]map[string]float64 `json:"rates"`
}

func (s *APIService) GetRate(param *RequestParameter) (float64, error) {
	if param.From == param.To {
		return 1.0, nil
	}

	start := time.Now().Format("2006-01-02")
	end := time.Now().Format("2006-01-02")

	url := fmt.Sprintf(
		"https://api.frankfurter.app/%s..%s?from=%s&to=%s",
		start,
		end,
		param.From,
		param.To,
	)

	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf(
			"error making HTTP request: %w",
			err,
		)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing response body")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf(
			"received non-200 response code: %d",
			resp.StatusCode,
		)
	}

	var response apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, fmt.Errorf(
			"error decoding JSON response: %w",
			err,
		)
	}

	// Get the latest available rate
	for _, rates := range response.Rates {
		if rate, ok := rates[param.To]; ok {
			return rate, nil
		}
	}

	return 0, fmt.Errorf(
		"rate for %s not found",
		param.To,
	)

}
