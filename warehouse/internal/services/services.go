package services

import (
	"warehouse/internal/repositories"
)


type Services struct {
	ItemService *ItemService
}

func InitializeServices(repo *repositories.Repositories) (*Services){
	itemService := NewItemService(repo.ItemRep)

	return &Services {
		ItemService: itemService,
	}
}