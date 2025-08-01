package controllers

import (
	"my-go-api/internal/config"
	"my-go-api/internal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateOrder(c *fiber.Ctx) error {
	var req struct {
		TableID uint `json:"table_id"`
		Items   []struct {
			ProductID uint `json:"product_id"`
			Quantity  int  `json:"quantity"`
		} `json:"items"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	var order models.Order

	// ✅ Cek jika sudah ada order pending di meja ini
	err := config.DB.Where("table_id = ? AND status = ?", req.TableID, models.Pending).
		Preload("Items.Product").
		Preload("Table").
		First(&order).Error

	if err == gorm.ErrRecordNotFound {
		// ✅ Jika tidak ada, buat order baru
		order = models.Order{
			TableID: req.TableID,
			Status:  models.Pending,
		}
		if err := config.DB.Create(&order).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create order"})
		}
	} else if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	// ✅ Tambahkan item ke order
	for _, item := range req.Items {
		if item.Quantity > 0 {
			newItem := models.OrderItem{
				OrderID:   order.ID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
			}
			if err := config.DB.Create(&newItem).Error; err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Failed to add item"})
			}
		}
	}

	// ✅ Tandai meja sebagai occupied
	config.DB.Model(&models.Table{}).Where("id = ?", req.TableID).Update("is_occupied", true)

	// ✅ Ambil ulang order lengkap dengan semua relasi
	if err := config.DB.
		Preload("Table").
		Preload("Items.Product").
		First(&order, order.ID).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch final order"})
	}

	return c.Status(201).JSON(order)
}

func GetOrdersByTable(c *fiber.Ctx) error {
	tableID := c.Params("table_id")
	var orders []models.Order

	config.DB.Preload("Table").Preload("Items").Preload("Items.Product").Where("table_id = ?", tableID).Find(&orders)
	return c.JSON(orders)
}
