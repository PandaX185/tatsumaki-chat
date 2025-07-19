package chats

import (
	"github.com/PandaX185/tatsumaki-chat/config"
	"github.com/jmoiron/sqlx"
)

type ChatRepository interface {
	Create(Chat) (*Chat, error)
}

type ChatRepositoryImpl struct {
	db *sqlx.DB
}

func NewRepository() ChatRepository {
	return &ChatRepositoryImpl{
		db: config.DbInstance,
	}
}

func (c *ChatRepositoryImpl) Create(chat Chat) (*Chat, error) {
	tx := c.db.MustBegin()

	if _, err := tx.NamedExec(`insert into chats (chat_name, chat_owner) values (:chat_name, :chat_owner)`, chat); err != nil {
		tx.Rollback()
		return nil, err
	}

	var res Chat
	if err := tx.Get(&res, `select * from chats where chat_owner = $1 limit 1`, chat.ChatOwner); err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return &res, nil
}
