package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"MrrCalc/pkg/currencies"
	"MrrCalc/pkg/models"
	"MrrCalc/pkg/myTools"
	"MrrCalc/pkg/rates"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type mockRateProvider struct{}

func (m *mockRateProvider) GetRate(param *rates.RequestParameter) (float64, error) {
	if param.From == "USD" && param.To == "USD" {
		return 1.0, nil
	}
	return 0, fmt.Errorf(
		"unsupported currency pair: %s to %s",
		param.From,
		param.To,
	)
}

func TestCalculateMRR(t *testing.T) {
	t.Run(
		"check if presentMRR is calculated",
		func(t *testing.T) {
			t.Parallel()
			today := time.Now()
			firstDayOfMonth := time.Date(
				today.Year(),
				today.Month(),
				1,
				0,
				0,
				0,
				0,
				today.Location(),
			)
			formattedFirstDay := firstDayOfMonth.Format(time.RFC3339)
			twelveMonthsBeforeToday := today.AddDate(
				0,
				-12,
				0,
			)
			firstDayOfTwelveMonthsBeforeToday := time.Date(
				twelveMonthsBeforeToday.Year(),
				twelveMonthsBeforeToday.Month(),
				1,
				0,
				0,
				0,
				0,
				twelveMonthsBeforeToday.Location(),
			)
			formattedFirstDayOfTwelveMonthsBeforeToday := firstDayOfTwelveMonthsBeforeToday.Format(time.RFC3339)
			firstDayOfPreviousMonth := time.Date(
				today.Year(),
				today.Month()-1,
				2,
				0,
				0,
				0,
				0,
				today.Location(),
			)
			formattedFirstDayOfPreviousMonth := firstDayOfPreviousMonth.Format(time.RFC3339)

			subscriptionsJSON := fmt.Sprintf(
				`
		[
			{"subscription_id": "sub_001", "customer_id": "cust_001", "start_at": "%s", "end_at": null, "amount": "100.00", "currency": "USD", "interval": "month", "status": "active", "cancelled_at": null},
			{"subscription_id": "sub_002", "customer_id": "cust_002", "start_at": "%s", "end_at": null, "amount": "200.00", "currency": "USD", "interval": "month", "status": "active", "cancelled_at": null},
			{"subscription_id": "sub_003", "customer_id": "cust_003", "start_at": "%s", "end_at": null, "amount": "300.00", "currency": "USD", "interval": "month", "status": "active", "cancelled_at": null},
			{"subscription_id": "sub_004", "customer_id": "cust_004", "start_at": "%s", "end_at": null, "amount": "400.00", "currency": "USD", "interval": "month", "status": "cancelled", "cancelled_at": "%s"},
			{"subscription_id": "sub_005", "customer_id": "cust_004", "start_at": "%s", "end_at": null, "amount": "100.00", "currency": "USD", "interval": "month", "status": "active", "cancelled_at": null},
			{"subscription_id": "sub_006", "customer_id": "cust_005", "start_at": "%s", "end_at": null, "amount": "75.00", "currency": "USD", "interval": "month", "status": "amended", "cancelled_at": null},
			{"subscription_id": "sub_007", "customer_id": "cust_005", "start_at": "%s", "end_at": null, "amount": "100.00", "currency": "USD", "interval": "month", "status": "active", "cancelled_at": null},
			{"subscription_id": "sub_008", "customer_id": "cust_006", "start_at": "%s", "end_at": null, "amount": "100.00", "currency": "USD", "interval": "month", "status": "amended", "cancelled_at": null},
			{"subscription_id": "sub_009", "customer_id": "cust_006", "start_at": "%s", "end_at": null, "amount": "50.00", "currency": "USD", "interval": "month", "status": "active", "cancelled_at": null}
		]`,
				formattedFirstDayOfTwelveMonthsBeforeToday,
				formattedFirstDayOfTwelveMonthsBeforeToday,
				formattedFirstDay,
				formattedFirstDayOfTwelveMonthsBeforeToday,
				formattedFirstDayOfPreviousMonth,
				formattedFirstDay,
				formattedFirstDayOfTwelveMonthsBeforeToday,
				formattedFirstDay,
				formattedFirstDayOfTwelveMonthsBeforeToday,
				formattedFirstDay,
			)

			var subscriptions []models.Subscription
			err := json.Unmarshal(
				[]byte(subscriptionsJSON),
				&subscriptions,
			)
			if err != nil {
				t.Fatalf(
					"Failed to parse subscriptions JSON: %v",
					err,
				)
			}

			provider := &mockRateProvider{}
			converter := currencies.NewConverter(provider)

			expectedPresentMRR := decimal.NewFromFloat(850.00)
			expectedNewBusiness := decimal.NewFromFloat(300.00)
			expectedUpgrades := decimal.NewFromFloat(25.0)
			expectedDowngrades := decimal.NewFromFloat(50.0)
			expectedChurn := decimal.NewFromFloat(400.0)
			expectedReactivations := decimal.NewFromFloat(100.0)

			presentMRR, newBusiness, upgrades, downgrades, churn, reactivations := calculateMRR(
				converter,
				subscriptions,
				"USD",
			)
			assert.True(
				t,
				expectedPresentMRR.Equal(presentMRR),
			)
			assert.True(
				t,
				expectedNewBusiness.Equal(newBusiness),
			)
			assert.True(
				t,
				expectedUpgrades.Equal(upgrades),
			)
			assert.True(
				t,
				expectedDowngrades.Equal(downgrades),
			)
			assert.True(
				t,
				expectedChurn.Equal(churn),
			)
			assert.True(
				t,
				expectedReactivations.Equal(reactivations),
			)
		},
	)
}

func TestDailyMRR(t *testing.T) {
	t.Run(
		"check if dailyMRR is not null and not giving error",
		func(t *testing.T) {
			t.Parallel()
			subscriptions, _ := myTools.ReadJsonFileAndUnmarshall("subscriptions.json")
			provider := &mockRateProvider{}
			converter := currencies.NewConverter(provider)
			result := calculateDailyMRR(
				converter,
				subscriptions,
				"USD",
				1,
			)
			assert.NotNil(
				t,
				result,
			)
		},
	)
}
