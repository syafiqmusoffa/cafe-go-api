package controllers

import (
	"my-go-api/internal/config"
	"my-go-api/internal/models"
	"my-go-api/internal/services"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func CreateProduct(c *fiber.Ctx) error {
	product := models.Product{
		ProductName: c.FormValue("productName"),
		Description: c.FormValue("description"),
		Category:    models.ProductCategory(c.FormValue("category")),
	}

	if priceStr := c.FormValue("price"); priceStr != "" {
		if price, err := strconv.Atoi(priceStr); err == nil {
			product.Price = price
		}
	}

	// ✅ Upload gambar jika ada
	file, err := c.FormFile("image")
	if err == nil {
		// ✅ Validasi ukuran & tipe file
		if file.Size > 2*1024*1024 {
			return c.Status(400).JSON(fiber.Map{"error": "File too large (max 2MB)"})
		}
		if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid file type"})
		}

		// ✅ Buka file & upload
		f, _ := file.Open()
		uploaded, uploadErr := services.UploadImage(f, file)
		f.Close() // ✅ Tutup setelah upload selesai

		if uploadErr != nil {
			return c.Status(500).JSON(fiber.Map{
				"error":   "Upload failed",
				"details": uploadErr.Error(),
			})
		}

		// ✅ Simpan URL & PublicID ke database
		product.ImgURL = uploaded.URL
		product.ImgPublicID = uploaded.PublicID
	}

	if err := config.DB.Create(&product).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save product"})
	}

	return c.Status(201).JSON(product)
}

func UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var product models.Product

	if err := config.DB.First(&product, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
	}

	if name := c.FormValue("productName"); name != "" {
		product.ProductName = name
	}
	if desc := c.FormValue("description"); desc != "" {
		product.Description = desc
	}
	if cat := c.FormValue("category"); cat != "" {
		product.Category = models.ProductCategory(cat)
	}
	if priceStr := c.FormValue("price"); priceStr != "" {
		if price, err := strconv.Atoi(priceStr); err == nil {
			product.Price = price
		}
	}

	// ✅ Jika ada file gambar baru
	file, err := c.FormFile("image")
	if err == nil {
		if file.Size > 2*1024*1024 {
			return c.Status(400).JSON(fiber.Map{"error": "File too large (max 2MB)"})
		}
		if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid file type"})
		}
		f, _ := file.Open()
		// ❌ Jangan dulu defer f.Close() sebelum upload berhasil

		uploaded, uploadErr := services.UploadImage(f, file)
		f.Close() // ✅ Tutup setelah dipakai
		if uploadErr != nil {
			return c.Status(500).JSON(fiber.Map{"error": uploadErr.Error()})
		}

		product.ImgURL = uploaded.URL
		product.ImgPublicID = uploaded.PublicID
	}

	config.DB.Save(&product)
	return c.JSON(product)
}

func UpdateProductStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	var body struct {
		Available bool `json:"available" gorm:"default:false"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	var product models.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not Found"})
	}
	product.Available = body.Available
	if err := config.DB.Save(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal memperbarui status produk",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Status produk berhasil diperbarui",
		"product": product,
	})
}

func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var product models.Product

	if err := config.DB.First(&product, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
	}

	if product.ImgPublicID != "" {
		services.DeleteImage(product.ImgPublicID)
	}

	if err := config.DB.Delete(&product).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete product"})
	}

	return c.JSON(fiber.Map{"message": "Product deleted successfully"})
}

func GetProducts(c *fiber.Ctx) error {
	var products []models.Product
	config.DB.Find(&products)
	return c.JSON(products)
}

func CreateTable(c *fiber.Ctx) error {
	var table models.Table

	if err := c.BodyParser(&table); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// ✅ Cek apakah table_number sudah ada & belum dihapus
	var existing models.Table
	if err := config.DB.
		Where("table_number = ?", table.TableNumber).
		First(&existing).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Table number already exists",
		})
	}

	// ✅ Buat meja baru
	if err := config.DB.Create(&table).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(table)
}

func DeleteTable(c *fiber.Ctx) error {
	id := c.Params("id") // ambil parameter id dari URL

	var table models.Table

	// Cek apakah table dengan id tersebut ada
	if err := config.DB.First(&table, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Table not found"})
	}

	// Hapus table
	if err := config.DB.Delete(&table).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Table deleted successfully",
	})
}

func GetTables(c *fiber.Ctx) error {
	var tables []models.Table
	config.DB.Find(&tables)
	return c.JSON(tables)
}

func GetTableByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var table models.Table
	if err := config.DB.First(&table, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Table not found",
		})
	}

	return c.JSON(table)
}
