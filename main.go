package main

import (
	"flag"
	"fmt"
	"log"

	"MrrCalc/pkg/myTools"
)

const (
	two = 2
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
		log.Printf(
			"Error reading and unmarshalling the file: %v\n",
			err,
		)
		return
	}
	presentMRR, newBusiness, upgrades, downgrades, churn, reactivations := calculateMRR(
		subscriptions,
		currency,
	)
	fmt.Printf(
		"Present MRR Net Value: %s %s\n",
		presentMRR.StringFixed(two),
		currency,
	)
	fmt.Println("Present MRR Breakdown:")
	fmt.Printf(
		"- New Business: %s %s\n",
		newBusiness.StringFixed(two),
		currency,
	)
	fmt.Printf(
		"- Upgrades: %s %s\n",
		upgrades.StringFixed(two),
		currency,
	)
	fmt.Printf(
		"- Downgrades: -%s %s\n",
		downgrades.StringFixed(two),
		currency,
	)
	fmt.Printf(
		"- Churn: -%s %s\n",
		churn.StringFixed(two),
		currency,
	)
	fmt.Printf(
		"- Reactivations: %s %s\n",
		reactivations.StringFixed(two),
		currency,
	)

	dailyMRRs := calculateDailyMRR(
		subscriptions,
		currency,
		period,
	)

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
			dailyMRR.MRR.StringFixed(two),
		)
	}
	fmt.Println("|------------|------------------|")
}
