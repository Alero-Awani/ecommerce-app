package domain

import "time"

type Address struct {
	ID           uint      `gorm:"PrimaryKey" json:"id"`
	AddressLine1 string    `json:"address_line_1"`
	AddressLine2 string    `json:"address_line_2"`
	City         string    `json:"city"`
	Postcode     string    `json:"postcode"`
	Country      string    `json:"country"`
	UserId       uint      `json:"user_id"`
	CreatedAt    time.Time `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:current_timestamp" json:"updated_at"`
}
