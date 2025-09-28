package services

import (
	"warehouse/internal/cfg"
	"warehouse/internal/repositories"
)

type Services struct {
	ItemService *ItemService
	JWTService  *JWTService
	UserService *UserService
}

func InitializeServices(cfg *cfg.Config, repo *repo.Repositories) *Services {
	itemService := NewItemService(repo.ItemRep)
	jwtService := NewJWTService(&cfg.Auth)
	userService := NewUserService(repo.UserRep, jwtService)

	return &Services{
		ItemService: itemService,
		JWTService:  jwtService,
		UserService: userService,
	}
}
