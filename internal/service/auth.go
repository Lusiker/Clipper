package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/lusiker/clipper/internal/model"
	"github.com/lusiker/clipper/internal/pkg/crypto"
	"github.com/lusiker/clipper/internal/repository"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(username, password string) (*model.User, error) {
	if existing, _ := s.userRepo.FindByUsername(username); existing != nil {
		return nil, errors.New("username already exists")
	}

	hashedPassword, err := crypto.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:       uuid.New().String(),
		Username: username,
		Password: hashedPassword,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}

func (s *AuthService) Login(username, password string) (*model.User, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	if !crypto.CheckPassword(password, user.Password) {
		return nil, errors.New("invalid password")
	}

	user.Password = ""
	return user, nil
}

func (s *AuthService) GetUserByID(userID string) (*model.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	user.Password = ""
	return user, nil
}