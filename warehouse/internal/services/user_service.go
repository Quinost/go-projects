package services

import (
	"fmt"
	"warehouse/internal/models"
	"warehouse/internal/repositories"
)

type UserService struct {
	repo *repo.UserRepository
	jwt  *JWTService
}

func NewUserService(repo *repo.UserRepository, jwt *JWTService) *UserService {
	return &UserService{
		repo: repo,
		jwt: jwt,
	}
}

func (s *UserService) CheckAndGetUser(username, password string) (*models.User, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil || user == nil {
		return nil, fmt.Errorf("user %s not found", username)
	}

	if !s.jwt.IsPasswordCorrect(password, user.Password) {
		return nil, fmt.Errorf("user %s not found", username)
	}

	return user, nil
}
