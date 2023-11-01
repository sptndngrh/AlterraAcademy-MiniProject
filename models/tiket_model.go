package models

import "time"

type Ticket struct {
	TiketID                uint             `gorm:"primaryKey;not null;autoIncrement" json:"tiket_id"`
	UserID                 uint             `json:"user_id"`
	OwnerID                uint             `json:"owner_id"`
	PropertiID             uint             `json:"properti_id"`
	LateTime               time.Time       `json:"tenggat_waktu"`
	Code_ticket            string           `json:"code_tiket"`
	DurasiBulan            int              `json:"durasi_bulan"`
	DurasiTahun            int              `json:"durasi_tahun"`
	Price                  int              `json:"harga"`
	Status                 TicketStatus     `json:"status" gorm:"type:enum('Belum Tersewa', 'Sudah Tersewa')"`
	Bulanan                bool             `json:"bulanan"`
	CreatedAt              time.Time       `json:"created_at"`
	PaymentMethod          string           `json:"metode_pembayaran"`
	PaymentDescription     string           `json:"deskripsi_pembayaran"`
	PaymentTotal           int              `json:"total_harga"`
	PaymentStatus          PembayaranStatus `json:"status_pembayaran" gorm:"type:enum('Sudah Lunas', 'Belum Lunas')"`
	LatePenaltyDescription string           `json:"deskripsi_denda"`
	LatePenaltyFee         int              `json:"denda_keterlambatan"`
	LatePenaltyStatus      PembayaranStatus `json:"status_denda" gorm:"type:enum('Sudah Lunas', 'Belum Lunas')"`
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

type PembayaranStatus string

const (
	// Status pembayaran enum
	StatusBelumLunas PembayaranStatus = "Belum Lunas"
	StatusSudahLunas PembayaranStatus = "Sudah Lunas"
)
