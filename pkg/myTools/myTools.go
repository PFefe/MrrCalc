package myTools

import (
	"fmt"
	"strconv"
	"time"
)

/*func ReadJsonFileAndUnmarshall(path string) ([]models.Subscription, error) {
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

// Custom unmarshaller for the Subscription struct to handle time fields
func (s *Subscription) UnmarshalJSON(data []byte) error {
	type Alias Subscription
	aux := &struct {
		StartAt     string  `json:"start_at"`
		EndAt       *string `json:"end_at"`
		CancelledAt *string `json:"cancelled_at"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(
		data,
		&aux,
	); err != nil {
		return err
	}
	var err error
	s.StartAt, err = time.Parse(
		time.RFC3339,
		aux.StartAt,
	)
	if err != nil {
		return err
	}
	if aux.EndAt != nil {
		endAt, err := time.Parse(
			time.RFC3339,
			*aux.EndAt,
		)
		if err != nil {
			return err
		}
		s.EndAt = &endAt
	}
	if aux.CancelledAt != nil {
		cancelledAt, err := time.Parse(
			time.RFC3339,
			*aux.CancelledAt,
		)
		if err != nil {
			return err
		}
		s.CancelledAt = &cancelledAt
	}
	return nil
}*/

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
