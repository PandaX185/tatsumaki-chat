package user

import (
	"github.com/PandaX185/tatsumaki-chat/config"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Save(user User) (*User, error)
	GetByUserName(username string) (*User, error)
}

type UserRepositoryImpl struct {
	Db *sqlx.DB
}

func NewRepository() UserRepository {
	return &UserRepositoryImpl{
		Db: config.DbInstance,
	}
}

func (r *UserRepositoryImpl) Save(user User) (*User, error) {
	tx := r.Db.MustBegin()
	_, err := tx.NamedExec(`insert into users (user_name, full_name, password) values (:user_name, :full_name, :password)`, user)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var res User
	if err = tx.Get(&res, `select * from users where user_name = $1`, user.UserName); err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return &res, nil
}

func (r *UserRepositoryImpl) GetByUserName(username string) (*User, error) {
	tx := r.Db.MustBegin()

	var res User
	if err := tx.Get(&res, `select * from users where user_name = $1`, username); err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return &res, nil
}
