package websocket

import (
	"realtime-chat-backend/pkg/models"
	"sync"
)

type MessageStore struct {
	messages []models.Message
	lock     sync.RWMutex
}

func NewMessageStore() *MessageStore {
	return &MessageStore{
		messages: make([]models.Message, 0),
	}
}

func (s *MessageStore) AddMessage(msg models.Message) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.messages = append(s.messages, msg)
}

func (s *MessageStore) GetLastMessages(n int) []models.Message {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if len(s.messages) < n {
		return s.messages
	}

	return s.messages[len(s.messages)-n:]
}
