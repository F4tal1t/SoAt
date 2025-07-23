package friendships

import (
	"SoAt/internals/cache"
	"SoAt/internals/dto"
	"SoAt/models/friendships"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Friendship struct {
	Friendship *dto.Friends
}

func New() *Friendship {
	return &Friendship{}
}

func (f *Friendship) Create(ctx context.Context) error {
	m := friendships.New()
	m.UserID = f.Friendship.UserID
	m.FriendID = f.Friendship.FriendID

	if err := m.Create(ctx); err != nil {
		return err
	}

	// Copy generated values back
	f.Friendship.ID = uint64(m.ID)
	f.Friendship.CreatedAt = &m.CreatedAt
	f.Friendship.UpdatedAt = &m.UpdatedAt

	return nil
}

func (f *Friendship) Get(ctx context.Context) error {
	m := friendships.New()
	m.ID = uint(f.Friendship.ID)

	if err := m.Get(ctx); err != nil {
		return err
	}

	// Copy retrieved data back
	f.Friendship.UserID = m.UserID
	f.Friendship.FriendID = m.FriendID
	f.Friendship.CreatedAt = &m.CreatedAt
	f.Friendship.UpdatedAt = &m.UpdatedAt

	return nil
}

func (f *Friendship) GetAllFriends(ctx context.Context, userID uuid.UUID) ([]dto.User, error) {
	var dtoFriends []dto.User
	val, err := cache.Client().Get(ctx, f.Friendship.UserID.String()).Result()
	if val != "" && err == nil {
		if err := json.Unmarshal([]byte(val), &dtoFriends); err != nil {
			return nil, err
		}
		return dtoFriends, nil
	}
	m := friendships.New()
	friends, err := m.GetAllFriends(ctx, userID)
	if err != nil {
		return nil, err
	}
	for _, friend := range friends {
		dtoFriends = append(dtoFriends, dto.User{
			ID:        friend.ID,
			Name:      friend.Name,
			Email:     friend.Email,
			CreatedAt: &friend.CreatedAt,
			UpdatedAt: &friend.UpdatedAt,
		})
	}
	b, _ := json.Marshal(dtoFriends)
	if err := cache.Client().Set(ctx, f.Friendship.UserID.String(), b, 24*time.Hour).Err(); err != nil {
		fmt.Println(err)
	}
	return dtoFriends, nil
}

func (f *Friendship) Delete(ctx context.Context) error {
	m := friendships.New()
	m.ID = uint(f.Friendship.ID)

	if err := m.Delete(ctx); err != nil {
		return err
	}

	return nil
}
