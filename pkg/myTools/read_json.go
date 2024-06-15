package myTools

import (
	"MrrCalc/pkg/models"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func ReadJsonFileAndUnmarshall(path string) ([]models.Subscription, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to read the file: %v",
			err,
		)
	}
	defer jsonFile.Close()

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
