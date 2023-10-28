package models

import "time"

type LatePenalty struct {
	LatePenaltyID uint   `gorm:"primaryKey;not null;autoIncrement" json:"denda_id"`
	TicketID      uint `json:"tiket_id"`
	Description   string `json:"deskripsi"`
	PenaltyFee    int    `gorm:"type:decimal(20,2)" json:"denda"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (u *LatePenalty) TableName() string {
	return "denda_keterlambatan"
}