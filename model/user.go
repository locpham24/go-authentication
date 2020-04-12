package model

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserForm struct {
	Username string `json:"username" validate:"required" binding:"required,min=4,max=255"`
	Password string `json:"password" validate:"required" binding:"required,min=8,max=255"`
}

func (u *User) Fill(f UserForm) {
	u.Username = f.Username
	u.Password = f.Password
}
