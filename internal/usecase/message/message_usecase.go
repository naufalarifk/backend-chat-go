package message

import (
	"context"
	"fmt"
	"realtime-chat-backend/internal/domain/entity"
	"realtime-chat-backend/internal/domain/repository"
	"time"

	"github.com/google/uuid"
)

// UseCase handles message-related business logic
type UseCase struct {
	repo  repository.MessageRepository
	cache repository.MessageCache
}

// NewUseCase creates a new message use case
func NewUseCase(repo repository.MessageRepository, cache repository.MessageCache) *UseCase {
	return &UseCase{
		repo:  repo,
		cache: cache,
	}
}

// CreateMessage creates a new message
func (uc *UseCase) CreateMessage(ctx context.Context, sender, content, roomID string) (*entity.Message, error) {
	if sender == "" {
		return nil, fmt.Errorf("sender is required")
	}
	if content == "" {
		return nil, fmt.Errorf("content is required")
	}
	if roomID == "" {
		return nil, fmt.Errorf("roomID is required")
	}

	message := &entity.Message{
		ID:        uuid.New().String(),
		Sender:    sender,
		Content:   content,
		RoomID:    roomID,
		Timestamp: time.Now(),
	}

	// Save to database
	if err := uc.repo.Create(ctx, message); err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	// Add to cache
	uc.cache.Add(message)

	return message, nil
}

// GetRecentMessages retrieves recent messages
func (uc *UseCase) GetRecentMessages(ctx context.Context, limit int) ([]*entity.Message, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	messages, err := uc.repo.GetRecent(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent messages: %w", err)
	}

	return messages, nil
}

// GetMessagesByRoom retrieves messages for a specific room
func (uc *UseCase) GetMessagesByRoom(ctx context.Context, roomID string, limit int) ([]*entity.Message, error) {
	if roomID == "" {
		return nil, fmt.Errorf("roomID is required")
	}

	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	messages, err := uc.repo.GetByRoomID(ctx, roomID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages for room: %w", err)
	}

	return messages, nil
}

// GetCachedRecentMessages retrieves recent messages from cache
func (uc *UseCase) GetCachedRecentMessages(limit int) []*entity.Message {
	if limit <= 0 {
		limit = 50
	}
	return uc.cache.GetRecent(limit)
}
