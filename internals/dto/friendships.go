package dto

import (
	"time"

	"github.com/google/uuid"
)

type FriendCreate struct {
	UserID   uuid.UUID `json:"user_id" validate:"required,max=100"`
	FriendID uuid.UUID `json:"friend_id" validate:"required,max=100"`
}

type Friends struct {
	ID        uint64     `json:"id" `
	UserID    uuid.UUID  `json:"user_id"`
	FriendID  uuid.UUID  `json:"friend_id"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
