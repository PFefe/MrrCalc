package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/shopspring/decimal"
	exchangerates "github.com/yusufthedragon/exchange-rates-go"
	"io"
	"os"
	"strconv"
	"time"
)

var (
	currency string
	period   int
	input    string
)

// Subscription struct which contains required fields

type Subscription struct {
	Subscription_id string  `json:"subscription_id"`
	Customer_id     string  `json:"customer_id"`
	Start_at        string  `json:"start_at"`
	End_at          *string `json:"end_at"`
	Amount          string  `json:"amount"`
	Currency        string  `json:"currency"`
	Interval        string  `json:"interval"`
	Status          string  `json:"status"`
	Cancelled_at    *string `json:"cancelled_at"`
}

func convertCurrency(amount float64, from string, to string) (float64, error) {
	var rate, err = exchangerates.ConvertCurrency(
		&exchangerates.RequestParameter{
			From:  from,
			To:    to,
			Value: 1,
		},
	)

	if err != nil {
		fmt.Errorf(
			"Error converting currency: %v\n",
			err,
		)
		return 0, err
	}
	return rate * amount, nil
}

func parseToFloat(value string) (float64, error) {
	amount, err := strconv.ParseFloat(value, 64)
	if err != nil {
		fmt.Printf("Error parsing amount: %v\n", err)
		return 0, err
	}
	return amount, nil
}

func parseToTime(value string) (time.Time, error) {
	date, err := time.Parse(time.RFC3339, value)
	if err != nil {
		fmt.Printf("Error parsing time: %v\n", err)
		return time.Time{}, err // Return zero time on error
	}
	return date, nil
}

func readJsonFileAndUnmarshall(path string) ([]Subscription, error) {
	// Open our jsonFile
	jsonFile, err := os.Open(path)
	// if the os.Open returns an error then handle it
	if err != nil {
		fmt.Printf(
			"Unable to read the file %v",
			err,
		)
	}
	fmt.Println("Json file read success")
	// defer the closing jsonFile
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {

		}
	}(jsonFile)
	// read opened jSon
	read, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read the JSON file: %v", err)
	}

	// initialize  Subscriptions array
	var subscriptions []Subscription

	// unmarshal byteArray into 'subscriptions'
	err = json.Unmarshal(
		read,
		&subscriptions,
	)
	if err != nil {
		fmt.Printf(
			"Unable to unmarshal the json file %v",
			err,
		)
	}
	return subscriptions, nil
}

