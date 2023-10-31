package serializers

import "camping-backend/models"

type User struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func UserSerializer(userModel *models.User) User {
	return User{
		ID:       userModel.ID,
		Email:    userModel.Email,
		Username: userModel.Username,
	}
}
