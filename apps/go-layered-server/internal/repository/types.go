package repository

import "time"

// DealRow is the normalized deal data returned by the repository layer.
// Repositories return their own row types — separate from request/response DTOs
// — so handlers can serialize them directly without the service layer needing
// to copy fields.
type DealRow struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Amount    int64     `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