func calculateMRR(subscriptions []Subscription, currency string) (
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

		if sub.End_at != nil && *sub.End_at != "" {
			endAt, _ := parseToTime(*sub.End_at)
			isExpired = endAt.Before(todaysDate)
		}

		amount, _ := parseToFloat(sub.Amount)

		value, _ := convertCurrency(
			amount,
			sub.Currency,
			currency,
		)
		convertedAmount := decimal.NewFromFloat(
			value,
		)

		startedAt, _:= parseToTime(sub.Start_at)

		switch status {
		case "cancelled":
			if sub.Cancelled_at != nil {
				cancelledAt, _ := parseToTime(*sub.Cancelled_at)

				if cancelledAt.After(firstDayOfPreviousMonth) {
					if interval == "month" {
						churn = churn.Add(convertedAmount)
					} else if interval == "year" {
						churn = churn.Add(convertedAmount.Div(decimal.NewFromInt(12)))
					}
				}
				if cancelledAt.Before(lastDayOfPreviousMonth) {
					previouslyCancelled[sub.Customer_id] = true
				}
			}

		case "amended":
			if startedAt.Before(lastDayOfPreviousMonth) {
				previouslyAmended[sub.Customer_id] = true
				previousAmounts[sub.Customer_id] = convertedAmount
			}

		case "active":
			if isExpired == false {
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
					if previouslyCancelled[sub.Customer_id] {
						if interval == "month" {
							reactivations = reactivations.Add(convertedAmount)
						} else if interval == "year" {
							reactivations = reactivations.Add(convertedAmount.Div(decimal.NewFromInt(12)))
						}
					}
					// Calculate new business
					if !previouslyCancelled[sub.Customer_id] && !previouslyAmended[sub.Customer_id] {
						if interval == "month" {
							newBusiness = newBusiness.Add(convertedAmount)
						} else if interval == "year" {
							newBusiness = newBusiness.Add(convertedAmount.Div(decimal.NewFromInt(12)))
						}
					}
					// Calculate upgrades and downgrades
					if previouslyAmended[sub.Customer_id] {
						prevAmount := previousAmounts[sub.Customer_id]
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
	MRR  string
}

func calculateDailyMRR(subscriptions []Subscription, currency string, period int)([]DailyMRR, error){
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

			if sub.End_at != nil && *sub.End_at != "" {
				endAt, _ := parseToTime(*sub.End_at)
				isExpiredAt = endAt.Before(periodStartAt)
			}

			amount, _ := parseToFloat(sub.Amount)

			value, _ := convertCurrency(
				amount,
				sub.Currency,
				currency,
			)
			convertedAmount := decimal.NewFromFloat(
				value,
			)

			startedAt, _ := parseToTime(sub.Start_at)

			if (status == "active" || status == "amended") && isExpiredAt == false && startedAt.Before(periodStartAt) {
				if interval == "month" {
					dailyMRR = dailyMRR.Add(convertedAmount)
				} else if interval == "year" {
					dailyMRR = dailyMRR.Add(convertedAmount.Div(decimal.NewFromInt(12)))
				}
			}
		}
		dailyMRRs = append(dailyMRRs, DailyMRR{
			Date: periodStartAt.Format("2006-01-02"),
			MRR:  dailyMRR.StringFixed(2),
		})

		periodStartAt = periodStartAt.AddDate(
			0,
			0,
			1,
		)
	}
	// Printing the daily MRR values
	fmt.Println("\n Daily MRR:")
	fmt.Println("|------------|------------------|")
	fmt.Printf(
		"| Date       | MRR Value (%s)  |\n",
		currency,
	)
	fmt.Println("|------------|------------------|")
	for _, dailyMRR := range dailyMRRs {
		fmt.Printf(
			"| %s |     %s       |\n",
			dailyMRR.Date,
			dailyMRR.MRR,
		)
	}

	fmt.Println("|------------|------------------|")

	return dailyMRRs, nil

}

func main() {
	flag.StringVar(
		&currency,
		"currency",
		"USD",
		"Currency code",
	)
	flag.IntVar(
		&period,
		"period",
		1,
		"period",
	)
	flag.StringVar(
		&input,
		"input",
		"subscriptions.json",
		"Path to subscriptions json file",
	)
	flag.Parse()

	subscriptions, _ := readJsonFileAndUnmarshall(input)
	presentMRR, newBusiness, upgrades, downgrades, churn, reactivations := calculateMRR(subscriptions, currency)
	fmt.Printf(
		"Present MRR Net Value: %s %s\n",
		presentMRR.StringFixed(2),
		currency,
	)
	fmt.Println("Present MRR Breakdown:")
	fmt.Printf(
		"- New Business: %s %s\n",
		newBusiness.StringFixed(2),
		currency,
	)
	fmt.Printf(
		"- Upgrades: %s %s\n",
		upgrades.StringFixed(2),
		currency,
	)
	fmt.Printf(
		"- Downgrades: -%s %s\n",
		downgrades.StringFixed(2),
		currency,
	)
	fmt.Printf(
		"- Churn: -%s %s\n",
		churn.StringFixed(2),
		currency,
	)
	fmt.Printf(
		"- Reactivations: %s %s\n",
		reactivations.StringFixed(2),
		currency,
	)

	calculateDailyMRR(subscriptions, currency, period)
	gbpToUsd, _ := convertCurrency(1.00 , "GBP", "USD")
	fmt.Printf("GBP to USD rate: %.2f \n", gbpToUsd)
	eurToUsd, _ := convertCurrency(1.00, "EUR", "USD")
	fmt.Printf("EUR to USD rate: %.2f \n", eurToUsd)
}
