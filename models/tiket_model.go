package models

import "time"

type Ticket struct {
	TiketID         uint          `gorm:"primaryKey;not null;autoIncrement" json:"tiket_id"`
	UserID          uint          `json:"user_id"`
	PropertiID      uint          `json:"properti_id"`
	LateTime        time.Time     `json:"tenggat_waktu"`
	Code_ticket     string        `json:"code_tiket"`
	PropertyAddress string        `json:"alamat_properti"`
	Description     string        `json:"deskripsi"`
	Status          TicketStatus  `json:"status" gorm:"type:enum('Belum Tersewa', 'Sudah Tersewa')"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"update_at"`
	Payments        Payment       `json:"payments" gorm:"foreignKey:TiketID"`
	LatePenalties   []LatePenalty `json:"LatePenalties" gorm:"foreignKey:TicketID"`
}

func (u *Ticket) TableName() string {
	return "tiket"
}

type TicketStatus string

const (
	// Tiket status enum
	StatusBelumTersewa TicketStatus = "Belum Tersewa"
	StatusSudahTersewa TicketStatus = "Sudah Tersewa"
)
