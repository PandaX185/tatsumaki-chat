package chats

import "time"

type Chat struct {
	Id        int       `json:"id" db:"id"`
	ChatName  string    `json:"chat_name" db:"chat_name"`
	ChatOwner int       `json:"chat_owner" db:"chat_owner"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type ChatRequest struct {
	Id          int       `json:"-" db:"id"`
	ChatName    string    `json:"chat_name" db:"chat_name"`
	ChatOwner   int       `json:"-" db:"chat_owner"`
	ChatMembers []string  `json:"chat_members" db:"chat_members"`
	CreatedAt   time.Time `json:"-" db:"created_at"`
}

type ChatResponse struct {
	Id          int       `json:"id"`
	ChatName    string    `json:"chat_name"`
	ChatOwner   int       `json:"chat_owner"`
	ChatMembers []string  `json:"chat_members"`
	CreatedAt   time.Time `json:"created_at"`
}
