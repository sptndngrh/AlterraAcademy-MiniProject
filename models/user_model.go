package models

import (
	"time"
)

type User struct {
	Id             uint      `gorm:"primaryKey;not null;autoIncrement" json:"id"`
	Nama           string    `json:"nama"`
	NoTelp         string    `json:"no_telp"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	Password       string    `json:"password"`
	OwnerRole      bool      `gorm:"default:false" json:"owner_role"`
	DoneVerify     bool      `gorm:"default:false" json:"done_verify"`
	JWTTokenVerify string    `json:"jwt_token_verify"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"update_at"`
	DeletedAt      time.Time `json:"delete_at"`
	Orders         []Order   `gorm:"foreignKey:UserId"`
}

func (u *User) TableName() string {
	return "user"
}

type ChangeUsernameRequest struct {
	CurrentUsername string `json:"currentUsername" binding:"required"`
	NewUsername     string `json:"newUsername" binding:"required"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required"`
}