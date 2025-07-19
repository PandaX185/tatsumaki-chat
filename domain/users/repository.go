package users

import (
	"github.com/PandaX185/tatsumaki-chat/config"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Save(User) (*User, error)
	GetByUserName(string) (*User, error)
}

type UserRepositoryImpl struct {
	db *sqlx.DB
}

func NewRepository() UserRepository {
	return &UserRepositoryImpl{
		db: config.DbInstance,
	}
}

func (r *UserRepositoryImpl) Save(user User) (*User, error) {
	tx := r.db.MustBegin()
	_, err := tx.NamedExec(`insert into users (user_name, full_name, password) values (:user_name, :full_name, :password)`, user)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var res User
	if err = tx.Get(&res, `select * from users where user_name = :user_name limit 1`, user); err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return &res, nil
}

func (r *UserRepositoryImpl) GetByUserName(username string) (*User, error) {
	tx := r.db.MustBegin()

	var res User
	if err := tx.Get(&res, `select * from users where user_name = $1 limit 1`, username); err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return &res, nil
}
