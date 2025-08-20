package chats

type ChatService struct {
	repository ChatRepository
}

func NewService(r ChatRepository) *ChatService {
	return &ChatService{
		repository: r,
	}
}

func (s *ChatService) Create(chat ChatRequest) (*ChatResponse, error) {
	res, err := s.repository.Create(chat)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *ChatService) GetAllChats(userId int) ([]Chat, error) {
	res, err := s.repository.GetAllChats(userId)
	if err != nil {
		return nil, err
	}

	return res, nil
}
