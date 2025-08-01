package config

import (
	"log"
	"my-go-api/internal/models"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	var err error
	dsn := os.Getenv("DB_URL")

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Gagal konek database:", err)
	}

	// Buat ENUM order_status kalau belum ada
	DB.Exec(`
	DO $$ 
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'order_status') THEN
			CREATE TYPE order_status AS ENUM (
				'pending', 'validated', 'processing', 'completed', 'paid', 'canceled'
			);
		END IF;
	END$$;
	`)

	// ✅ Lakukan migrasi dengan urutan yang benar agar relasi tidak gagal
	if err := DB.AutoMigrate(
		&models.User{},
		&models.Table{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
	); err != nil {
		log.Fatal("❌ AutoMigrate error:", err)
	}

	log.Println("✅ Database connected & migrated")

	SeedUsers()
	return DB
}

func SeedUsers() {
	var count int64
	DB.Model(&models.User{}).Count(&count)
	if count == 0 {
		log.Println("✅ Seeding default users...")
		adminPass, _ := bcrypt.GenerateFromPassword([]byte("admin123"), 10)
		cashierPass, _ := bcrypt.GenerateFromPassword([]byte("kasir123"), 10)

		DB.Create(&models.User{
			Username: "admin",
			Password: string(adminPass),
			Role:     models.RoleAdmin,
		})
		DB.Create(&models.User{
			Name:     "kasir 1",
			Username: "kasir",
			Password: string(cashierPass),
			Role:     models.RoleCashier,
		})
		log.Println("✅ Admin & Kasir created")
	}
}
