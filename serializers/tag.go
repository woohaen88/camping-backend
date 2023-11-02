package serializers

import (
	"camping-backend/models"
	"time"
)

type Tag struct {
	ID          uint
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UserId      uint
	User        User
}

func TagSerializer(tag models.Tag, serializedUser User) Tag {
	return Tag{
		ID:          tag.ID,
		Name:        tag.Name,
		Description: tag.Description,
		CreatedAt:   tag.CreatedAt,
		UpdatedAt:   tag.UpdatedAt,
		UserId:      tag.UserId,
		User:        serializedUser,
	}
}
