package models

import "gorm.io/gorm"

type ProductCategory string

const (
	Coffee    ProductCategory = "coffee"
	NonCoffee ProductCategory = "non_coffee"
	Snack     ProductCategory = "snack"
	HeavyMeal ProductCategory = "heavy_meal"
)

type Product struct {
	gorm.Model
	ID          uint            `json:"id" gorm:"primaryKey"`
	ProductName string          `json:"productName"`
	Description string          `json:"description"`
	Price       int             `json:"price"`
	Category    ProductCategory `json:"category"`
	ImgURL      string          `json:"img_url"`
	Available   bool            `json:"available" gorm:"default:false"`
	ImgPublicID string          `json:"-"`
}
