package chats

import "time"

type ChatRequest struct {
	Id          int       `json:"-" db:"id"`
	ChatName    string    `json:"chat_name" db:"chat_name"`
	ChatOwner   int       `json:"-" db:"chat_owner"`
	ChatMembers []int     `json:"chat_members" db:"chat_members"`
	CreatedAt   time.Time `json:"-" db:"created_at"`
	UpdatedAt   time.Time `json:"-" db:"updated_at"`
}

type ChatResponse struct {
	Id          int       `json:"id"`
	ChatName    string    `json:"chat_name"`
	ChatOwner   int       `json:"chat_owner"`
	ChatMembers []int     `json:"chat_members"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
