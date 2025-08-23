package messages

import (
	"github.com/PandaX185/tatsumaki-chat/config"
	"github.com/PandaX185/tatsumaki-chat/domain/shared"
	"github.com/jmoiron/sqlx"
)

type MessageRepository interface {
	GetAll(int, int) ([]shared.Message, error)
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

	if err := tx.Select(&res, `select messages.* from messages join users_chats on messages.cid = users_chats.cid where users_chats.cid = $1 and uid = $2`, chat_id, user_id); err != nil {
		return nil, err
	}

	return res, nil
}
