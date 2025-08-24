package shared

import "time"

type Chat struct {
	Id        int    `json:"id" db:"id"`
	ChatName  string `json:"chat_name" db:"chat_name"`
	ChatOwner int    `json:"chat_owner" db:"chat_owner"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Message struct {
	Id       int    `json:"id" db:"id"`
	ChatId   int    `json:"chat_id" db:"cid"`
	SenderId int    `json:"sender_id" db:"sender_id"`
	UserName string `json:"username" db:"username"`
	Content  string `json:"content" db:"content"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UnreadMessagesCount struct {
	Count  int `json:"unread_count" db:"unread_count"`
	ChatId int `json:"chat_id" db:"cid"`
}
