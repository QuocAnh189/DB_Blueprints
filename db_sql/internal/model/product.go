package model

import "time"

type Product struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	OwnerID   int64     `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Owner     *User     `json:"user"`
}
