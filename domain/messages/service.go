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
	res, err := s.sharedRepo.SendMessage(m)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *MessageService) GetAll(chat_id, user_id int) ([]shared.Message, error) {
	res, err := s.repository.GetAll(chat_id, user_id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *MessageService) GetUnreadMessagesCount(user_id int) ([]UnreadMessagesCount, error) {
	count, err := s.repository.GetUnreadMessagesCount(user_id)
	if err != nil {
		return nil, err
	}
	return count, nil
}
