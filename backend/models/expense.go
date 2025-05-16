package models

import "time"

type Expense struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Category  string    `json:"category"`
	Amount    float64   `json:"amount"`
	Note      string    `json:"note"`
	SpentAt   time.Time `json:"spent_at"`
	CreatedAt time.Time `json:"created_at"`
}
