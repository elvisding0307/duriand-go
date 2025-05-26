package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID           uint   `json:"id" gorm:"primaryKey"`
	Username     string `json:"username" gorm: "unique"`
	Password     string `json:"password"`
	CorePassword string `json:"core_password"`
}
