package messages

import (
	"github.com/PandaX185/tatsumaki-chat/config"
	"github.com/PandaX185/tatsumaki-chat/domain/chats"
	"github.com/jmoiron/sqlx"
)

type MessageRepository interface {
	Send(Message) (string, error)
	GetAll(int, int) ([]Message, error)
}

type MessageRepositoryImpl struct {
	db *sqlx.DB
}

func NewRepository() MessageRepository {
	return &MessageRepositoryImpl{
		db: config.DbInstance,
	}
}

func (r *MessageRepositoryImpl) Send(m Message) (string, error) {
	tx := r.db.MustBegin()
	if err := tx.Get(&chats.Chat{}, `select * from chats where id = $1 limit 1`, m.ChatId); err != nil {
		tx.Rollback()
		return "", err
	}

	if _, err := tx.NamedExec(`insert into messages (content, sender_id, cid) values (:content, :sender_id, :cid)`, m); err != nil {
		tx.Rollback()
		return "", err
	}

	tx.Commit()
	return "Message sent", nil
}

func (r *MessageRepositoryImpl) GetAll(chat_id, user_id int) ([]Message, error) {
	tx := r.db.MustBegin()
	var res []Message

	if err := tx.Select(&res, `select messages.* from messages join users_chats on messages.cid = users_chats.cid where users_chats.cid = $1 and uid = $2`, chat_id, user_id); err != nil {
		return nil, err
	}

	return res, nil
}
