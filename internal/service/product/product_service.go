package product_service

import (
	"errors"
	"mime/multipart"
	"test-rakamin/internal/models"
	product_repository "test-rakamin/internal/repository/product"
	product_photo_repository "test-rakamin/internal/repository/product_photo"
	"test-rakamin/utils"
)

type ProductService interface {
	GetAllProducts(nama, categoryID, tokoID, minHarga, maxHarga string) ([]models.Product, error)
	GetProductByID(id uint) (*models.Product, error)
	CreateProduct(product *models.Product, photos []*multipart.FileHeader) (*models.Product, error)
	UpdateProduct(id uint, updatedProduct *models.Product, photos []*multipart.FileHeader) (*models.Product, error)
	DeleteProduct(id uint) error
}

type productServiceImpl struct {
	productRepo      product_repository.ProductRepository
	productPhotoRepo product_photo_repository.ProductPhotoRepository
}

func NewProductService(repo product_repository.ProductRepository, photoRepo product_photo_repository.ProductPhotoRepository) ProductService {
	return &productServiceImpl{productRepo: repo, productPhotoRepo: photoRepo}
}

func (s *productServiceImpl) GetAllProducts(nama, categoryID, tokoID, minHarga, maxHarga string) ([]models.Product, error) {
	return s.productRepo.FindAllWithFilter(nama, categoryID, tokoID, minHarga, maxHarga)
}

func (s *productServiceImpl) GetProductByID(id uint) (*models.Product, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}
	return product, nil
}

func (s *productServiceImpl) CreateProduct(product *models.Product, photos []*multipart.FileHeader) (*models.Product, error) {
	err := s.productRepo.Create(product)
	if err != nil {
		return nil, err
	}

	for _, photo := range photos {
		photoURL, err := utils.SaveUploadedFile(photo)
		if err != nil {
			return nil, err
		}

		photoModel := models.ProductPhoto{
			ProductID: product.ID,
			URL:       photoURL,
		}

		if err := s.productPhotoRepo.Create(&photoModel); err != nil {
			return nil, err
		}
	}

	return product, nil
}

func (s *productServiceImpl) UpdateProduct(id uint, updatedProduct *models.Product, photos []*multipart.FileHeader) (*models.Product, error) {
	existingProduct, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if existingProduct == nil {
		return nil, errors.New("product not found")
	}

	existingProduct.NamaProduct = updatedProduct.NamaProduct
	existingProduct.HargaReseller = updatedProduct.HargaReseller
	existingProduct.HargaKonsumen = updatedProduct.HargaKonsumen
	existingProduct.Stok = updatedProduct.Stok
	existingProduct.Deskripsi = updatedProduct.Deskripsi

	err = s.productRepo.Update(existingProduct)
	if err != nil {
		return nil, err
	}
	return existingProduct, nil
}

func (s *productServiceImpl) DeleteProduct(id uint) error {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("product not found")
	}

	if err := s.productPhotoRepo.DeleteByProductID(id); err != nil {
		return err
	}

	return s.productRepo.Delete(id)
}
