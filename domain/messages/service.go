package messages

import "github.com/PandaX185/tatsumaki-chat/domain/shared"

type MessageService struct {
	repository MessageRepository
	sharedRepo shared.SharedRepository
}

func NewService(r MessageRepository, s shared.SharedRepository) *MessageService {
	return &MessageService{
		repository: r,
		sharedRepo: s,
	}
}

func (s *MessageService) Send(m shared.Message) (*shared.Message, error) {
	return s.sharedRepo.SendMessage(m)
}

func (s *MessageService) GetAll(chat_id, user_id int) ([]shared.Message, error) {
	return s.repository.GetAll(chat_id, user_id)
}

func (s *MessageService) GetUnreadMessagesCount(user_id int) ([]shared.UnreadMessagesCount, error) {
	return s.sharedRepo.GetUnreadMessagesCount(user_id)
}

func (s *MessageService) MarkAsRead(chat_id, user_id int) error {
	return s.repository.MarkAsRead(chat_id, user_id)
}
