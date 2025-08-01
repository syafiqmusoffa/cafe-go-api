package models

import "gorm.io/gorm"

type OrderStatus string

const (
	Pending    OrderStatus = "pending"
	Validated  OrderStatus = "validated"
	Processing OrderStatus = "processing"
	Completed  OrderStatus = "completed"
	Paid       OrderStatus = "paid"
	Canceled   OrderStatus = "canceled"
)

type Order struct {
	gorm.Model
	ID      uint        `json:"id" gorm:"primaryKey"`
	TableID uint        `json:"table_id"`
	Table   Table       `json:"Table" gorm:"foreignKey:TableID"`
	Status  OrderStatus `json:"status" gorm:"type:varchar(20);default:'pending'"`
	Items   []OrderItem `json:"Items"`
}
