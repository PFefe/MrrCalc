package main


import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"fmt"
)


func TestConvertCurrency(t *testing.T) {
	t.Run("check if currency is converted", func (t *testing.T) {
	// Test case 1
		amount := 100.00
		from := "EUR"
		to := "USD"
		result, err := convertCurrency(amount, from, to)
		assert.NoError(t, err)
		assert.NotEqual(t, amount, result)
		})


	t.Run("check if same currency is converted with same amount", func (t *testing.T) {
		// Test case 1
		amount := 100.00
		from := "EUR"
		to := "EUR"
		result, err := convertCurrency(amount, from, to)
		assert.NoError(t, err)
		assert.Equal(t, amount, result)
		})
}

func TestParseToFloat(t *testing.T) {
	t.Run("check if any string is parsed into a float", func(t *testing.T) {
		// Test case 1
		value := "100.00"
		expectedResult := 100.00
		result, err := parseToFloat(value)
		t.Logf("Parsed value: %v, Error: %v", result, err)
		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("check if invalid string returns an error", func(t *testing.T) {
		// Test case 2
		value := "invalid"
		result, err := parseToFloat(value)
		t.Logf("Parsed value: %v, Error: %v", result, err)
		assert.Error(t, err)
		assert.Equal(t, float64(0), result) // Ensure that result is 0 when an error occurs
	})
}

func TestParseToTime(t *testing.T) {
	t.Run("check if valid RFC3339 string is parsed into a time.Time", func(t *testing.T) {
		// Test case 1
		value := "2024-06-12T12:00:00Z"
		expectedResult, _ := time.Parse(time.RFC3339, value)
		result, err := parseToTime(value)
		t.Logf("Parsed value: %v, Error: %v", result, err)
		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("check if invalid string returns an error", func(t *testing.T) {
		// Test case 2
		value := "invalid"
		result, err := parseToTime(value)
		t.Logf("Parsed value: %v, Error: %v", result, err)
		assert.Error(t, err)
		assert.True(t, result.IsZero()) // Ensure that result is zero time when an error occurs
	})
}

func TestReadJsonFileAndUnmarshall(t *testing.T) {
	t.Run("check if JSON file is read and unmarshalled correctly", func(t *testing.T) {
		// Assuming you have a test JSON file `test_subscriptions.json` in the same directory
		subscriptions, err := readJsonFileAndUnmarshall("subscriptions.json")
		t.Logf("Subscriptions: %v, Error: %v", subscriptions, err)
		assert.NoError(t, err)
		assert.NotEmpty(t, subscriptions)
	})

	t.Run("check if non-existent file returns an error", func(t *testing.T) {
		_, err := readJsonFileAndUnmarshall("non_existent_file.json")
		t.Logf("Error: %v", err)
		assert.Error(t, err)
	})
}

func TestCalculateMRR(t *testing.T) {
	t.Run("check if presentMRR is calculated", func(t *testing.T) {
		// Test case 1
		today := time.Now()
		firstDayOfMonth := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
		formattedFirstDay := firstDayOfMonth.Format(time.RFC3339)
		twelveMonthsBeforeToday := today.AddDate(0, -12, 0)
		firstDayOfTwelveMonthsBeforeToday := time.Date(twelveMonthsBeforeToday.Year(), twelveMonthsBeforeToday.Month(), 1, 0, 0, 0, 0, twelveMonthsBeforeToday.Location())
		formattedFirstDayOfTwelveMonthsBeforeToday := firstDayOfTwelveMonthsBeforeToday.Format(time.RFC3339)
		firstDayOfPreviousMonth := time.Date(today.Year(), today.Month()-1, 2, 0, 0, 0, 0, today.Location())
		formattedFirstDayOfPreviousMonth := firstDayOfPreviousMonth.Format(time.RFC3339)
		// Parse subscriptions JSON
		subscriptionsJSON := fmt.Sprintf(`
		[
			{
				"subscription_id": "sub_001",
				"customer_id": "cust_001",
				"start_at": "%s",
				"end_at": null,
				"amount": "100.00",
				"currency": "USD",
				"interval": "month",
				"status": "active",
				"cancelled_at": null
			},
			{
				"subscription_id": "sub_002",
				"customer_id": "cust_002",
				"start_at": "%s",
				"end_at": null,
				"amount": "200.00",
				"currency": "USD",
				"interval": "month",
				"status": "active",
				"cancelled_at": null
			},
			{
				"subscription_id": "sub_003",
				"customer_id": "cust_003",
				"start_at": "%s",
				"end_at": null,
				"amount": "300.00",
				"currency": "USD",
				"interval": "month",
				"status": "active",
				"cancelled_at": null
			},
			{
				"subscription_id": "sub_004",
				"customer_id": "cust_004",
				"start_at": "%s",
				"end_at": null,
				"amount": "400.00",
				"currency": "USD",
				"interval": "month",
				"status": "cancelled",
				"cancelled_at": "%s"
			},
			{
				"subscription_id": "sub_005",
				"customer_id": "cust_004",
				"start_at": "%s",
				"end_at": null,
				"amount": "100.00",
				"currency": "USD",
				"interval": "month",
				"status": "active",
				"cancelled_at": null
			},
  			{
  			  "subscription_id": "sub_006",
  			  "customer_id": "cust_005",
  			  "start_at": "%s",
  			  "end_at": null,
  			  "amount": "75.00",
  			  "currency": "USD",
  			  "interval": "month",
  			  "status": "amended",
  			  "cancelled_at": null
  			},
  			{
  			  "subscription_id": "sub_007",
  			  "customer_id": "cust_005",
  			  "start_at": "%s",
  			  "end_at": null,
  			  "amount": "100.00",
  			  "currency": "USD",
  			  "interval": "month",
  			  "status": "active",
  			  "cancelled_at": null
  			},
  			{
  			  "subscription_id": "sub_008",
  			  "customer_id": "cust_006",
  			  "start_at": "%s",
  			  "end_at": null,
  			  "amount": "100.00",
  			  "currency": "USD",
  			  "interval": "month",
  			  "status": "amended",
  			  "cancelled_at": null
  			},
  			{
  			  "subscription_id": "sub_09",
  			  "customer_id": "cust_006",
  			  "start_at": "%s",
  			  "end_at": null,
  			  "amount": "50.00",
  			  "currency": "USD",
  			  "interval": "month",
  			  "status": "active",
  			  "cancelled_at": null
  			}

		]`, formattedFirstDayOfTwelveMonthsBeforeToday, formattedFirstDayOfTwelveMonthsBeforeToday, formattedFirstDay, formattedFirstDayOfTwelveMonthsBeforeToday, formattedFirstDayOfPreviousMonth, formattedFirstDay, formattedFirstDayOfTwelveMonthsBeforeToday, formattedFirstDay, formattedFirstDayOfTwelveMonthsBeforeToday, formattedFirstDay,)

		var subscriptions []Subscription
		err := json.Unmarshal([]byte(subscriptionsJSON), &subscriptions)
		if err != nil {
			t.Fatalf("Failed to parse subscriptions JSON: %v", err)
		}

		// Test case 1
		expectedPresentMRR := decimal.NewFromFloat(850.00)
		expectedNewBusiness := decimal.NewFromFloat(300.00)
		expectedUpgrades := decimal.NewFromFloat(25.0)
		expectedDowngrades := decimal.NewFromFloat(50.0)
		expectedChurn := decimal.NewFromFloat(400.0)
		expectedReactivations := decimal.NewFromFloat(100.0)
		presentMRR, newBusiness, upgrades, downgrades, churn, reactivations := calculateMRR(subscriptions, "USD")
		assert.True(t, expectedPresentMRR.Equals(presentMRR))
		assert.True(t, expectedNewBusiness.Equals(newBusiness))
		assert.True(t, expectedUpgrades.Equals(upgrades))
		assert.True(t, expectedDowngrades.Equals(downgrades))
		assert.True(t, expectedChurn.Equals(churn))
		assert.True(t, expectedReactivations.Equals(reactivations))
	})
}
func TestDailyMRR(t *testing.T) {
	t.Run("check if dailyMRR is not null and not giving error", func(t *testing.T) {
		subscriptions, _ := readJsonFileAndUnmarshall("subscriptions.json")
		result, err := calculateDailyMRR(subscriptions, "USD", 1)
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}