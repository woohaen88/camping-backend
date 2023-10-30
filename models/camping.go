package models

import "time"

type Camping struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Address   string `json:"address"`
	UserID    uint
	CreatedAt time.Time
	UpdatedAt time.Time
}
