package handlers

import (
	"camping-backend/database"
	"camping-backend/enums"
	"camping-backend/models"
	"camping-backend/serializers"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateCamping(c *fiber.Ctx) error {
	user, err := authUser(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}

	var request = models.Camping{
		CreatedAt: database.DB.NowFunc(),
		UpdatedAt: database.DB.NowFunc(),
		UserID:    user.ID,
	}

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

	database.DB.Create(&request)

	responseUser := serializers.UserSerializer(user)
	responseCamping := serializers.CampingSerializer(&request, responseUser)

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
		if err := FindUserById(&owner, int(camping.UserID)); err != nil {
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

	if err := FindUserById(&user, int(camping.UserID)); err != nil {
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

	owner, err := authUser(c)
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
	owner, err := authUser(c)
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
	if owner.ID != camping.UserID {
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
