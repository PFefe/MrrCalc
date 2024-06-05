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
		panic(err.Error())
	}

	return rate * amount
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

	// Open our jsonFile
	jsonFile, err := os.Open(input)
	// if the os.Open returns an error then handle it
	if err != nil {
		fmt.Printf(
			"Unable to read the file %v",
			err,
		)
	}
	fmt.Println("Json file read success")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	// read our opened jSon as a byte array.
	byteValue, _ := io.ReadAll(jsonFile)

	// we initialize our Subscriptions array
	var subscriptions []Subscription

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'subscriptions' which we defined above
	err = json.Unmarshal(
		byteValue,
		&subscriptions,
	)
	if err != nil {
		fmt.Printf(
			"Unable to unmarshal the json file %v",
			err,
		)
	}

	// we iterate through every subscription within our subscriptions array and
	// print out the fields

	// Calculate Present MRR Net Value:
	todaysDate := time.Now()
	presentMRR := 0.00
	for i := 0; i < len(subscriptions); i++ {
		status := subscriptions[i].Status
		interval := subscriptions[i].Interval
		isExpired := false
		if subscriptions[i].End_at == nil || *subscriptions[i].End_at == "" {
			isExpired = false
		} else {
			endAt, err := time.Parse(
				time.RFC3339,
				*subscriptions[i].End_at,
			)
			if err != nil {
				fmt.Printf(
					"Error parsing the end date for subscription %s: %v\n",
					subscriptions[i].Subscription_id,
					err,
				)
			}
			isExpired = endAt.Before(todaysDate)
		}

		amount, err := strconv.ParseFloat(
			subscriptions[i].Amount,
			64,
		)
		if err != nil {
			fmt.Printf(
				"Error parsing float: %v\n",
				err,
			)
		}

		convertedAmount := convertCurrency(
			amount,
			subscriptions[i].Currency,
			currency,
		)
		fmt.Printf(
			"Subscription ID: %s, Amount: %f, Converted Amount: %f, Status: %s, IsExpired: %t, Interval: %s\n",
			subscriptions[i].Subscription_id,
			amount,
			convertedAmount,
			status,
			isExpired,
			interval,
		)
		if (status == "active" || status == "amended") && !isExpired {
			if interval == "monthly" {
				presentMRR += convertedAmount
				fmt.Printf(
					"Added %f %s to presentMRR\n",
					convertedAmount,
					currency,
				)
			} else if interval == "yearly" {
				presentMRR += convertedAmount / 12
				fmt.Printf(
					"Added %f %s to presentMRR\n",
					convertedAmount/12,
					currency,
				)
			}
		}
	}

	for i := 0; i < len(subscriptions); i++ {
		fmt.Println(
			"subscription_id: ",
			subscriptions[i].Subscription_id,
		)
		fmt.Println(
			"customer_id: ",
			subscriptions[i].Customer_id,
		)
		fmt.Println(
			"start_at: ",
			subscriptions[i].Start_at,
		)
		fmt.Println(
			"end_at: ",
			subscriptions[i].End_at,
		)
		fmt.Println(
			"amount: ",
			subscriptions[i].Amount,
		)
		fmt.Println(
			"currency: ",
			subscriptions[i].Currency,
		)
		fmt.Println(
			"interval: ",
			subscriptions[i].Interval,
		)
		fmt.Println(
			"status: ",
			subscriptions[i].Status,
		)
		fmt.Println(
			"cancelled_at: ",
			subscriptions[i].Cancelled_at,
		)
	}
	fmt.Println(
		"currency: ",
		currency,
	)
	fmt.Println(
		"period: ",
		input,
	)

	// convert the currency test
	fmt.Println(
		"Converted currency: 100 usd to usd:",
		convertCurrency(
			100.00,
			"USD",
			"USD",
		),
	)

	// print: Present MRR Net Value: 2230.00 USD
	fmt.Println(
		"Present MRR Net Value: ",
		presentMRR,
		currency,
	)

}
