package models

import (
	"time"
)

type Property struct {
	Id            uint       `gorm:"primaryKey;not null;autoIncrement" json:"id"`
	OwnerId        uint       `gorm:"foreignKey:UserId" json:"owner_id"`
	Judul         string     `json:"judul"`
	Tipe          string     `json:"tipe"`
	Harga         int        `json:"harga"` // Harga untuk property untuk setiap bulannya
	Lokasi        string     `json:"lokasi"`
	StatusTersewa bool       `gorm:"default:false" json:"status_tersewa"` // default false jka tersedia, jika sudah ada yang menyewa status akan menjadi true dan tidak bisa di sewa lagi, sampai user mengubah statusnya menjadi false lagi/ saat tanggal sewa pembeli sebelumnya sudah selesai
	Deskripsi     string     `json:"deskripsi"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"update_at"`
	DeletedAt     *time.Time `json:"delete_at"`
}

func (u *Property) TableName() string {
	return "properti"
}

type PropertyStatus string
