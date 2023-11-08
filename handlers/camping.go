package handlers

import (
	commonErrors "camping-backend/common/errors"
	"camping-backend/database"
	"camping-backend/enums"
	"camping-backend/middleware"
	"camping-backend/models"
	"camping-backend/serializers"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/structs"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"log"
)

func CreateCamping(c *fiber.Ctx) error {

	user, err := middleware.GetAuthUser(c)

	if err != nil {
		return commonErrors.ErrorHandler(c, fiber.StatusUnauthorized, err)
	}

	var camping = models.Camping{
		CreatedAt: database.Database.Conn.NowFunc(),
		UpdatedAt: database.Database.Conn.NowFunc(),
		UserId:    user.ID,
	}
	request := struct {
		Title          string         `json:"title"`
		Address        string         `json:"address"`
		Description    string         `json:"description"`
		View           enums.ViewKind `json:"view"`
		IsEvCharge     enums.Status   `json:"is_ev_charge"`
		MannerTime     string         `json:"manner_time"`
		IsSideParking  enums.Status   `json:"is_side_parking"`
		IsPetFriendly  enums.Status   `json:"is_pet_friendly"`
		VisitedStartAt string         `json:"visited_start_at"`
		VisitedEndAt   string         `json:"visited_end_at"`
		Tags           []uint         `json:"tags"`
		Amenities      []uint         `json:"amenities"`
	}{}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn't parse createCamping",
			"data":    err.Error(),
		})
	}

	// enum check
	if err := setEnumView(request.View); err != nil {
		return commonErrors.ErrorHandler(c, fiber.StatusBadRequest, err)
	}

	if err := setEnumStatus(request.IsEvCharge, request.IsSideParking, request.IsPetFriendly); err != nil {
		return commonErrors.ErrorHandler(c, fiber.StatusBadRequest, err)
	}

	// requestt -> camping
	camping.Title = request.Title
	camping.Address = request.Address
	camping.Description = request.Description
	camping.View = request.View
	camping.IsEvCharge = request.IsEvCharge
	camping.MannerTime = request.MannerTime
	camping.IsSideParking = request.IsSideParking
	camping.IsPetFriendly = request.IsPetFriendly
	camping.VisitedStartAt = request.VisitedStartAt
	camping.VisitedEndAt = request.VisitedEndAt

	var tagModels []models.Tag

	for _, tagId := range request.Tags {
		var tagModel models.Tag
		err = database.Database.Conn.First(&tagModel, "id = ?", tagId).Error
		if err != nil {
			return commonErrors.ErrorHandler(c, fiber.StatusBadRequest, err, "could not parse tag")
		}

		tagModels = append(tagModels, tagModel)
	}

	var serializedTags []serializers.Tag
	for i := 0; i < len(tagModels); i++ {
		tagModel := tagModels[i]
		var tagUser models.User
		err = FindUserById(&tagUser, int(tagModel.UserId))
		if err != nil {
			return commonErrors.ErrorHandler(c, fiber.StatusNotFound, err)
		}
		serializedTag := serializers.TagSerializer(tagModel, serializers.UserSerializer(user))
		serializedTags = append(serializedTags, serializedTag)
	}

	// amenity
	var amenities []models.Amenity
	var serializedAmenities []serializers.Amenity
	for _, amenityId := range request.Amenities {
		amenity := models.Amenity{}
		var amenityCreator models.User
		if err := database.Database.Conn.First(&amenity, "id = ?", amenityId).Error; err != nil {
			return commonErrors.ErrorHandler(c, fiber.StatusNotFound, err)
		}

		if err := FindUserById(&amenityCreator, int(amenity.UserId)); err != nil {
			return commonErrors.ErrorHandler(c, fiber.StatusNotFound, err)
		}

		serializedAmenity := serializers.AmenitySerializer(amenity, serializers.UserSerializer(&amenityCreator))
		serializedAmenities = append(serializedAmenities, serializedAmenity)
		amenities = append(amenities, amenity)
	}

	database.Database.Conn.Create(&camping)

	if err := database.Database.Conn.Model(&camping).Association("Tags").Append(tagModels); err != nil {
		return commonErrors.ErrorHandler(c, fiber.StatusBadRequest, err, "failed many to many field camping=>tags")
	}

	if err := database.Database.Conn.Model(&camping).Association("Amenities").Append(amenities); err != nil {
		return commonErrors.ErrorHandler(c, fiber.StatusBadRequest, err, "failed many to many field camping=>amenities")
	}
	responseUser := serializers.UserSerializer(user)
	responseCamping := serializers.CampingSerializer(&camping, responseUser, serializedTags, serializedAmenities)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "success",
		"data":    responseCamping,
	})

}

