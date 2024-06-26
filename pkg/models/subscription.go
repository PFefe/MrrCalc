package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// Subscription represents a subscription with dates parsed as time.Time
type Subscription struct {
	SubscriptionID string     `json:"subscription_id"`
	CustomerID     string     `json:"customer_id"`
	StartAt        time.Time  `json:"start_at"`
	EndAt          *time.Time `json:"end_at"`
	Amount         string     `json:"amount"`
	Currency       string     `json:"currency"`
	Interval       string     `json:"interval"`
	Status         string     `json:"status"`
	CancelledAt    *time.Time `json:"cancelled_at"`
}

// UnmarshalJSON custom unmarshaller to handle time fields
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
		return fmt.Errorf(
			"unmarshall error %w ",
			err,
		)
	}
	var err error
	s.StartAt, err = time.Parse(
		time.RFC3339,
		aux.StartAt,
	)
	if err != nil {
		return fmt.Errorf(
			"error parsing start_at %w",
			err,
		)
	}
	if aux.EndAt != nil {
		endAt, err := time.Parse(
			time.RFC3339,
			*aux.EndAt,
		)
		if err != nil {
			return fmt.Errorf(
				"error parsing end_at %w",
				err,
			)
		}
		s.EndAt = &endAt
	}
	if aux.CancelledAt != nil {
		cancelledAt, err := time.Parse(
			time.RFC3339,
			*aux.CancelledAt,
		)
		if err != nil {
			return fmt.Errorf(
				"error parsing cancelled_at %w",
				err,
			)
		}
		s.CancelledAt = &cancelledAt
	}
	return nil
}
