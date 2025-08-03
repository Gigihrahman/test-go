package user_service

import (
	"errors"
	"fmt"
	"os"
	"time"

	"test-rakamin/internal/models"
	user_repository "test-rakamin/internal/repository/user"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterUser(user *models.User) (*models.User, error)
	LoginUser(noTelp, password string) (string, error)
	GetUserProfile(userID uint) (*models.User, error)
	UpdateUserProfile(userID uint, updatedUser *models.User) (*models.User, error)
}

type userServiceImpl struct {
	userRepo user_repository.UserRepository
}

func NewUserService(repo user_repository.UserRepository) UserService {
	return &userServiceImpl{userRepo: repo}
}

func (s *userServiceImpl) RegisterUser(user *models.User) (*models.User, error) {

	existingUserByEmail, _ := s.userRepo.FindByEmail(user.Email)
	if existingUserByEmail != nil {
		return nil, errors.New("email already exists")
	}
	existingUserByNoTelp, _ := s.userRepo.FindByNoTelp(user.NoTelp)
	if existingUserByNoTelp != nil {
		return nil, errors.New("phone number already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.KataSandi), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to encrypt password")
	}
	user.KataSandi = string(hashedPassword)

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userServiceImpl) LoginUser(noTelp, password string) (string, error) {

	user, err := s.userRepo.FindByNoTelp(noTelp)
	if err != nil {
		return "", errors.New("no telp or password is wrong")
	}
	if user == nil {
		return "", errors.New("no telp or password is wrong")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.KataSandi), []byte(password)); err != nil {
		return "", errors.New("no telp or password is wrong")
	}

	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", fmt.Errorf("JWT_SECRET not set in environment")
	}
	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", errors.New("failed to create token")
	}

	return t, nil
}

func (s *userServiceImpl) GetUserProfile(userID uint) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *userServiceImpl) UpdateUserProfile(userID uint, updatedUser *models.User) (*models.User, error) {
	existingUser, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
		return nil, errors.New("user not found")
	}

	existingUser.Nama = updatedUser.Nama
	existingUser.Email = updatedUser.Email
	if updatedUser.KataSandi != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.KataSandi), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("failed to encrypt new password")
		}
		existingUser.KataSandi = string(hashedPassword)
	}
	existingUser.NoTelp = updatedUser.NoTelp
	existingUser.TanggalLahir = updatedUser.TanggalLahir
	existingUser.Pekerjaan = updatedUser.Pekerjaan
	existingUser.IDProvinsi = updatedUser.IDProvinsi
	existingUser.IDKota = updatedUser.IDKota

	err = s.userRepo.Update(existingUser)
	if err != nil {
		return nil, err
	}
	return existingUser, nil
}