func ListCamping(c *fiber.Ctx) error {
	var campings []models.Camping
	var owner models.User

	database.Database.Conn.Find(&campings)

	var responseCampings []serializers.Camping

	for _, camping := range campings {
		if err := FindUserById(&owner, int(camping.UserId)); err != nil {
			return commonErrors.ErrorHandler(c, fiber.StatusBadRequest, err)
		}
		responseUser := serializers.UserSerializer(&owner)
		responseCamping := serializers.CampingSerializer(&camping, responseUser)
		responseCampings = append(responseCampings, responseCamping)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "success",
		"data":    responseCampings,
	})
}

func GetCamping(c *fiber.Ctx) error {
	campingId, err := c.ParamsInt("campingId")
	if err != nil {
		return commonErrors.ErrorHandler(c, fiber.StatusBadRequest, err)
	}

	var camping models.Camping
	var user models.User

	if err := database.Database.Conn.First(&camping, "id = ?", campingId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return commonErrors.ErrorHandler(c, fiber.StatusNotFound, err)
		}
		return commonErrors.ErrorHandler(c, fiber.StatusBadRequest, err)
	}

	if err := FindUserById(&user, int(camping.UserId)); err != nil {
		return commonErrors.ErrorHandler(c, fiber.StatusNotFound, err)
	}

	var tags []*models.Tag
	if err := database.Database.Conn.Model(&camping).Association("Tags").Find(&tags); err != nil {
		return commonErrors.ErrorHandler(c, fiber.StatusBadRequest, err)
	}

	var amenities []*models.Amenity
	if err := database.Database.Conn.Model(&camping).Association("Amenities").Find(&amenities); err != nil {
		return commonErrors.ErrorHandler(c, fiber.StatusBadRequest, err)
	}

	camping.Tags = tags
	camping.Amenities = amenities

	var serializedTags []serializers.Tag
	for _, tag := range camping.Tags {
		tag, err := FindByTagId(tag.ID)
		if err != nil {
			return commonErrors.ErrorHandler(c, fiber.StatusNotFound, err)
		}

		var tagOwner models.User
		err = FindUserById(&tagOwner, int(tag.UserId))
		if err != nil {
			return commonErrors.ErrorHandler(c, fiber.StatusNotFound, err)
		}

		serializedTag := serializers.TagSerializer(*tag, serializers.UserSerializer(&tagOwner))
		serializedTags = append(serializedTags, serializedTag)
	}

	var serializedAmenities []serializers.Amenity
	for _, amenity := range camping.Amenities {
		amenityItem, err := database.Database.FindByAmenityId(int(amenity.Id))
		if err != nil {
			return commonErrors.ErrorHandler(c, fiber.StatusNotFound, err)
		}

		amenityOwner, err := database.Database.FindByUserId(int(amenityItem.UserId))
		if err != nil {
			return commonErrors.ErrorHandler(c, fiber.StatusNotFound, err)
		}

		serializedAmenity := serializers.AmenitySerializer(*amenityItem, serializers.UserSerializer(amenityOwner))
		serializedAmenities = append(serializedAmenities, serializedAmenity)

	}

	serializedUser := serializers.UserSerializer(&user)
	serializedCamping := serializers.CampingSerializer(&camping, serializedUser, serializedTags, serializedAmenities)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "success",
		"data":    serializedCamping,
	})
}

