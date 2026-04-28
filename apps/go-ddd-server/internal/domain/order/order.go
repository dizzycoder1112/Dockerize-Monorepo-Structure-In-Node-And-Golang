package order

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID           string    `json:"id"`
	CustomerName string    `json:"customer_name"`
	Amount       int64     `json:"amount"`
	Status       Status    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func New(customerName string, amount int64) (*Order, error) {
	if customerName == "" {
		return nil, ErrInvalidInput
	}
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	return &Order{
		ID:           id.String(),
		CustomerName: customerName,
		Amount:       amount,
		Status:       StatusPending,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func (o *Order) Pay() error {
	if o.Status != StatusPending {
		return ErrInvalidTransition
	}
	o.Status = StatusPaid
	o.UpdatedAt = time.Now()
	return nil
}

func (o *Order) Cancel() error {
	if o.Status == StatusPaid {
		return ErrInvalidTransition
	}
	o.Status = StatusCancelled
	o.UpdatedAt = time.Now()
	return nil
}

func (o *Order) ApplyUpdate(customerName string, amount int64) error {
	if customerName == "" {
		return ErrInvalidInput
	}
	if amount <= 0 {
		return ErrInvalidAmount
	}
	o.CustomerName = customerName
	o.Amount = amount
	o.UpdatedAt = time.Now()
	return nil
}
