package repositories

import "database/sql"

type Repositories struct {
	ItemRep *ItemRepository
}

func InitializeRepositories(db *sql.DB) (*Repositories){
	itemRep := NewItemRepository(db)

	return &Repositories {
		ItemRep: itemRep,
	}
}