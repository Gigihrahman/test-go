package category_service

import (
	"errors"
	"test-rakamin/internal/models"
	category_repository "test-rakamin/internal/repository/category"
)

type CategoryService interface {
	GetAllCategories() ([]models.Category, error)
	GetCategoryByID(id uint) (*models.Category, error)
	CreateCategory(category *models.Category) (*models.Category, error)
	UpdateCategory(id uint, updatedCategory *models.Category) (*models.Category, error)
	DeleteCategory(id uint) error
}

type categoryServiceImpl struct {
	categoryRepo category_repository.CategoryRepository
}

func NewCategoryService(repo category_repository.CategoryRepository) CategoryService {
	return &categoryServiceImpl{categoryRepo: repo}
}

func (s *categoryServiceImpl) GetAllCategories() ([]models.Category, error) {
	return s.categoryRepo.FindAll()
}

func (s *categoryServiceImpl) GetCategoryByID(id uint) (*models.Category, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("category not found")
	}
	return category, nil
}

func (s *categoryServiceImpl) CreateCategory(category *models.Category) (*models.Category, error) {
	err := s.categoryRepo.Create(category)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *categoryServiceImpl) UpdateCategory(id uint, updatedCategory *models.Category) (*models.Category, error) {
	existingCategory, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if existingCategory == nil {
		return nil, errors.New("category not found")
	}
	existingCategory.NamaCategory = updatedCategory.NamaCategory
	err = s.categoryRepo.Update(existingCategory)
	if err != nil {
		return nil, err
	}
	return existingCategory, nil
}

func (s *categoryServiceImpl) DeleteCategory(id uint) error {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return err
	}
	if category == nil {
		return errors.New("category not found")
	}
	return s.categoryRepo.Delete(id)
}
