package services

import (
	"fmt"
	"log"
	"warehouse/internal/models"
	"warehouse/internal/repositories"
	"warehouse/internal/validator"

	"github.com/google/uuid"
)

type ItemService struct {
	repo *repo.ItemRepository
}

func NewItemService(repo *repo.ItemRepository) *ItemService {
	return &ItemService{
		repo: repo,
	}
}

func (s *ItemService) GetById(uid uuid.UUID) (*models.Item, error) {
	item, err := s.repo.GetById(uid)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("item not found")
	}

	return item, nil
}

func (s *ItemService) GetAll(filter string, page int, limit int) ([]models.Item, error) {
	items, err := s.repo.GetAll(filter, page, limit)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("error while loading list")
	}
	return items, nil
}

func (s *ItemService) Add(item *models.Item) (*uuid.UUID, error) {
	if err := validator.Validate(item); err != nil {
		return nil, err
	}

	exists := s.repo.CheckIfNameExist(item.Name)
	if exists {
		return nil, fmt.Errorf("name exist")
	}

	item.Id = uuid.New()
	if err := s.repo.Add(item); err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed adding item")
	}

	return &item.Id, nil
}

func (s *ItemService) Update(item models.Item) error {
	if err := validator.Validate(item); err != nil {
		return err
	}

	extItem, err := s.repo.GetById(item.Id)
	if err != nil || extItem == nil {
		log.Println(err)
		return fmt.Errorf("item not found")
	}

	if err := s.repo.Update(&item); err != nil {
		log.Println(err)
		return fmt.Errorf("failed to update item")
	}

	return nil
}

func (s *ItemService) Delete(uid uuid.UUID) error {
	if item, err := s.repo.GetById(uid); err != nil || item == nil {
		return fmt.Errorf("item not found")
	}

	err := s.repo.Delete(uid)
	if err != nil {
		return fmt.Errorf("failed while deleting item %s", uid)
	}
	return nil
}
