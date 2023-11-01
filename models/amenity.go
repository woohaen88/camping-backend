package models

import (
	"gorm.io/gorm"
	"time"
)

type Amenity struct {
	gorm.Model
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Campings    []*Camping `json:"campings" gorm:"many2many:camping_amenities"`
}
