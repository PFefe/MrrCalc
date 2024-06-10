package main

import (
	"encoding/json"
	"flag"
	"fmt"
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

func convertCurrency(amount float64, from string, to string) float64 {
	var rate, err = exchangerates.ConvertCurrency(
		&exchangerates.RequestParameter{
			From:  from,
			To:    to,
			Value: 1,
		},
	)

	if err != nil {
		panic(
			fmt.Errorf(
				"error converting currency: %v",
				err.Error,
			),
		)
	}
	return rate * amount
}

func parseToFloat(value string) float64 {
	amount, err := strconv.ParseFloat(
		value,
		64,
	)
	if err != nil {
		fmt.Printf(
			"Error parsing amount: %v\n",
			err,
		)
	}
	return amount
}

func parseToTime(value string) time.Time {
	date, err := time.Parse(
		time.RFC3339,
		value,
	)
	if err != nil {
		fmt.Printf(
			"Error parsing time: %v\n",
			err,
		)
	}
	return date
}

func readJsonFileAndUnmarshall(path string) []Subscription {
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
	read, _ := io.ReadAll(jsonFile)

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
	return subscriptions
}

func main() {
	flag.StringVar(
		&currency,
		"currency",
		"USD",
		"please enter the currency",
	)
	flag.IntVar(
		&period,
		"period",
		1,
		"please enter the period",
	)
	flag.StringVar(
		&input,
		"input",
		"subscriptions.json",
		"please enter the path to the json file",
	)
	flag.Parse()

	subscriptions := readJsonFileAndUnmarshall(input)

	// Calculate Present MRR Net Value:
	todaysDate := time.Now()
	presentMRR := 0.00
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

	newBusiness := 0.00
	upgrades := 0.00
	downgrades := 0.00
	churn := 0.00
	reactivations := 0.00
	previousAmounts := make(map[string]float64)
	previouslyCancelled := make(map[string]bool)
	previouslyAmended := make(map[string]bool)
	for i := 0; i < len(subscriptions); i++ {
		sub := subscriptions[i]
		status := sub.Status
		interval := sub.Interval
		isExpired := false

		if sub.End_at != nil && *sub.End_at != "" {
			endAt := parseToTime(*sub.End_at)
			isExpired = endAt.Before(todaysDate)
		}

		amount := parseToFloat(sub.Amount)

		convertedAmount := convertCurrency(
			amount,
			sub.Currency,
			currency,
		)

		startedAt := parseToTime(sub.Start_at)

		switch status {
		case "cancelled":
			if sub.Cancelled_at != nil {
				cancelledAt, err := time.Parse(
					time.RFC3339,
					*sub.Cancelled_at,
				)
				if err != nil {
					fmt.Printf(
						"Error parsing the cancelled date for subscription %s: %v\n",
						sub.Subscription_id,
						err,
					)
					continue
				}
				if cancelledAt.After(firstDayOfPreviousMonth) {
					if interval == "month" {
						churn += convertedAmount
					} else if interval == "year" {
						churn += convertedAmount / 12
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
						presentMRR += convertedAmount
					} else if interval == "year" {
						presentMRR += convertedAmount / 12
					}
				}
				if startedAt.After(lastDayOfPreviousMonth) {
					// Check if the customer has previously cancelled
					// Calculate reactivations
					if previouslyCancelled[sub.Customer_id] {
						if interval == "month" {
							reactivations += convertedAmount
						} else if interval == "year" {
							reactivations += convertedAmount / 12
						}
					}
					// Calculate new business
					if !previouslyCancelled[sub.Customer_id] {
						if interval == "month" {
							newBusiness += convertedAmount
						} else if interval == "year" {
							newBusiness += convertedAmount / 12
						}
					}
					// Calculate upgrades and downgrades
					if previouslyAmended[sub.Customer_id] {
						if previousAmounts[sub.Customer_id] < convertedAmount {
							if interval == "month" {
								upgrades += convertedAmount - previousAmounts[sub.Customer_id]
							} else if interval == "year" {
								upgrades += (convertedAmount - previousAmounts[sub.Customer_id]) / 12
							}
						}
						if previousAmounts[sub.Customer_id] > convertedAmount {
							if interval == "month" {
								downgrades += previousAmounts[sub.Customer_id] - convertedAmount
							} else if interval == "year" {
								downgrades += (previousAmounts[sub.Customer_id] - convertedAmount) / 12
							}
						}
					}
				}
			}
		}
	}
	fmt.Printf(
		"Present MRR Net Value: %.2f %s\n",
		presentMRR,
		currency,
	)
	fmt.Println("Present MRR Breakdown:")
	fmt.Printf(
		"- New Business: %.2f %s\n",
		newBusiness,
		currency,
	)
	fmt.Printf(
		"- Upgrades: %.2f %s\n",
		upgrades,
		currency,
	)
	fmt.Printf(
		"- Downgrades: -%.2f %s\n",
		downgrades,
		currency,
	)
	fmt.Printf(
		"- Churn: -%.2f %s\n",
		churn,
		currency,
	)
	fmt.Printf(
		"- Reactivations: %.2f %s\n",
		reactivations,
		currency,
	)

	// print daily MRR report
	// loop the days in a month/period
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

	fmt.Println("\n Daily MRR:")
	fmt.Println("|------------|------------------|")
	fmt.Printf(
		"| Date       | MRR Value (%s)  |",
		currency,
	)
	fmt.Println("\n|------------|------------------|")
	for periodStartAt.Before(periodEndsAt) {
		dailyMRR := 0.00
		for i := 0; i < len(subscriptions); i++ {
			sub := subscriptions[i]
			status := sub.Status
			interval := sub.Interval
			isExpiredAt := false

			if sub.End_at != nil && *sub.End_at != "" {
				endAt := parseToTime(*sub.End_at)
				isExpiredAt = endAt.Before(periodStartAt)
			}

			amount := parseToFloat(sub.Amount)

			convertedAmount := convertCurrency(
				amount,
				sub.Currency,
				currency,
			)

			startedAt := parseToTime(sub.Start_at)

			if (status == "active" || status == "amended") && isExpiredAt == false && startedAt.Before(periodStartAt) {
				if interval == "month" {
					dailyMRR += convertedAmount
				} else if interval == "year" {
					dailyMRR += convertedAmount / 12
				}
			}
		}
		fmt.Printf(
			"| %s |     %.2f       | \n",
			periodStartAt.Format("2006-01-02"),
			dailyMRR,
		)
		periodStartAt = periodStartAt.AddDate(
			0,
			0,
			1,
		)
	}
	fmt.Println("|------------|------------------|")
}
