package models

import "gorm.io/gorm"

type OrderItem struct {
	gorm.Model
	OrderID uint  `json:"order_id"`
	Order   Order `json:"-"`

	ProductID uint    `json:"product_id"`
	Product   Product `json:"product"`

	Quantity int `json:"quantity"`
}
