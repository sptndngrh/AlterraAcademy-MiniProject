package models

import "time"

type Payment struct {
	PaymentID     uint             `gorm:"primaryKey;not null;autoIncrement" json:"pembayaran_id"`
	TicketID      uint             `json:"tiket_id"`
	OrderAmount   int              `json:"jumlah_pemesanan"`
	PriceTotal    int              `gorm:"type:decimal(20,2)" json:"total_harga"`
	PaymentStatus PembayaranStatus `json:"status_pembayaran" gorm:"type:enum('Sudah Lunas', 'Belum Lunas')"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"update_at"`
}

func (u *Payment) TableName() string {
	return "pembayaran"
}

type PembayaranStatus string

const (
	// Tiket status enum
	StatusBelumLunas PembayaranStatus = "Belum Lunas"
	StatusSudahLunas PembayaranStatus = "Sudah Lunas"
)
