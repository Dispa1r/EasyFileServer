package model

import (
	"time"
)

type TblUser struct {
	ID  uint `gorm:"column:id"`
	UserName string `gorm:"column:user_name"`
	UserPwd string`gorm:"column:user_pwd"`
	Email string `gorm:"column:email"`
	Phone string  `gorm:"column:phone"`
	EmailValidated int `gorm:"column:email_validated"`
	PhoneValidated int `gorm:"column:phone_validated"`
	LastActive time.Time `gorm:"column:last_active;default:null"`
	Profile string `gorm:"column:profile"`
	Status int `gorm:"column:status"`
	SignupAt time.Time `gorm:"column:signup_at;default:null"`
}