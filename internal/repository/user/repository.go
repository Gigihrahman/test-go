package user_repository

import (
	"test-rakamin/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindAll() ([]models.User, error)
	FindByID(id uint) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByNoTelp(noTelp string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
}

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepositoryImpl) FindAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepositoryImpl) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

func (r *userRepositoryImpl) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

func (r *userRepositoryImpl) FindByNoTelp(noTelp string) (*models.User, error) {
	var user models.User
	err := r.db.Where("no_telp = ?", noTelp).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

func (r *userRepositoryImpl) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}
