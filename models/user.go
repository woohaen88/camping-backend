package models

import (
	"camping-backend/common"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID uint `json:"id" gorm:"primaryKey"`
	Email string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}


func (u *User) PaswordHash(){
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	common.CheckErr(err)
	u.Password = string(hashPassword)
}