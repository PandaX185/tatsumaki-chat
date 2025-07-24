package chats

import "time"

type Chat struct {
	Id                int       `json:"id" db:"id"`
	ChatName          string    `json:"chat_name" db:"chat_name"`
	ChatOwner         int       `json:"chat_owner" db:"chat_owner"`
	LastMessage       string    `json:"last_message" db:"-"`
	LastMessageSender string    `json:"sender" db:"-"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
}
