package chats

import (
	"fmt"

	"github.com/PandaX185/tatsumaki-chat/config"
	"github.com/PandaX185/tatsumaki-chat/domain/shared"
	"github.com/PandaX185/tatsumaki-chat/domain/users"
	"github.com/jmoiron/sqlx"
)

type ChatRepository interface {
	Create(ChatRequest) (*ChatResponse, error)
	Delete(int, int) error
	Edit(int, ChatRequest) (*ChatResponse, error)
	GetAllChats(int) ([]shared.Chat, error)
	GetChatMembers(int) (users.UserSlice, error)
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

func (r *ChatRepositoryImpl) Delete(chatId int, userId int) error {
	tx := r.db.MustBegin()

	if _, err := tx.Exec(`delete from unread_chats where cid = $1`, chatId); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(`delete from users_chats where cid = $1 and uid = $2`, chatId, userId); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(`delete from chats where id = $1`, chatId); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *ChatRepositoryImpl) Edit(chatId int, chat ChatRequest) (*ChatResponse, error) {
	tx := r.db.MustBegin()

	if chat.ChatName != "" {
		if _, err := tx.NamedExec(`update chats set chat_name = :chat_name, updated_at = now() where id = :id`, map[string]interface{}{
			"id":        chatId,
			"chat_name": chat.ChatName,
		}); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if len(chat.ChatMembers) > 0 {
		chat.ChatMembers = append(chat.ChatMembers, chat.ChatOwner)
		chat.ChatMembers = shared.RemoveDuplicateMembers(chat.ChatMembers)
		if _, err := tx.Exec(`delete from users_chats where cid = $1`, chatId); err != nil {
			tx.Rollback()
			return nil, err
		}

		for _, member := range chat.ChatMembers {
			if _, err := tx.Exec(`insert into users_chats (cid, uid) values ($1, $2)`, chatId, member); err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}
	tx.Commit()

	var res shared.Chat
	if err := r.db.Get(&res, `select * from chats where id = $1 limit 1`, chatId); err != nil {
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

func (r *ChatRepositoryImpl) GetChatMembers(chatId int) (users.UserSlice, error) {
	tx := r.db.MustBegin()

	var res users.UserSlice
	if err := tx.Select(&res, `select u.* from users u join users_chats uc on u.id = uc.uid where uc.cid = $1`, chatId); err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return res, nil
}
