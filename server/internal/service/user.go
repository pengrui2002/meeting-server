package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"meeting-server/internal/model"
	"meeting-server/internal/repository"
)

var ErrUserExists = errors.New("user already exists")
var ErrInvalidCredentials = errors.New("invalid credentials")

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(username, password, displayName string) (*model.User, error) {
	existing, _ := s.repo.GetByUsername(username)
	if existing != nil {
		return nil, ErrUserExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username:     username,
		PasswordHash: string(hash),
		DisplayName:  displayName,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(username, password string) (*model.User, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (s *UserService) GetByID(id string) (*model.User, error) {
	return s.repo.GetByID(id)
}