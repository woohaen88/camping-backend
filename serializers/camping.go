package serializers

import (
	"camping-backend/enums"
	"camping-backend/models"
	"time"
)

type Camping struct {
	ID             uint
	Title          string
	Address        string
	Description    string
	View           enums.ViewKind
	IsEvCharge     enums.Status
	MannerTime     string // TODO: parseDateTIme
	IsSideParking  enums.Status
	IsPetFriendly  enums.Status
	VisitedStartAt string // TODO: parseDateTIme
	VisitedEndAt   string // TODO: parseDateTIme

	UserID uint
	User   User // Serializer

	CreatedAt time.Time
	UpdatedAt time.Time
}

func CampingSerializer(camping *models.Camping, serializeUser User) Camping {
	return Camping{
		ID:             camping.ID,
		Title:          camping.Title,
		Address:        camping.Address,
		Description:    camping.Description,
		View:           camping.View,
		IsEvCharge:     camping.IsEvCharge,
		MannerTime:     camping.MannerTime,
		IsSideParking:  camping.IsSideParking,
		IsPetFriendly:  camping.IsPetFriendly,
		VisitedStartAt: camping.VisitedStartAt,
		VisitedEndAt:   camping.VisitedEndAt,

		UserID: serializeUser.ID,
		User:   serializeUser,

		CreatedAt: camping.CreatedAt,
		UpdatedAt: camping.UpdatedAt,
	}
}
