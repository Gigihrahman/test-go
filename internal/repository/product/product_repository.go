package product_repository

import (
	"test-rakamin/internal/models"

	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *models.Product) error
	FindAllWithFilter(nama, categoryID, tokoID, minHarga, maxHarga string) ([]models.Product, error)
	FindByID(id uint) (*models.Product, error)
	Update(product *models.Product) error
	Delete(id uint) error
}

type productRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepositoryImpl{db: db}
}

func (r *productRepositoryImpl) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepositoryImpl) FindAllWithFilter(nama, categoryID, tokoID, minHarga, maxHarga string) ([]models.Product, error) {
	var products []models.Product
	query := r.db.Model(&models.Product{}).Preload("Category").Preload("Toko").Preload("FotoProduk")

	if nama != "" {
		query = query.Where("nama_produk LIKE ?", "%"+nama+"%")
	}
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if tokoID != "" {
		query = query.Where("toko_id = ?", tokoID)
	}
	if minHarga != "" {
		query = query.Where("harga_konsumen >= ?", minHarga)
	}
	if maxHarga != "" {
		query = query.Where("harga_konsumen <= ?", maxHarga)
	}

	err := query.Find(&products).Error
	return products, err
}

func (r *productRepositoryImpl) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Category").Preload("Toko").Preload("FotoProduk").First(&product, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &product, err
}

func (r *productRepositoryImpl) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}
