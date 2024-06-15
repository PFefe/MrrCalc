package main

import (
	"MrrCalc/pkg/currencies"
	"MrrCalc/pkg/myTools"
	"flag"
	"fmt"
)

var (
	currency string
	period   int
	input    string
)

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

	subscriptions, err := myTools.ReadJsonFileAndUnmarshall(input)
	if err != nil {
		fmt.Printf(
			"Error reading and unmarshalling the file: %v\n",
			err,
		)
		return
	}
	/*	fmt.Println("Json file read success")
		for _, sub := range subscriptions {
			fmt.Printf(
				"%+v\n",
				sub,
			)
		}*/
	presentMRR, newBusiness, upgrades, downgrades, churn, reactivations := calculateMRR(
		subscriptions,
		currency,
	)
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

	dailyMRRs, err := calculateDailyMRR(
		subscriptions,
		currency,
		period,
	)
	if err != nil {
		fmt.Printf(
			"Error calculating daily MRR: %v\n",
			err,
		)
		return
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
			dailyMRR.MRR.StringFixed(2),
		)
	}
	fmt.Println("|------------|------------------|")

	gbpToUsd, _ := currencies.ConvertCurrency(
		&currencies.Currency{
			From:   "GBP",
			To:     "USD",
			Amount: 1.00,
		},
	)
	fmt.Printf(
		"GBP to USD rate: %.2f \n",
		gbpToUsd,
	)
	eurToUsd, _ := currencies.ConvertCurrency(
		&currencies.Currency{
			From:   "EUR",
			To:     "USD",
			Amount: 1.00,
		},
	)
	fmt.Printf(
		"EUR to USD rate: %.2f \n",
		eurToUsd,
	)
}
