package repositories

import (
	"database/sql"
	"warehouse/internal/models"

	"github.com/google/uuid"
)

type ItemRepository struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{
		db: db,
	}
}

func (r *ItemRepository) GetById(id uuid.UUID) (*models.Item, error) {
	query := `SELECT id, name, description FROM items WHERE id = $1`
	item := &models.Item{}
	err := r.db.QueryRow(query, id).Scan(&item.Id, &item.Name, &item.Description)
	return item, err
}

func (r *ItemRepository) GetAll(filter string, page, limit int) ([]models.Item, error) {
	query := `
	SELECT id, name, description
	FROM items
	WHERE name ILIKE $1
	LIMIT $2 OFFSET $3
	`
	searchTerm := "%" + filter + "%"
	rows, err := r.db.Query(query, searchTerm, limit, page)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var items []models.Item
	for rows.Next() {
		var item models.Item
		rows.Scan(&item.Id, &item.Name, &item.Description)
		items = append(items, item)
	}
	return items, nil
}

func (r *ItemRepository) Add(item *models.Item) error {
	query := `INSERT INTO items (id, name, description) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, item.Id, item.Name, item.Description)
	return err
}

func (r *ItemRepository) Update(item *models.Item) error {
	query := `
        UPDATE items
        SET name = $1, description = $2
        WHERE id = $3
    ` 
	_, err := r.db.Exec(query, item.Name, item.Description, item.Id)
	return err
}

func (r *ItemRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM items WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
