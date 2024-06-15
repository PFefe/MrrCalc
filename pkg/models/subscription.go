package models

import (
	"encoding/json"
	"time"
)

// Subscription represents a subscription with dates parsed as time.Time
type Subscription struct {
	SubscriptionId string     `json:"subscription_id"`
	CustomerId     string     `json:"customer_id"`
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
}