func UpdateCamping(c *fiber.Ctx) error {

	owner, err := middleware.GetAuthUser(c)
	if err != nil {
		return commonErrors.ErrorHandler(c, fiber.StatusUnauthorized, err)
	}

	// urlparsing
	campingId, err := c.ParamsInt("campingId")
	if err != nil {
		return commonErrors.ErrorHandler(c, fiber.StatusBadRequest, err)
	}

	// FindDB
	var camping models.Camping
	if err := database.Database.Conn.First(&camping, "id = ?", campingId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return commonErrors.ErrorHandler(c, fiber.StatusNotFound, err)
		}
		return commonErrors.ErrorHandler(c, fiber.StatusBadRequest, err)
	}

	// user가 같은지 체크
	if err := checkUserEqualsRecordUser(owner, camping); err != nil {
		return commonErrors.ErrorHandler(c, fiber.StatusBadRequest, err)
	}

	type UpdatePayload struct {
		Title          string         `json:"title"`
		Address        string         `json:"address"`
		Description    string         `json:"description"`
		View           enums.ViewKind `json:"view"`
		IsEvCharge     enums.Status   `json:"is_ev_charge"`
		MannerTime     string         `json:"manner_time"`
		IsSideParking  enums.Status   `json:"is_side_parking"`
		IsPetFriendly  enums.Status   `json:"is_pet_friendly"`
		VisitedStartAt string         `json:"visited_start_at"`
		VisitedEndAt   string         `json:"visited_end_at"`
		Tags           []int          `json:"tags"`
		Amenities      []int          `json:"amenities"`
	}

	payload := &UpdatePayload{}

	if err := c.BodyParser(&payload); err != nil {
		return commonErrors.ErrorHandler(c, fiber.StatusBadRequest, err)
	}

	if len(payload.View) > 0 {
		err := setEnumView(payload.View)
		if err != nil {
			return commonErrors.ErrorHandler(c, fiber.StatusBadRequest, err)
		}

	}

	if len(payload.IsEvCharge) > 0 {
		err := setEnumStatus(payload.IsEvCharge)
		if err != nil {
			return commonErrors.ErrorHandler(c, fiber.StatusBadRequest, err)
		}

	}

	if len(payload.IsSideParking) > 0 {
		err := setEnumStatus(payload.IsSideParking)
		if err != nil {
			return commonErrors.ErrorHandler(c, fiber.StatusBadRequest, err)
		}

	}

	if len(payload.IsPetFriendly) > 0 {
		err := setEnumStatus(payload.IsPetFriendly)
		if err != nil {
			return commonErrors.ErrorHandler(c, fiber.StatusBadRequest, err)
		}

	}

	var serializedTags []serializers.Tag
	var tags []*models.Tag
	if len(payload.Tags) > 0 {
		for _, tagId := range payload.Tags {
			tag, err := FindByTagId(tagId)
			if err != nil {
				return commonErrors.ErrorHandler(c, fiber.StatusNotFound, err)
			}

			var tagOwner models.User
			err = FindUserById(&tagOwner, int(tag.UserId))
			if err != nil {
				return commonErrors.ErrorHandler(c, fiber.StatusNotFound, err)
			}

			serializedTag := serializers.TagSerializer(*tag, serializers.UserSerializer(&tagOwner))
			tags = append(tags, tag)
			serializedTags = append(serializedTags, serializedTag)
		}
	}

	m := structs.Map(&payload)
	fields := structs.Fields(&payload)
	for _, field := range fields {
		if field.IsZero() {
			delete(m, field.Name())
		}
	}

	m["Tags"] = tags

	database.Database.Conn.Model(&camping).Updates(m)
	if err := database.Database.Conn.Model(&camping).Association("Tags").Replace(tags); err != nil {
		return commonErrors.ErrorHandler(c, fiber.StatusBadRequest, err)
	}
	serializedUser := serializers.UserSerializer(owner)
	serializedCamping := serializers.CampingSerializer(&camping, serializedUser, serializedTags)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "success",
		"data":    serializedCamping,
	})

}

func DeleteCamping(c *fiber.Ctx) error {
	owner, err := middleware.GetAuthUser(c)
	if err != nil {
		return commonErrors.ErrorHandler(c, fiber.StatusUnauthorized, err)
	}

	campingId, err := c.ParamsInt("campingId")
	if err != nil {
		return commonErrors.ErrorHandler(c, fiber.StatusBadRequest, err)
	}

	// database
	var camping models.Camping
	if err := database.Database.Conn.First(&camping, "id = ?", campingId).Error; err != nil {
		return commonErrors.ErrorHandler(c, fiber.StatusNotFound, err)
	}

	if err := checkUserEqualsRecordUser(owner, camping); err != nil {
		return commonErrors.ErrorHandler(c, fiber.StatusForbidden, err)
	}

	database.Database.Conn.Delete(&camping)

	return c.SendStatus(fiber.StatusOK)
}

func checkUserEqualsRecordUser(owner *models.User, camping models.Camping) error {
	if owner.ID != camping.UserId {
		return errors.New("당신은 작성자가 아니군요!!")
	}

	return nil

}

func setEnumView(view enums.ViewKind) error {
	err := view.String()
	if err != nil {
		return errors.New(fmt.Sprintf("%s: %s", view, err.Error()))
	}
	return nil

}

func setEnumStatus(status ...enums.Status) error {
	for _, s := range status {
		err := s.String()
		if err != nil {
			return errors.New(fmt.Sprintf("%s: %s", s, err.Error()))
		}
	}
	return nil
}

func pprint(arg interface{}) string {
	str, err := json.MarshalIndent(arg, "", "\t")
	if err != nil {
		log.Println(err)
	}

	return string(str)
}
