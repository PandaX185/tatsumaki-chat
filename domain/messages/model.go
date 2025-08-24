package messages

type UnreadMessagesCount struct {
	Count  int `json:"unread_count" db:"unread_count"`
	ChatID int `json:"chat_id" db:"cid"`
}
