package category_repository

import (
	"test-rakamin/internal/models"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(category *models.Category) error
	FindAll() ([]models.Category, error)
	FindByID(id uint) (*models.Category, error)
	Update(category *models.Category) error
	Delete(id uint) error
}

type categoryRepositoryImpl struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepositoryImpl{db: db}
}

func (r *categoryRepositoryImpl) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepositoryImpl) FindAll() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Find(&categories).Error
	return categories, err
}

func (r *categoryRepositoryImpl) FindByID(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &category, err
}

func (r *categoryRepositoryImpl) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}
