package myTools

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseToFloat(t *testing.T) {
	t.Run(
		"check if any string is parsed into a float",
		func(t *testing.T) {
			t.Parallel()
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
			t.Parallel()
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
			t.Parallel()
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
			t.Parallel()
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

func TestReadJsonFileAndUnmarshall(t *testing.T) {
	t.Run(
		"check if JSON file is read and unmarshalled correctly",
		func(t *testing.T) {
			t.Parallel()
			// Assuming you have a test JSON file `test_subscriptions.json` in the same directory
			subscriptions, err := ReadJsonFileAndUnmarshall("../../testdata/subscriptions.json")
			t.Logf(
				"Subscriptions: %v, Error: %v",
				subscriptions,
				err,
			)
			assert.NoError(
				t,
				err,
			)
			assert.NotEmpty(
				t,
				subscriptions,
			)
		},
	)

	t.Run(
		"check if non-existent file returns an error",
		func(t *testing.T) {
			t.Parallel()
			_, err := ReadJsonFileAndUnmarshall("non_existent_file.json")
			t.Logf(
				"Error: %v",
				err,
			)
			assert.Error(
				t,
				err,
			)
		},
	)
}
