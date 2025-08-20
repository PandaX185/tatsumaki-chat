package users

import (
	"github.com/PandaX185/tatsumaki-chat/config"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Save(User) (*User, error)
	GetByExactUserName(string) (*User, error)
	SearchByUserName(string) (UserSlice, error)
	Login(string, string) (*User, error)
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
	if err = tx.Get(&res, `select * from users where user_name = $1 limit 1`, user.UserName); err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return &res, nil
}

func (r *UserRepositoryImpl) SearchByUserName(username string) (UserSlice, error) {
	var res UserSlice
	if err := r.db.Select(&res, `select * from users where user_name ILIKE $1`, "%"+username+"%"); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *UserRepositoryImpl) GetByExactUserName(username string) (*User, error) {
	var res User
	if err := r.db.Get(&res, `select * from users where user_name = $1`, username); err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *UserRepositoryImpl) Login(username, password string) (*User, error) {
	var res User
	if err := r.db.Get(&res, `select * from users where user_name = $1 and password = $2 limit 1`, username, password); err != nil {
		return nil, err
	}
	return &res, nil
}
