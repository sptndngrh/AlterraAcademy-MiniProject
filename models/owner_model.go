package models

import (
	"time"
)

type Owner struct {
	OwnerID          uint       `gorm:"primaryKey;not null;autoIncrement" json:"owner_id"`
	Nama             string     `json:"nama"`
	NoTelp           string     `json:"no_telp"`
	Username         string     `json:"username"`
	Email            string     `json:"email"`
	Password         string     `json:"password"`
	DoneVerify       bool       `gorm:"default:false" json:"done_verify"`
	OwnerTokenVerify string     `json:"owner_token_verify"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"update_at"`
	DeletedAt        time.Time  `json:"delete_at"`
	Properties       []Property `json:"properties" gorm:"foreignKey:OwnerID"`
	Tickets          []Ticket   `json:"tickets" gorm:"foreignKey:OwnerID"`
}

func (u *Owner) TableName() string {
	return "owner"
}

type ChangeNameOwnerRequest struct {
	CurrentNama string `json:"currentNama" binding:"required"`
	NewNama     string `json:"newNama" binding:"required"`
}

type ChangePasswordOwnerRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required"`
}

type ChangeUsernameOwnerRequest struct {
	CurrentUsername string `json:"currentUsername" binding:"required"`
	NewUsername     string `json:"newUsername" binding:"required"`
}
