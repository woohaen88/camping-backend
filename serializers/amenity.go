package serializers

import (
	"camping-backend/models"
	"time"
)

type Amenity struct {
	Id          uint
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	User        User
}

func AmenitySerializer(amenity models.Amenity, serializedUser User) Amenity {
	return Amenity{
		Id:          amenity.Id,
		Name:        amenity.Name,
		Description: amenity.Description,
		CreatedAt:   amenity.CreatedAt,
		UpdatedAt:   amenity.UpdatedAt,
		User:        serializedUser,
	}
}
