package models

import (
	"gorm.io/gorm"
	"time"
)

type Tag struct {
	gorm.Model
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"not null;unique"`
	Description string `json:"description"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Campings    []*Camping `json:"campings" gorm:"many2many:camping_tags;"`
	UserId      uint
}
