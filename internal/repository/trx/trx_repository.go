package trx_repository

import (
	"test-rakamin/internal/models"

	"gorm.io/gorm"
)

type TrxRepository interface {
	Create(trx *models.Trx) error
	FindByUserID(userID uint) ([]models.Trx, error)
	FindByIDAndUserID(id uint, userID uint) (*models.Trx, error)
}

type trxRepositoryImpl struct {
	db *gorm.DB
}

func NewTrxRepository(db *gorm.DB) TrxRepository {
	return &trxRepositoryImpl{db: db}
}

func (r *trxRepositoryImpl) Create(trx *models.Trx) error {
	return r.db.Create(trx).Error
}

func (r *trxRepositoryImpl) FindByUserID(userID uint) ([]models.Trx, error) {
	var trxList []models.Trx
	err := r.db.Preload("DetailTrx").Where("id_user = ?", userID).Find(&trxList).Error
	return trxList, err
}

func (r *trxRepositoryImpl) FindByIDAndUserID(id uint, userID uint) (*models.Trx, error) {
	var trx models.Trx
	err := r.db.Preload("DetailTrx").Where("id = ? AND id_user = ?", id, userID).First(&trx).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &trx, err
}
