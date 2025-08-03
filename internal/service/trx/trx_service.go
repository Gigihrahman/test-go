package trx_service

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"test-rakamin/internal/models"
	product_repository "test-rakamin/internal/repository/product"
	trx_repository "test-rakamin/internal/repository/trx"
)

type TrxService interface {
	GetAllTrxByUserID(userID uint) ([]models.Trx, error)
	GetTrxByID(id uint, userID uint) (*models.Trx, error)
	CreateTrx(userID uint, payload *models.TrxPayload) (*models.Trx, error)
}

type trxServiceImpl struct {
	trxRepo     trx_repository.TrxRepository
	productRepo product_repository.ProductRepository
}

func NewTrxService(repo trx_repository.TrxRepository, productRepo product_repository.ProductRepository) TrxService {
	return &trxServiceImpl{trxRepo: repo, productRepo: productRepo}
}

func (s *trxServiceImpl) GetAllTrxByUserID(userID uint) ([]models.Trx, error) {
	return s.trxRepo.FindByUserID(userID)
}

func (s *trxServiceImpl) GetTrxByID(id uint, userID uint) (*models.Trx, error) {
	trx, err := s.trxRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		return nil, err
	}
	if trx == nil {
		return nil, errors.New("transaction not found")
	}
	return trx, nil
}

func (s *trxServiceImpl) CreateTrx(userID uint, payload *models.TrxPayload) (*models.Trx, error) {
	var totalHarga int
	var detailTrxList []models.DetailTrx

	for _, item := range payload.DetailTrx {
		product, err := s.productRepo.FindByID(item.ProductID)
		if err != nil || product == nil {
			return nil, fmt.Errorf("product with ID %d not found", item.ProductID)
		}

		if product.Stok < item.Kuantitas {
			return nil, fmt.Errorf("stock for product %s is insufficient", product.NamaProduct)
		}

		itemTotal := product.HargaKonsumen * item.Kuantitas
		totalHarga += itemTotal

		detailTrxList = append(detailTrxList, models.DetailTrx{

			IDToko:     product.IDToko,
			Kuantitas:  item.Kuantitas,
			HargaTotal: itemTotal,
		})
	}

	rand.Seed(time.Now().UnixNano())
	kodeInvoice := fmt.Sprintf("INV-%d", rand.Intn(1000000))

	newTrx := &models.Trx{
		IDUser:           userID,
		AlamatPengiriman: payload.AlamatKirim,
		KodeInvoice:      kodeInvoice,
		MethodBayar:      payload.MethodBayar,
		HargaTotal:       totalHarga,
		DetailTrx:        detailTrxList,
	}

	err := s.trxRepo.Create(newTrx)
	if err != nil {
		return nil, err
	}

	return newTrx, nil
}
