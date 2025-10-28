package message_test

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"realtime-chat-backend/internal/domain/entity"
	"realtime-chat-backend/internal/usecase/message"
)

type repoMock struct {
	mu        sync.Mutex
	created   []*entity.Message
	recent    []*entity.Message
	byRoom    map[string][]*entity.Message
	createErr error
	getErr    error
	deleteErr error
}

func (r *repoMock) Create(ctx context.Context, m *entity.Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.created = append(r.created, m)
	return r.createErr
}

func (r *repoMock) GetRecent(ctx context.Context, limit int) ([]*entity.Message, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.getErr != nil {
		return nil, r.getErr
	}
	// return up to limit of r.recent
	if limit <= 0 || limit >= len(r.recent) {
		out := make([]*entity.Message, len(r.recent))
		copy(out, r.recent)
		return out, nil
	}
	out := make([]*entity.Message, limit)
	copy(out, r.recent[:limit])
	return out, nil
}

func (r *repoMock) GetByID(ctx context.Context, id string) (*entity.Message, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, m := range r.created {
		if m.ID == id {
			return m, nil
		}
	}
	return nil, fmt.Errorf("message not found")
}

func (r *repoMock) GetByRoomID(ctx context.Context, roomID string, limit int) ([]*entity.Message, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.getErr != nil {
		return nil, r.getErr
	}
	list := r.byRoom[roomID]
	if limit <= 0 || limit >= len(list) {
		out := make([]*entity.Message, len(list))
		copy(out, list)
		return out, nil
	}
	out := make([]*entity.Message, limit)
	copy(out, list[:limit])
	return out, nil
}

func (r *repoMock) Created() []*entity.Message {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]*entity.Message, len(r.created))
	copy(out, r.created)
	return out
}

func (r *repoMock) SetRecent(messages []*entity.Message) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.recent = make([]*entity.Message, len(messages))
	copy(r.recent, messages)
}

func (r *repoMock) SetByRoom(roomID string, messages []*entity.Message) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.byRoom == nil {
		r.byRoom = make(map[string][]*entity.Message)
	}
	r.byRoom[roomID] = make([]*entity.Message, len(messages))
	copy(r.byRoom[roomID], messages)
}

func (r *repoMock) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.deleteErr != nil {
		return r.deleteErr
	}

	filter := func(src []*entity.Message) []*entity.Message {
		if len(src) == 0 {
			return src
		}
		out := src[:0]
		for _, m := range src {
			if m != nil && m.ID != id {
				out = append(out, m)
			}
		}
		return out
	}

	r.created = filter(r.created)
	r.recent = filter(r.recent)

	// Remove from each room list
	for room, list := range r.byRoom {
		r.byRoom[room] = filter(list)
	}

	return nil
}

// cacheMock implements the methods used by the UseCase's cache interface.
type cacheMock struct {
	mu     sync.Mutex
	added  []*entity.Message
	recent []*entity.Message
}

func (c *cacheMock) Add(m *entity.Message) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.added = append(c.added, m)
}

func (c *cacheMock) GetRecent(limit int) []*entity.Message {
	c.mu.Lock()
	defer c.mu.Unlock()
	if limit <= 0 || limit >= len(c.recent) {
		out := make([]*entity.Message, len(c.recent))
		copy(out, c.recent)
		return out
	}
	out := make([]*entity.Message, limit)
	copy(out, c.recent[:limit])
	return out
}

func (c *cacheMock) Added() []*entity.Message {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]*entity.Message, len(c.added))
	copy(out, c.added)
	return out
}

func (c *cacheMock) SetRecent(messages []*entity.Message) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.recent = make([]*entity.Message, len(messages))
	copy(c.recent, messages)
}

func (c *cacheMock) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.added = nil
	c.recent = nil
}

func (c *cacheMock) Count() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.recent)
}

func (c *cacheMock) GetByRoomID(roomID string, limit int) []*entity.Message {
	c.mu.Lock()
	defer c.mu.Unlock()
	var result []*entity.Message
	for _, m := range c.recent {
		if m.RoomID == roomID {
			result = append(result, m)
			if len(result) >= limit {
				break
			}
		}
	}
	return result
}

func TestCreateMessage(t *testing.T) {
	ctx := context.Background()
	rm := &repoMock{}
	cm := &cacheMock{}
	uc := message.NewUseCase(rm, cm)

	message, err := uc.CreateMessage(ctx, "user1", "Hello, World!", "room1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Add to cache
	if message.Sender != "user1" || message.Content != "Hello, World!" || message.RoomID != "room1" {
		t.Errorf("message fields not set correctly: %+v", message)
	}
	fmt.Println(message)
	cm.Add(message)

}
