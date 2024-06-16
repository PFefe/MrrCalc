package main

import (
	"MrrCalc/pkg/currencies"
	"MrrCalc/pkg/models"
	"MrrCalc/pkg/myTools"
	"github.com/shopspring/decimal"
	"time"
)

func calculateMRR(subscriptions []models.Subscription, currency string) (
	decimal.Decimal, decimal.Decimal, decimal.Decimal, decimal.Decimal, decimal.Decimal, decimal.Decimal) {
	todaysDate := time.Now()
	previousMonth := todaysDate.AddDate(
		0,
		-1,
		0,
	)
	lastDayOfPreviousMonth := time.Date(
		previousMonth.Year(),
		previousMonth.Month(),
		previousMonth.Day(),
		23,
		59,
		50,
		0,
		previousMonth.Location(),
	)
	firstDayOfPreviousMonth := time.Date(
		previousMonth.Year(),
		previousMonth.Month(),
		1,
		0,
		0,
		0,
		1,
		previousMonth.Location(),
	)
	presentMRR := decimal.NewFromFloat(.0)
	newBusiness := decimal.NewFromFloat(.0)
	upgrades := decimal.NewFromFloat(.0)
	downgrades := decimal.NewFromFloat(.0)
	churn := decimal.NewFromFloat(.0)
	reactivations := decimal.NewFromFloat(.0)
	previousAmounts := make(map[string]decimal.Decimal)
	previouslyCancelled := make(map[string]bool)
	previouslyAmended := make(map[string]bool)
	for i := 0; i < len(subscriptions); i++ {
		sub := subscriptions[i]
		status := sub.Status
		interval := sub.Interval
		isExpired := false

		if sub.EndAt != nil {
			isExpired = sub.EndAt.Before(todaysDate)
		} /*else {
			isExpired = false
		}*/

		amount, _ := myTools.ParseToFloat(sub.Amount)

		value, _ := currencies.ConvertCurrency(
			&currencies.Currency{
				From:   sub.Currency,
				To:     currency,
				Amount: amount,
			},
		)
		convertedAmount := decimal.NewFromFloat(
			value,
		)

		startedAt := sub.StartAt

		switch status {
		case "cancelled":
			if sub.CancelledAt != nil {
				cancelledAt := *sub.CancelledAt

				if cancelledAt.After(firstDayOfPreviousMonth) {
					if interval == "month" {
						churn = churn.Add(convertedAmount)
					} else if interval == "year" {
						churn = churn.Add(convertedAmount.Div(decimal.NewFromInt(12)))
					}
				}
				if cancelledAt.Before(lastDayOfPreviousMonth) {
					previouslyCancelled[sub.CustomerId] = true
				}
			}

		case "amended":
			if startedAt.Before(lastDayOfPreviousMonth) {
				previouslyAmended[sub.CustomerId] = true
				previousAmounts[sub.CustomerId] = convertedAmount
			}

		case "active":
			if !isExpired {
				// Calculate total MRR
				if startedAt.Before(todaysDate) {
					if interval == "month" {
						presentMRR = presentMRR.Add(convertedAmount)
					} else if interval == "year" {
						presentMRR = presentMRR.Add(convertedAmount.Div(decimal.NewFromInt(12)))
					}
				}
				if startedAt.After(lastDayOfPreviousMonth) {
					// Check if the customer has previously cancelled
					// Calculate reactivations
					if previouslyCancelled[sub.CustomerId] {
						if interval == "month" {
							reactivations = reactivations.Add(convertedAmount)
						} else if interval == "year" {
							reactivations = reactivations.Add(convertedAmount.Div(decimal.NewFromInt(12)))
						}
					}
					// Calculate new business
					if !previouslyCancelled[sub.CustomerId] && !previouslyAmended[sub.CustomerId] {
						if interval == "month" {
							newBusiness = newBusiness.Add(convertedAmount)
						} else if interval == "year" {
							newBusiness = newBusiness.Add(convertedAmount.Div(decimal.NewFromInt(12)))
						}
					}
					// Calculate upgrades and downgrades
					if previouslyAmended[sub.CustomerId] {
						prevAmount := previousAmounts[sub.CustomerId]
						if prevAmount.LessThan(convertedAmount) {
							if interval == "month" {
								upgrades = upgrades.Add(convertedAmount.Sub(prevAmount))
							} else if interval == "year" {
								upgrades = upgrades.Add(convertedAmount.Sub(prevAmount).Div(decimal.NewFromInt(12)))
							}
						}
						if prevAmount.GreaterThan(convertedAmount) {
							if interval == "month" {
								downgrades = downgrades.Add(prevAmount.Sub(convertedAmount))
							} else if interval == "year" {
								downgrades = downgrades.Add(prevAmount.Sub(convertedAmount).Div(decimal.NewFromInt(12)))
							}
						}
					}
				}
			}
		}
	}
	return presentMRR, newBusiness, upgrades, downgrades, churn, reactivations
}

type DailyMRR struct {
	Date string
	MRR  decimal.Decimal
}

func calculateDailyMRR(subscriptions []models.Subscription, currency string, period int) ([]DailyMRR, error) {
	todaysDate := time.Now()
	periodStartAt := time.Date(
		todaysDate.Year(),
		time.Month(period),
		1,
		0,
		0,
		0,
		0,
		todaysDate.Location(),
	)
	periodEndsAt := time.Date(
		todaysDate.Year(),
		time.Month(period+1),
		1,
		0,
		0,
		0,
		0,
		todaysDate.Location(),
	)

	var dailyMRRs []DailyMRR

	for periodStartAt.Before(periodEndsAt) {
		dailyMRR := decimal.NewFromFloat(.0)
		for i := 0; i < len(subscriptions); i++ {
			sub := subscriptions[i]
			status := sub.Status
			interval := sub.Interval
			isExpiredAt := false
			isCancelled := false

			if sub.EndAt != nil {
				isExpiredAt = sub.EndAt.Before(periodStartAt)
			}
			if status == "cancelled" && sub.CancelledAt != nil {
				isCancelled = sub.CancelledAt.Before(periodStartAt)
			}

			amount, _ := myTools.ParseToFloat(sub.Amount)

			value, _ := currencies.ConvertCurrency(
				&currencies.Currency{
					From:   sub.Currency,
					To:     currency,
					Amount: amount,
				},
			)
			convertedAmount := decimal.NewFromFloat(
				value,
			)

			startedAt := sub.StartAt

			if isExpiredAt == false && startedAt.Before(periodStartAt) && isCancelled == false {
				if interval == "month" {
					dailyMRR = dailyMRR.Add(convertedAmount)
				} else if interval == "year" {
					dailyMRR = dailyMRR.Add(convertedAmount.Div(decimal.NewFromInt(12)))
				}
			}
		}
		dailyMRRs = append(
			dailyMRRs,
			DailyMRR{
				Date: periodStartAt.Format("2006-01-02"),
				MRR:  dailyMRR,
			},
		)

		periodStartAt = periodStartAt.AddDate(
			0,
			0,
			1,
		)
	}
	return dailyMRRs, nil
}
