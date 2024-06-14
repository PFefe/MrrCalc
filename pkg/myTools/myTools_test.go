package myTools

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParseToFloat(t *testing.T) {
	t.Run(
		"check if any string is parsed into a float",
		func(t *testing.T) {
			// Test case 1
			value := "100.00"
			expectedResult := 100.00
			result, err := ParseToFloat(value)
			t.Logf(
				"Parsed value: %v, Error: %v",
				result,
				err,
			)
			assert.NoError(
				t,
				err,
			)
			assert.Equal(
				t,
				expectedResult,
				result,
			)
		},
	)

	t.Run(
		"check if invalid string returns an error",
		func(t *testing.T) {
			// Test case 2
			value := "invalid"
			result, err := ParseToFloat(value)
			t.Logf(
				"Parsed value: %v, Error: %v",
				result,
				err,
			)
			assert.Error(
				t,
				err,
			)
			assert.Equal(
				t,
				float64(0),
				result,
			) // Ensure that result is 0 when an error occurs
		},
	)
}

func TestParseToTime(t *testing.T) {
	t.Run(
		"check if valid RFC3339 string is parsed into a time.Time",
		func(t *testing.T) {
			// Test case 1
			value := "2024-06-12T12:00:00Z"
			expectedResult, _ := time.Parse(
				time.RFC3339,
				value,
			)
			result, err := ParseToTime(value)
			t.Logf(
				"Parsed value: %v, Error: %v",
				result,
				err,
			)
			assert.NoError(
				t,
				err,
			)
			assert.Equal(
				t,
				expectedResult,
				result,
			)
		},
	)

	t.Run(
		"check if invalid string returns an error",
		func(t *testing.T) {
			// Test case 2
			value := "invalid"
			result, err := ParseToTime(value)
			t.Logf(
				"Parsed value: %v, Error: %v",
				result,
				err,
			)
			assert.Error(
				t,
				err,
			)
			assert.True(
				t,
				result.IsZero(),
			) // Ensure that result is zero time when an error occurs
		},
	)
}
