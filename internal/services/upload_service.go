package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"my-go-api/internal/config"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
)

type UploadedImage struct {
	URL      string `json:"url"`
	PublicID string `json:"public_id"`
}

func UploadImage(file multipart.File, fileHeader *multipart.FileHeader) (UploadedImage, error) {
	ctx := context.Background()

	// ✅ Buat PublicID unik agar tidak tertimpa
	publicID := fmt.Sprintf("products/%d_%s_%s",
		time.Now().Unix(),
		uuid.New().String(),
		fileHeader.Filename,
	)

	uploadResp, err := config.Cloud.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID: publicID,
		Folder:   "products",
	})
	if err != nil {
		return UploadedImage{}, fmt.Errorf("cloudinary upload failed: %w", err)
	}

	return UploadedImage{
		URL:      uploadResp.SecureURL,
		PublicID: uploadResp.PublicID,
	}, nil
}

func DeleteImage(publicID string) error {
	if publicID == "" {
		return nil
	}

	ctx := context.Background()
	_, err := config.Cloud.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return fmt.Errorf("cloudinary delete failed: %w", err)
	}
	return nil
}
