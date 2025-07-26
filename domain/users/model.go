package users

import (
	"database/sql"
	"strconv"
	"time"
)

type User struct {
	Id        int          `json:"id" db:"id"`
	FullName  string       `json:"full_name" db:"full_name"`
	UserName  string       `json:"user_name" db:"user_name"`
	Password  string       `json:"password" db:"password"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt time.Time    `json:"updated_at" db:"updated_at"`
	DeletedAt sql.NullTime `json:"-" db:"deleted_at"`
}

func (u User) ToApiResponse() map[string]string {
	return map[string]string{
		"id":       strconv.FormatInt(int64(u.Id), 10),
		"username": u.UserName,
		"fullname": u.FullName,
	}
}
