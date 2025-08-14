package messages

type MessageService struct {
	repository MessageRepository
}

func NewService(r MessageRepository) *MessageService {
	return &MessageService{
		repository: r,
	}
}

func (s *MessageService) Send(m Message) (*Message, error) {
	res, err := s.repository.Send(m)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *MessageService) GetAll(chat_id, user_id int) ([]Message, error) {
	res, err := s.repository.GetAll(chat_id, user_id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *MessageService) GetChatMembers(chat_id int) []int {
	return s.repository.GetChatMembers(chat_id)
}
