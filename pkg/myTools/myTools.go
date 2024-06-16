package myTools

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"MrrCalc/pkg/models"
)

func ParseToFloat(value string) (float64, error) {
	amount, err := strconv.ParseFloat(
		value,
		64,
	)
	if err != nil {
		log.Printf(
			"Error parsing amount: %v\n",
			err,
		)
		return 0, err
	}
	return amount, nil
}

func ParseToTime(value string) (time.Time, error) {
	date, err := time.Parse(
		time.RFC3339,
		value,
	)
	if err != nil {
		log.Printf(
			"Error parsing time: %v\n",
			err,
		)
		return time.Time{}, err // Return zero time on error
	}
	return date, nil
}

func ReadJsonFileAndUnmarshall(path string) ([]models.Subscription, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to read the file: %w",
			err,
		)
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			fmt.Println(
				"Unable to close : ",
				path,
			)
		}
	}(jsonFile)

	fmt.Println("Json file read success")
	read, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to read the JSON file: %w",
			err,
		)
	}

	var output []models.Subscription
	err = json.Unmarshal(
		read,
		&output,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal the json file: %w",
			err,
		)
	}
	return output, nil
}
