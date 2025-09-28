package repo

import (
	"database/sql"
	"warehouse/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetByUsername(name string) (*models.User, error) {
	query := `SELECT id, username, password FROM users WHERE username = $1`
	user := &models.User{}
	err := r.db.QueryRow(query, name).Scan(&user.Id, &user.Username, &user.Password)
	return user, err
}