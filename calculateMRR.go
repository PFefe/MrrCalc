package main

import (
	"time"

	"MrrCalc/pkg/currencies"
	"MrrCalc/pkg/models"
	"MrrCalc/pkg/myTools"

	"github.com/shopspring/decimal"
)

const (
	zero          = 0.0
	months        = 12
	YearInterval  = "year"
	MonthInterval = "month"
)

func calculateMRR(subscriptions []models.Subscription, currency string) (
	decimal.Decimal, decimal.Decimal, decimal.Decimal, decimal.Decimal, decimal.Decimal, decimal.Decimal,
) {
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
	presentMRR := decimal.NewFromFloat(zero)
	newBusiness := decimal.NewFromFloat(zero)
	upgrades := decimal.NewFromFloat(zero)
	downgrades := decimal.NewFromFloat(zero)
	churn := decimal.NewFromFloat(zero)
	reactivations := decimal.NewFromFloat(zero)
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
					if interval == MonthInterval {
						churn = churn.Add(convertedAmount)
					} else if interval == YearInterval {
						churn = churn.Add(convertedAmount.Div(decimal.NewFromInt(months)))
					}
				}
				if cancelledAt.Before(lastDayOfPreviousMonth) {
					previouslyCancelled[sub.CustomerID] = true
				}
			}

		case "amended":
			if startedAt.Before(lastDayOfPreviousMonth) {
				previouslyAmended[sub.CustomerID] = true
				previousAmounts[sub.CustomerID] = convertedAmount
			}

		case "active":
			if !isExpired {
				// Calculate total MRR
				if startedAt.Before(todaysDate) {
					if interval == MonthInterval {
						presentMRR = presentMRR.Add(convertedAmount)
					} else if interval == YearInterval {
						presentMRR = presentMRR.Add(convertedAmount.Div(decimal.NewFromInt(months)))
					}
				}
				if startedAt.After(lastDayOfPreviousMonth) {
					// Check if the customer has previously cancelled
					// Calculate reactivations
					if previouslyCancelled[sub.CustomerID] {
						if interval == MonthInterval {
							reactivations = reactivations.Add(convertedAmount)
						} else if interval == YearInterval {
							reactivations = reactivations.Add(convertedAmount.Div(decimal.NewFromInt(months)))
						}
					}
					// Calculate new business
					if !previouslyCancelled[sub.CustomerID] && !previouslyAmended[sub.CustomerID] {
						if interval == MonthInterval {
							newBusiness = newBusiness.Add(convertedAmount)
						} else if interval == YearInterval {
							newBusiness = newBusiness.Add(convertedAmount.Div(decimal.NewFromInt(months)))
						}
					}
					// Calculate upgrades and downgrades
					if previouslyAmended[sub.CustomerID] {
						prevAmount := previousAmounts[sub.CustomerID]
						if prevAmount.LessThan(convertedAmount) {
							if interval == MonthInterval {
								upgrades = upgrades.Add(convertedAmount.Sub(prevAmount))
							} else if interval == YearInterval {
								upgrades = upgrades.Add(convertedAmount.Sub(prevAmount).Div(decimal.NewFromInt(months)))
							}
						}
						if prevAmount.GreaterThan(convertedAmount) {
							if interval == MonthInterval {
								downgrades = downgrades.Add(prevAmount.Sub(convertedAmount))
							} else if interval == YearInterval {
								downgrades = downgrades.Add(prevAmount.Sub(convertedAmount).Div(decimal.NewFromInt(months)))
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

func calculateDailyMRR(subscriptions []models.Subscription, currency string, period int) []DailyMRR {
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
		dailyMRR := decimal.NewFromFloat(zero)
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

			if !isExpiredAt && startedAt.Before(periodStartAt) && !isCancelled {
				if interval == MonthInterval {
					dailyMRR = dailyMRR.Add(convertedAmount)
				} else if interval == YearInterval {
					dailyMRR = dailyMRR.Add(convertedAmount.Div(decimal.NewFromInt(months)))
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
	return dailyMRRs
}
