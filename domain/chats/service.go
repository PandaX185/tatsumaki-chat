package chats

import "github.com/PandaX185/tatsumaki-chat/domain/shared"

type ChatService struct {
	repository ChatRepository
}

func NewService(r ChatRepository) *ChatService {
	return &ChatService{
		repository: r,
	}
}

func (s *ChatService) Create(chat ChatRequest) (*ChatResponse, error) {
	return s.repository.Create(chat)
}

func (s *ChatService) GetAllChats(userId int) ([]shared.Chat, error) {
	return s.repository.GetAllChats(userId)
}
