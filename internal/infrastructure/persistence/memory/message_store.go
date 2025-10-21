package memory

import (
	"realtime-chat-backend/internal/domain/entity"
	"realtime-chat-backend/internal/domain/repository"
	"sync"
)

type messageCache struct {
	messages []*entity.Message
	lock     sync.RWMutex
	maxSize  int
}

// NewMessageCache creates a new in-memory message cache
func NewMessageCache(maxSize int) repository.MessageCache {
	if maxSize <= 0 {
		maxSize = 1000
	}
	return &messageCache{
		messages: make([]*entity.Message, 0),
		maxSize:  maxSize,
	}
}

func (c *messageCache) Add(message *entity.Message) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.messages = append(c.messages, message)

	// Keep only the most recent maxSize messages
	if len(c.messages) > c.maxSize {
		c.messages = c.messages[len(c.messages)-c.maxSize:]
	}
}

func (c *messageCache) GetRecent(n int) []*entity.Message {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if len(c.messages) == 0 {
		return []*entity.Message{}
	}

	if n <= 0 || n > len(c.messages) {
		n = len(c.messages)
	}

	start := len(c.messages) - n
	result := make([]*entity.Message, n)
	copy(result, c.messages[start:])

	return result
}

func (c *messageCache) GetByRoomID(roomID string, limit int) []*entity.Message {
	c.lock.RLock()
	defer c.lock.RUnlock()

	var result []*entity.Message

	// Iterate from newest to oldest
	for i := len(c.messages) - 1; i >= 0 && len(result) < limit; i-- {
		if c.messages[i].RoomID == roomID {
			result = append(result, c.messages[i])
		}
	}

	return result
}

func (c *messageCache) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.messages = c.messages[:0]
}

func (c *messageCache) Count() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return len(c.messages)
}
