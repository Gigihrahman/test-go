package toko_service

import (
	"errors"
	"mime/multipart"
	"test-rakamin/internal/models"
	toko_repository "test-rakamin/internal/repository/toko"
	"test-rakamin/utils"
)

type TokoService interface {
	GetTokoByUserID(userID uint) (*models.Toko, error)
	GetAllToko() ([]models.Toko, error)
	GetTokoByID(id uint) (*models.Toko, error)
	UpdateToko(id uint, namaToko string, photo *multipart.FileHeader) (*models.Toko, error)
}

type tokoServiceImpl struct {
	tokoRepo toko_repository.TokoRepository
}

func NewTokoService(repo toko_repository.TokoRepository) TokoService {
	return &tokoServiceImpl{tokoRepo: repo}
}

func (s *tokoServiceImpl) GetTokoByUserID(userID uint) (*models.Toko, error) {
	toko, err := s.tokoRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if toko == nil {
		return nil, errors.New("toko not found")
	}
	return toko, nil
}

func (s *tokoServiceImpl) GetAllToko() ([]models.Toko, error) {
	return s.tokoRepo.FindAll()
}

func (s *tokoServiceImpl) GetTokoByID(id uint) (*models.Toko, error) {
	toko, err := s.tokoRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if toko == nil {
		return nil, errors.New("toko not found")
	}
	return toko, nil
}

func (s *tokoServiceImpl) UpdateToko(id uint, namaToko string, photo *multipart.FileHeader) (*models.Toko, error) {
	existingToko, err := s.tokoRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if existingToko == nil {
		return nil, errors.New("toko not found")
	}

	existingToko.NamaToko = namaToko

	if photo != nil {

		photoURL, err := utils.SaveUploadedFile(photo)
		if err != nil {
			return nil, err
		}
		existingToko.URLFotoToko = photoURL
	}

	err = s.tokoRepo.Update(existingToko)
	if err != nil {
		return nil, err
	}

	return existingToko, nil
}
