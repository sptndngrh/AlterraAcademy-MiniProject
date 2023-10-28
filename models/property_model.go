package models

import (
	"time"
)

type Property struct {
	PropertiID  uint      `gorm:"primaryKey;not null;autoIncrement" json:"properti_id"`
	OwnerID     uint      `json:"owner_id"`
	Title       string    `json:"judul"`
	Type        string    `json:"tipe"`
	Price       int       `gorm:"type:decimal(20,2)" json:"harga"`
	Description string    `json:"deskripsi"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"update_at"`
	DeletedAt   time.Time `json:"delete_at"`
	Ticket      Ticket    `json:"tickets" gorm:"foreignKey:PropertiID"`  
}

func (u *Property) TableName() string {
	return "properti"
}
