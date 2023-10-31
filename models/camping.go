package models

import (
	"camping-backend/enums"
	"time"
)

type Camping struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	Title          string         `json:"title" gorm:"not null"`
	Address        string         `json:"address" gorm:"not null"`
	Description    string         `json:"description"`
	View           enums.ViewKind `json:"view" gorm:"type:view_kind"`                              // enum: Forest, Sea, Lake, Mountain, Other
	IsEvCharge     enums.Status   `json:"is_ev_charge" gorm:"type:Status;default:'CANT';not null"` // enum: OK, CANT, OTHER
	MannerTime     string         `json:"manner_time"`                                             // TODO string -> datetime
	IsSideParking  enums.Status   `json:"is_side_parking" gorm:"type:Status;not null"`             // enum: OK, CANT, OTHER
	IsPetFriendly  enums.Status   `json:"is_pet_friendly" gorm:"type:Status;not null"`             // enum: OK, CANT, OTHER
	VisitedStartAt string         `json:"visited_start_at"`                                        // TODO string -> datetime
	VisitedEndAt   string         `json:"visited_end_at"`

	UserID uint

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// TODO
	// Amenity, ManytoMany
	// Tag, ManytoMany

}
