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

type UserSlice []User

func (u UserSlice) ToShortUserResponse() []map[string]string {
	var response []map[string]string
	for _, user := range u {
		response = append(response, user.ToShortUserResponse())
	}
	return response
}

func (u User) ToShortUserResponse() map[string]string {
	return map[string]string{
		"id":       strconv.FormatInt(int64(u.Id), 10),
		"username": u.UserName,
		"fullname": u.FullName,
	}
}

func (u User) ToDetailedUserResponse() map[string]interface{} {
	return map[string]interface{}{
		"id":         strconv.FormatInt(int64(u.Id), 10),
		"full_name":  u.FullName,
		"user_name":  u.UserName,
		"created_at": u.CreatedAt.Format(time.RFC3339),
		"updated_at": u.UpdatedAt.Format(time.RFC3339),
	}
}
