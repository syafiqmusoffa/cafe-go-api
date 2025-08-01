package models

import "gorm.io/gorm"

type Table struct {
	gorm.Model
	ID          uint `json:"id" gorm:"primaryKey"`
	TableNumber int  `json:"table_number" `
	IsOccupied  bool `json:"is_occupied"`
}
