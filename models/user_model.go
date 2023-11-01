package models

import (
	"time"
)

type User struct {
	UserID          uint      `gorm:"primaryKey;not null;autoIncrement" json:"user_id"`
	Nama            string    `json:"nama"`
	NoTelp          string    `json:"no_telp"`
	Username        string    `json:"username"`
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	DoneVerify      bool      `gorm:"default:false" json:"done_verify"`
	UserTokenVerify string    `json:"user_token_verify"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"update_at"`
	DeletedAt       time.Time `json:"delete_at"`
	Tickets         []Ticket  `json:"tickets" gorm:"foreignKey:UserID"`
}

func (u *User) TableName() string {
	return "user"
}

type ChangeNameUserRequest struct {
	CurrentNama string `json:"currentNama" binding:"required"`
	NewNama     string `json:"newNama" binding:"required"`
}

type ChangePasswordUserRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required"`
}

type ChangeUsernameUserRequest struct {
	CurrentUsername string `json:"currentUsername" binding:"required"`
	NewUsername     string `json:"newUsername" binding:"required"`
}