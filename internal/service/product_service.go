package service

import (
	"context"
	"fmt"
	"log/slog"
	"mime/multipart"
	"path/filepath"
	"sync"
	"venturo-core/internal/adapter/storage"
	"venturo-core/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProductService handles the business logic for products.
type ProductService struct {
	db             *gorm.DB
	wg             *sync.WaitGroup
	storageAdapter storage.StorageAdapter
}

// NewProductService creates a new product service.
func NewProductService(db *gorm.DB, wg *sync.WaitGroup, storageAdapter storage.StorageAdapter) *ProductService {
	return &ProductService{db: db, wg: wg, storageAdapter: storageAdapter}
}

// CreateProductInput is the data needed to create a new product.
type CreateProductInput struct {
	Name  string
	Price int32
	Stock int16
	Image *multipart.FileHeader
}

// CreateProduct creates a product and asynchronously uploads its image.
func (s *ProductService) CreateProduct(ctx context.Context, input CreateProductInput) (*model.Product, error) {
	product := model.Product{
		Name:  input.Name,
		Price: input.Price,
		Stock: input.Stock,
	}

	// If an image is provided, prepare for upload.
	if input.Image != nil {
		imageName := fmt.Sprintf("%s%s", uuid.NewString(), filepath.Ext(input.Image.Filename))
		product.ImageURL = imageName
		product.ImageStatus = "uploading"
	}

	// Save the initial product record. This is fast and synchronous.
	if err := product.Save(s.db); err != nil {
		return nil, err
	}

	// If there's an image, start the background upload process.
	if input.Image != nil {
		s.wg.Add(1)
		go s.uploadProductImage(product.ID, input.Image, product.ImageURL)
	}

	return &product, nil
}

// uploadProductImage is the background worker.
func (s *ProductService) uploadProductImage(productID uuid.UUID, file *multipart.FileHeader, objectName string) {
	defer s.wg.Done()
	bgCtx := context.Background()

	// 1. Upload the file.
	publicURL, err := s.storageAdapter.Upload(bgCtx, file, objectName)
	if err != nil {
		slog.Error("Failed to upload product image", "productID", productID, "error", err)
		// Update status to 'failed'
		s.updateImageStatus(productID, "failed", "")
		return
	}

	// 2. Update status to 'done' on success.
	slog.Info("Successfully uploaded product image", "productID", productID, "url", publicURL)
	s.updateImageStatus(productID, "done", publicURL)
}

// updateImageStatus is a helper to update the product record.
func (s *ProductService) updateImageStatus(productID uuid.UUID, status string, url string) {
	var product model.Product
	// Using a map to update specific fields.
	updates := map[string]interface{}{"image_status": status}
	if url != "" {
		updates["image_url"] = url
	}

	if err := s.db.Model(&product).Where("id = ?", productID).Updates(updates).Error; err != nil {
		slog.Error("Failed to update product image status", "productID", productID, "error", err)
	}
}
