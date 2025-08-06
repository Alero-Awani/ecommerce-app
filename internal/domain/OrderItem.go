package domain

import "time"

type OrderItem struct {
	ID 				 uint    		`json:"id" gorm:"primaryKey"`
	OrderId 	 uint    		`json:"order_id"`
	ProductId  uint    		`json:"product_id"`
	Name 			 string  		`json:"name"`
	SellerId 	 uint    		`json:"seller_id"`
	ImageUrl 	 string  		`json:"image_url"`
	Qty 			 uint     	`json:"qty"`
	Price 		 float64 		`json:"price"`
	CreatedAt  time.Time  `gorm:"default:current_timestamp"`
	UpdatedAt  time.Time  `gorm:"default:current_timestamp"`	
}
