package handlers

import (
	commonErrors "camping-backend/common/errors"
	"camping-backend/database"
	"camping-backend/enums"
	"camping-backend/middleware"
	"camping-backend/models"
	"camping-backend/serializers"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateCamping(c *fiber.Ctx) error {

	user, err := middleware.GetAuthUser(c)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}

	var camping = models.Camping{
		CreatedAt: database.DB.NowFunc(),
		UpdatedAt: database.DB.NowFunc(),
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": fmt.Sprintf("%s", err.Error()),
		})
	}

	if err := setEnumStatus(request.IsEvCharge, request.IsSideParking, request.IsPetFriendly); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": fmt.Sprintf("%s", err.Error()),
		})
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
		err = database.DB.First(&tagModel, "id = ?", tagId).Error
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "could not pasing tag",
				"data":    err.Error(),
			})
		}

		tagModels = append(tagModels, tagModel)
	}

	var serializedTags []serializers.Tag
	for i := 0; i < len(tagModels); i++ {
		tagModel := tagModels[i]
		var tagUser models.User
		err = FindUserById(&tagUser, int(tagModel.UserId))
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "could not User",
				"data":    err.Error(),
			})
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
		if err := database.DB.First(&amenity, "id = ?", amenityId).Error; err != nil {
			return commonErrors.ErrorHandler(c, fiber.StatusNotFound, err)
		}

		if err := FindUserById(&amenityCreator, int(amenity.UserId)); err != nil {
			return commonErrors.ErrorHandler(c, fiber.StatusNotFound, err)
		}

		serializedAmenity := serializers.AmenitySerializer(amenity, serializers.UserSerializer(&amenityCreator))
		serializedAmenities = append(serializedAmenities, serializedAmenity)
		amenities = append(amenities, amenity)
	}

	database.DB.Create(&camping)

	if err := database.DB.Model(&camping).Association("Tags").Append(tagModels); err != nil {
		return commonErrors.ErrorHandler(c, fiber.StatusBadRequest, err, "failed many to many field camping=>tags")
	}

	if err := database.DB.Model(&camping).Association("Amenities").Append(amenities); err != nil {
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

	database.DB.Find(&campings)

	var responseCampings []serializers.Camping

	for _, camping := range campings {
		if err := FindUserById(&owner, int(camping.UserId)); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Couldn't change password",
				"data":    err.Error(),
			})
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "error",
			"data":    err.Error(),
		})
	}

	var camping models.Camping
	var user models.User

	if err := database.DB.First(&camping, "id = ?", campingId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "error",
				"data":    err.Error(),
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "error",
			"data":    err.Error(),
		})
	}

	if err := FindUserById(&user, int(camping.UserId)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "error",
			"data":    err.Error(),
		})
	}

	serializedUser := serializers.UserSerializer(&user)
	serializedCamping := serializers.CampingSerializer(&camping, serializedUser)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "success",
		"data":    serializedCamping,
	})
}

func UpdateCamping(c *fiber.Ctx) error {

	owner, err := middleware.GetAuthUser(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "Unauthorized",
			"message": "error",
			"data":    err.Error(),
		})
	}

	// urlparsing
	campingId, err := c.ParamsInt("campingId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "error",
			"data":    err.Error(),
		})
	}

	// FindDB
	var camping models.Camping
	if err := database.DB.First(&camping, "id = ?", campingId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "error",
				"data":    err.Error(),
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "error",
			"data":    err.Error(),
		})
	}

	// user가 같은지 체크
	if err := checkUserEqualsRecordUser(owner, camping); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "error",
			"data":    err.Error(),
		})
	}

	//payload
	//업데이트 시간 변경
	type Payload struct {
		Title          string
		Address        string
		Description    string
		View           enums.ViewKind
		IsEvCharge     enums.Status
		MannerTime     string
		IsSideParking  enums.Status
		IsPetFriendly  enums.Status
		VisitedStartAt string
		VisitedEndAt   string
	}

	payload := &Payload{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "error",
			"data":    err.Error(),
		})
	}

	if len(payload.View) > 0 {
		err := setEnumView(payload.View)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "error",
				"data":    err.Error(),
			})
		}
	}

	if len(payload.IsEvCharge) > 0 {
		err := setEnumStatus(payload.IsEvCharge)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "error",
				"data":    err.Error(),
			})
		}
	}

	if len(payload.IsSideParking) > 0 {
		err := setEnumStatus(payload.IsSideParking)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "error",
				"data":    err.Error(),
			})
		}
	}

	if len(payload.IsPetFriendly) > 0 {
		err := setEnumStatus(payload.IsPetFriendly)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "error",
				"data":    err.Error(),
			})
		}
	}

	database.DB.Model(&camping).Updates(payload)

	serializedUser := serializers.UserSerializer(owner)
	serializedCamping := serializers.CampingSerializer(&camping, serializedUser)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "success",
		"data":    serializedCamping,
	})
}

func DeleteCamping(c *fiber.Ctx) error {
	owner, err := middleware.GetAuthUser(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "success",
			"message": "success",
			"data":    err.Error(),
		})
	}

	campingId, err := c.ParamsInt("campingId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "error",
			"data":    err.Error(),
		})
	}

	// database
	var camping models.Camping
	if err := database.DB.First(&camping, "id = ?", campingId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "error",
			"data":    err.Error(),
		})
	}

	if err := checkUserEqualsRecordUser(owner, camping); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "error",
			"data":    err.Error(),
		})
	}

	database.DB.Delete(&camping)

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
