package repo

import "database/sql"

type Repositories struct {
	ItemRep *ItemRepository
	UserRep *UserRepository
}

func InitializeRepositories(db *sql.DB) (*Repositories){
	itemRep := NewItemRepository(db)
	userRep := NewUserRepository(db)

	return &Repositories {
		ItemRep: itemRep,
		UserRep: userRep,
	}
}