package myTools

import (
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

func ReadJsonFileAndUnmarshall(path string) ([]string, error) {
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
		return nil, fmt.Errorf(
			"unable to read the JSON file: %v",
			err,
		)
	}

	// initialize  Subscriptions array
	var output []string

	// unmarshal byteArray into 'subscriptions'
	err = json.Unmarshal(
		read,
		&output,
	)
	if err != nil {
		fmt.Printf(
			"Unable to unmarshal the json file %v",
			err,
		)
	}
	return output, nil
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
