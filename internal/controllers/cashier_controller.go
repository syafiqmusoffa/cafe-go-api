package controllers

import (
	"my-go-api/internal/config"
	"my-go-api/internal/models"

	"github.com/gofiber/fiber/v2"
)

func UpdateOrderStatus(c *fiber.Ctx) error {
	// ✅ Ambil user_id & role dari JWT (dengan pengecekan aman)
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized or invalid user",
		})
	}
	role := c.Locals("role").(string)

	// ✅ Pastikan hanya kasir
	if role != "cashier" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Only cashiers can update orders",
		})
	}

	// ✅ Cek apakah kasir aktif
	var cashier models.User
	if err := config.DB.First(&cashier, userID).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Cashier not found",
		})
	}

	if cashier.Status != true {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Only active cashiers can update orders",
		})
	}

	// ✅ Ambil order berdasarkan ID
	id := c.Params("id")
	var order models.Order
	if err := config.DB.Preload("Table").First(&order, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Order not found"})
	}

	// ✅ Parse status baru
	var req struct {
		Status models.OrderStatus `json:"status"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// ✅ Update status order
	order.Status = req.Status
	config.DB.Save(&order)

	// ✅ Update meja
	switch req.Status {
	case models.Canceled, models.Paid:
		config.DB.Model(&order.Table).Update("is_occupied", false)
	default:
		config.DB.Model(&order.Table).Update("is_occupied", true)
	}

	return c.JSON(order)
}

// controllers/order.go
func GetActiveOrders(c *fiber.Ctx) error {
	var orders []models.Order

	// Ambil semua order aktif dengan semua relasi
	if err := config.DB.
		Preload("Table").
		Preload("Items.Product").
		Where("status NOT IN ?", []string{"paid", "cancelled"}).
		Find(&orders).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal mengambil data order aktif",
		})
	}

	return c.JSON(orders)
}

func GetAllOrders(c *fiber.Ctx) error {
	var orders []models.Order

	if err := config.DB.
		Preload("Items.Product").
		Preload("Table").
		Order("created_at DESC").Limit(10).
		Find(&orders).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch orders: " + err.Error(),
		})
	}

	return c.JSON(orders)
}
