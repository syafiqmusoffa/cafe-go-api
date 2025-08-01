package config

import (
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/joho/godotenv"
)

var Cloud *cloudinary.Cloudinary

func InitCloudinary() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ .env not found, using system env")
	}

	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		log.Fatal("❌ Failed to init Cloudinary:", err)
	}
	Cloud = cld

	log.Println("✅ Cloudinary initialized")
}
