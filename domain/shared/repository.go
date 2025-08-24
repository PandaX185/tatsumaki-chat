package shared

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/PandaX185/tatsumaki-chat/config"
	"github.com/jmoiron/sqlx"
)

type SharedRepository interface {
	SendMessage(Message) (*Message, error)
	GetChatMembers(int) []int
}

type SharedRepositoryImpl struct {
	db *sqlx.DB
}

func NewRepository() SharedRepository {
	return &SharedRepositoryImpl{
		db: config.DbInstance,
	}
}

func (r *SharedRepositoryImpl) GetChatMembers(chat_id int) []int {
	tx := r.db.MustBegin()
	var res []int

	if err := tx.Select(&res, `select uid from users_chats where cid = $1`, chat_id); err != nil {
		return []int{}
	}

	res = RemoveDuplicateMembers(res)
	return res
}

func (r *SharedRepositoryImpl) SendMessage(m Message) (*Message, error) {
	tx := r.db.MustBegin()
	if err := tx.Get(&Chat{}, `select * from chats where id = $1 limit 1`, m.ChatId); err != nil {
		tx.Rollback()
		return nil, err
	}

	if _, err := tx.NamedExec(`insert into messages (content, sender_id, cid) values (:content, :sender_id, :cid)`, m); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Get(&m, `select * from messages where cid = $1 and sender_id = $2 order by id desc limit 1`, m.ChatId, m.SenderId); err != nil {
		tx.Rollback()
		return nil, err
	}

	if _, err := tx.Exec(`update chats set updated_at = now() where id = $1`, m.ChatId); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Get(&m.UserName, `select username from users where id = $1`, m.SenderId); err != nil {
		tx.Rollback()
		return nil, err
	}

	rds := config.GetRedis()
	messageJson, _ := json.Marshal(m)

	for _, userId := range r.GetChatMembers(m.ChatId) {
		channelName := fmt.Sprintf("messages:%d", userId)
		if err := rds.Publish(context.Background(), channelName, string(messageJson)).Err(); err != nil {
			fmt.Printf("Error publishing message to channel %v: %v\n", channelName, err)
		}
		if _, err := tx.Exec(`update unread_chats set unread_count = unread_count + 1 where cid = $1 and uid = $2`, m.ChatId, userId); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()
	return &m, nil
}
