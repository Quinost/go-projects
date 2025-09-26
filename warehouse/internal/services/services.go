package services

import (
	"warehouse/internal/cfg"
	"warehouse/internal/repositories"
)

type Services struct {
	ItemService *ItemService
	JWTService  *JWTService
}

func InitializeServices(cfg *cfg.Config, repo *repositories.Repositories) *Services {
	itemService := NewItemService(repo.ItemRep)
	jwtService := NewJWTService(&cfg.Auth)

	return &Services{
		ItemService: itemService,
		JWTService:  jwtService,
	}
}
