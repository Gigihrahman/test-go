package toko_repository

import (
	"test-rakamin/internal/models"

	"gorm.io/gorm"
)

type TokoRepository interface {
	Create(toko *models.Toko) error
	FindAll() ([]models.Toko, error)
	FindByID(id uint) (*models.Toko, error)
	FindByUserID(userID uint) (*models.Toko, error)
	Update(toko *models.Toko) error
	Delete(id uint) error
}

type tokoRepositoryImpl struct {
	db *gorm.DB
}

func NewTokoRepository(db *gorm.DB) TokoRepository {
	return &tokoRepositoryImpl{db: db}
}

func (r *tokoRepositoryImpl) Create(toko *models.Toko) error {
	return r.db.Create(toko).Error
}

func (r *tokoRepositoryImpl) FindAll() ([]models.Toko, error) {
	var tokoList []models.Toko
	err := r.db.Find(&tokoList).Error
	return tokoList, err
}

func (r *tokoRepositoryImpl) FindByID(id uint) (*models.Toko, error) {
	var toko models.Toko
	err := r.db.First(&toko, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &toko, err
}

func (r *tokoRepositoryImpl) FindByUserID(userID uint) (*models.Toko, error) {
	var toko models.Toko
	err := r.db.Where("id_user = ?", userID).First(&toko).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &toko, err
}

func (r *tokoRepositoryImpl) Update(toko *models.Toko) error {
	return r.db.Save(toko).Error
}

func (r *tokoRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Toko{}, id).Error
}
