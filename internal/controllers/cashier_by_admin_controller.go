package controllers

import (
	"my-go-api/internal/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CashierController struct {
	DB *gorm.DB
}

// ✅ CREATE Kasir Baru (dengan password & role cashier)
func (cc *CashierController) CreateCashier(c *fiber.Ctx) error {
	var body struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	cashier := models.User{
		Name:     body.Name,
		Username: body.Username,
		Password: string(hashedPassword),
		Role:     "cashier",
		Status:   bool(false),
	}

	if err := cc.DB.Create(&cashier).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Cashier created successfully",
		"data":    cashier,
	})
}

// ✅ update data kasir
func (cc *CashierController) UpdateCashier(c *fiber.Ctx) error {
	id := c.Params("id")

	// ✅ Ambil data kasir berdasarkan ID & pastikan role = cashier
	var cashier models.User
	if err := cc.DB.Where("role = ?", "cashier").First(&cashier, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Cashier not found",
		})
	}

	// ✅ Ambil input dari body (JSON)
	var body struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// ✅ Update data jika ada input
	if body.Name != "" {
		cashier.Name = body.Name
	}
	if body.Username != "" {
		cashier.Username = body.Username
	}
	if body.Password != "" {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		cashier.Password = string(hashed)
	}

	cashier.UpdatedAt = time.Now()

	// ✅ Simpan perubahan ke database
	if err := cc.DB.Save(&cashier).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Cashier updated successfully",
		"data":    cashier,
	})
}

// ✅ READ Semua Kasir
func (cc *CashierController) GetCashiers(c *fiber.Ctx) error {
	var cashiers []models.User
	if err := cc.DB.Where("role = ?", "cashier").Find(&cashiers).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(cashiers)
}

// ✅ UPDATE Status Kasir oleh Admin
func (cc *CashierController) UpdateCashierStatus(c *fiber.Ctx) error {
	id := c.Params("id")

	var body struct {
		Status bool `json:"status"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var cashier models.User
	if err := cc.DB.Where("role = ?", "cashier").First(&cashier, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Cashier not found"})
	}

	cashier.Status = body.Status

	if err := cc.DB.Save(&cashier).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Cashier status updated successfully",
		"data":    cashier,
	})
}

// ✅ DELETE Kasir
func (cc *CashierController) DeleteCashier(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := cc.DB.Where("role = ?", "cashier").Delete(&models.User{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Cashier deleted successfully"})
}
