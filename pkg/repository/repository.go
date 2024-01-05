package repository

import (
	"database/sql"
	"github/PandaX185/tatsumaki-chat/pkg/models"
)

type Repository interface {
	GetUser(id int) (models.User, error)
	CreateUser(user models.User) error
}

type repository struct {
	db *sql.DB
}

func (r *repository) GetUser(id int) (models.User, error) {
	var user models.User
	err := r.db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *repository) CreateUser(user models.User) error {
	_, err := r.db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, user.Password)

	if err != nil {
		return err
	}

	return nil
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}
