package models

import "gorm.io/gorm"

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleCashier UserRole = "cashier"
)

type User struct {
	gorm.Model
	Name     string   `json:"name"`
	Username string   `json:"username" gorm:"unique"`
	Password string   `json:"-"`
	Role     UserRole `json:"role" gorm:"type:varchar(20)"`
	Status   bool     `json:"status" gorm:"default:false"`
}
