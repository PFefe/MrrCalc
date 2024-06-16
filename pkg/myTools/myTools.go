package myTools

import (
	"MrrCalc/pkg/models"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

func ParseToFloat(value string) (float64, error) {
	amount, err := strconv.ParseFloat(
		value,
		64,
	)
	if err != nil {
		fmt.Printf(
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
		fmt.Printf(
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
			"unable to read the file: %v",
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
			"unable to read the JSON file: %v",
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
			"unable to unmarshal the json file: %v",
			err,
		)
	}
	return output, nil
}
