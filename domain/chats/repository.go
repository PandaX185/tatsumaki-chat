package chats

import (
	"github.com/PandaX185/tatsumaki-chat/config"
	"github.com/jmoiron/sqlx"
)

type ChatRepository interface {
	Create(ChatRequest) (*ChatResponse, error)
	GetAllChats(int) ([]Chat, error)
}

type ChatRepositoryImpl struct {
	db *sqlx.DB
}

func NewRepository() ChatRepository {
	return &ChatRepositoryImpl{
		db: config.DbInstance,
	}
}

func (c *ChatRepositoryImpl) Create(chat ChatRequest) (*ChatResponse, error) {
	tx := c.db.MustBegin()

	if _, err := tx.NamedExec(`insert into chats (chat_name, chat_owner) values (:chat_name, :chat_owner)`, chat); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Get(&chat, `select * from chats where chat_owner = $1 and chat_name = $2 order by created_at desc limit 1`, chat.ChatOwner, chat.ChatName); err != nil {
		tx.Rollback()
		return nil, err
	}

	if _, err := tx.NamedExec(`insert into users_chats (cid, uid) values (:id, :chat_owner)`, chat); err != nil {
		tx.Rollback()
		return nil, err
	}

	var res Chat
	if err := tx.Get(&res, `select * from chats where chat_owner = $1 limit 1`, chat.ChatOwner); err != nil {
		tx.Rollback()
		return nil, err
	}

	query, args, err := sqlx.In(`SELECT id FROM users WHERE user_name IN (?)`, chat.ChatMembers)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	query = tx.Rebind(query)
	var userIds []int
	if err := tx.Select(&userIds, query, args...); err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, userId := range userIds {
		if _, err := tx.Exec(`INSERT INTO users_chats (uid, cid) VALUES ($1, $2)`, userId, res.Id); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()
	return &ChatResponse{
		Id:          res.Id,
		ChatName:    res.ChatName,
		ChatOwner:   res.ChatOwner,
		ChatMembers: chat.ChatMembers,
		CreatedAt:   res.CreatedAt,
	}, nil
}

func (c *ChatRepositoryImpl) GetAllChats(userId int) ([]Chat, error) {
	tx := c.db.MustBegin()

	var res []Chat
	if err := tx.Select(&res, `select chats.* from chats join users_chats on cid = chats.id where uid = $1 order by created_at desc`, userId); err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return res, nil
}
