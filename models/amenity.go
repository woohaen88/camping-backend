package models

import (
	"time"
)

type Amenity struct {
	Id          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"not null;unique"`
	Description string `json:"description"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Campings    []*Camping `json:"campings" gorm:"many2many:camping_amenities"`
	UserId      uint       `json:"user_id" gorm:"not null"`
}
