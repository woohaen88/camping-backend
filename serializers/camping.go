package serializers

import (
	"camping-backend/models"
	"time"
)

type Camping struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Address   string `json:"address"`
	CreatedAt time.Time
	UpdatedAt time.Time
	User      User
}

func CampingSerializer(camping models.Camping, user User) Camping {
	return Camping{
		ID:        camping.ID,
		Title:     camping.Title,
		Address:   camping.Address,
		CreatedAt: camping.CreatedAt,
		UpdatedAt: camping.UpdatedAt,
		User:      user,
	}
}
