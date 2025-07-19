package chats

type ChatService struct {
	repository ChatRepository
}

func NewService(r ChatRepository) *ChatService {
	return &ChatService{
		repository: r,
	}
}

func (s *ChatService) Create(chat Chat) (*Chat, error) {
	res, err := s.repository.Create(chat)
	if err != nil {
		return nil, err
	}

	return res, nil
}
