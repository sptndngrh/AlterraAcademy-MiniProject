package models

import (
	"time"
)

type Property struct {
	PropertiID uint       `gorm:"primaryKey;not null;autoIncrement" json:"properti_id"`
	OwnerID    uint       `json:"owner_id"`
	Judul      string     `json:"judul"`
	Tipe       string     `json:"tipe"`
	Harga      int        `gorm:"not null" json:"harga"`
	Lokasi     string     `json:"lokasi"`
	Deskripsi  string     `json:"deskripsi"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"update_at"`
	DeletedAt  *time.Time `json:"delete_at"`
	Ticket     Ticket     `json:"tickets" gorm:"foreignKey:PropertiID"`
}

func (u *Property) TableName() string {
	return "properti"
}

type PropertyStatus string
