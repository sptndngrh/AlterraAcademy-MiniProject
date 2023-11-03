package models

import "time"

type Order struct {
	Id            uint       `gorm:"primaryKey;not null;autoIncrement" json:"id"`
	UserId        uint       `gorm:"foreignKey:UserId" json:"user_id"`
	PropertiId    uint       `gorm:"foreignKey:PropertiId" json:"properti_id"`
	Bulanan       int        `json:"bulanan"`
	PaymentTotal  int        `json:"total_harga"`
	PaymentStatus bool       `gorm:"default:false" json:"payment_status"`
	StartDate     time.Time  `json:"start_date"`
	EndDate       time.Time  `json:"end_date"`
	CreatedAt     *time.Time `json:"created_at"`
	DeletedAt     *time.Time `json:"delete_at"`
}

func (u *Order) TableName() string {
	return "order"
}

