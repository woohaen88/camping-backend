package models

import (
	"camping-backend/common"
	"camping-backend/enums"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uint       `json:"id" gorm:"primaryKey"`
	Email    string     `json:"email"`
	Password string     `json:"password"`
	Username string     `json:"username"`
	Campings []Camping  `gorm:"foreignKey:UserID"`
	Role     enums.Role `json:"role" gorm:"default:'client'"`
}

func (u *User) PaswordHash(password string) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	common.CheckErr(err)
	u.Password = string(hashPassword)
}
