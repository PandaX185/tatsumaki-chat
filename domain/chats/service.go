package chats

import (
	"github.com/PandaX185/tatsumaki-chat/domain/shared"
	"github.com/PandaX185/tatsumaki-chat/domain/users"
)

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

func (s *ChatService) Delete(chatId int, userId int) error {
	return s.repository.Delete(chatId, userId)
}

func (s *ChatService) Edit(chatId int, chat ChatRequest) (*ChatResponse, error) {
	return s.repository.Edit(chatId, chat)
}

func (s *ChatService) GetChatMembers(chatId int) (users.UserSlice, error) {
	return s.repository.GetChatMembers(chatId)
}
