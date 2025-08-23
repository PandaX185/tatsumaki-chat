package chats

import (
	"fmt"

	"github.com/PandaX185/tatsumaki-chat/config"
	"github.com/PandaX185/tatsumaki-chat/domain/shared"
	"github.com/jmoiron/sqlx"
)

type ChatRepository interface {
	Create(ChatRequest) (*ChatResponse, error)
	GetAllChats(int) ([]shared.Chat, error)
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

	chat.ChatMembers = append(chat.ChatMembers, chat.ChatOwner)
	chat.ChatMembers = shared.RemoveDuplicateMembers(chat.ChatMembers)
	for _, member := range chat.ChatMembers {
		fmt.Printf("member: %v\n", member)
		if _, err := tx.Exec(`insert into users_chats (cid, uid) values ($1, $2)`, chat.Id, member); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	var res shared.Chat
	if err := tx.Get(&res, `select * from chats where chat_owner = $1 and chat_name = $2 limit 1`, chat.ChatOwner, chat.ChatName); err != nil {
		tx.Rollback()
		return nil, err
	}

	query, args, err := sqlx.In(`SELECT id FROM users WHERE username IN (?)`, chat.ChatMembers)
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

	_, err = shared.NewRepository().SendMessage(shared.Message{
		ChatId:   res.Id,
		SenderId: chat.ChatOwner,
		Content:  fmt.Sprintf("Welcome to the new chat! (%s)", res.ChatName),
	})
	if err != nil {
		return nil, err
	}

	return &ChatResponse{
		Id:          res.Id,
		ChatName:    res.ChatName,
		ChatOwner:   res.ChatOwner,
		ChatMembers: chat.ChatMembers,
		CreatedAt:   res.CreatedAt,
	}, nil
}

func (c *ChatRepositoryImpl) GetAllChats(userId int) ([]shared.Chat, error) {
	tx := c.db.MustBegin()

	var res []shared.Chat
	if err := tx.Select(&res, `select chats.* from chats join users_chats on cid = chats.id where uid = $1 order by updated_at desc`, userId); err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return res, nil
}
