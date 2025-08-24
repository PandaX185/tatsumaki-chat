package messages

import (
	"context"
	"fmt"

	"github.com/PandaX185/tatsumaki-chat/config"
	"github.com/PandaX185/tatsumaki-chat/domain/shared"
	"github.com/jmoiron/sqlx"
)

type MessageRepository interface {
	GetAll(int, int) ([]shared.Message, error)
	MarkAsRead(int, int) error
}

type MessageRepositoryImpl struct {
	db *sqlx.DB
}

func NewRepository() MessageRepository {
	return &MessageRepositoryImpl{
		db: config.DbInstance,
	}
}

func (r *MessageRepositoryImpl) GetAll(chat_id, user_id int) ([]shared.Message, error) {
	tx := r.db.MustBegin()
	var res []shared.Message

	if err := tx.Select(&res, `
	       select messages.*, users.username 
	       from messages 
	       join users on messages.sender_id = users.id 
	       join users_chats on messages.cid = users_chats.cid 
	       where users_chats.cid = $1 and uid = $2
		   order by messages.created_at
       `, chat_id, user_id); err != nil {
		return nil, err
	}

	return res, nil
}

func (r *MessageRepositoryImpl) MarkAsRead(chat_id, user_id int) error {
	tx := r.db.MustBegin()

	if _, err := tx.Exec(`update unread_chats set unread_count = 0 where cid = $1 and uid = $2`, chat_id, user_id); err != nil {
		tx.Rollback()
		return err
	}

	if err := config.GetRedis().Publish(context.Background(), fmt.Sprintf("read:%d", user_id), chat_id).Err(); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
