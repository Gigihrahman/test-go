package product_photo_repository

import (
	"test-rakamin/internal/models"

	"gorm.io/gorm"
)

type ProductPhotoRepository interface {
	Create(photo *models.ProductPhoto) error
	DeleteByProductID(productID uint) error
}

type productPhotoRepositoryImpl struct {
	db *gorm.DB
}

func NewProductPhotoRepository(db *gorm.DB) ProductPhotoRepository {
	return &productPhotoRepositoryImpl{db: db}
}

func (r *productPhotoRepositoryImpl) Create(photo *models.ProductPhoto) error {
	return r.db.Create(photo).Error
}

func (r *productPhotoRepositoryImpl) DeleteByProductID(productID uint) error {
	return r.db.Where("product_id = ?", productID).Delete(&models.ProductPhoto{}).Error
}
