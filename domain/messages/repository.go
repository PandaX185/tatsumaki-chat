package messages

import (
	"strings"

	"github.com/PandaX185/tatsumaki-chat/config"
	"github.com/PandaX185/tatsumaki-chat/domain/shared"
	"github.com/jmoiron/sqlx"
)

type MessageRepository interface {
	GetAll(int, int) ([]shared.Message, error)
	GetUnreadMessagesCount(int) ([]UnreadMessagesCount, error)
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

func (r *MessageRepositoryImpl) GetUnreadMessagesCount(user_id int) ([]UnreadMessagesCount, error) {
	var count []UnreadMessagesCount

	if err := r.db.Select(&count, `
		SELECT unread_count, cid
		FROM unread_chats
		WHERE uid = $1
		GROUP BY cid, unread_count
	`, user_id); err != nil {
		if !strings.Contains(err.Error(), "no rows") {
			return nil, err
		}
	}

	return count, nil
}
