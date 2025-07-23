package notifications

import (
	"SoAt/models/friendships"
	"SoAt/models/users"
	"context"
	"fmt"

	"sync"

	"github.com/google/uuid"
)

// store
var Store map[uuid.UUID]chan string

// mutex
var mu sync.Mutex

func InitNotificationsSystem() {
	Store = make(map[uuid.UUID]chan string)
}
func Register(userID uuid.UUID) {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := Store[userID]; !ok {
		Store[userID] = make(chan string)
	}
}

func ListenForNotifications(ctx context.Context, userID uuid.UUID) {
	mu.Lock()
	channel, ok := Store[userID]
	mu.Unlock()
	if !ok {
		fmt.Printf("No notification channel registered for user %v", userID)
		return
	}
	for {
		select {
		case message := <-channel:
			fmt.Printf("Received a new notification %v\n", message)

		case <-ctx.Done():
			fmt.Printf("Stopping Notification channel for user %v", userID)
			return
		}
	}
}

func NotifyUsers(ctx context.Context, userID uuid.UUID, msg string) {
	// get all friends and notify
	fs := friendships.New()
	friends, err := fs.GetAllFriends(ctx, userID)
	if err != nil {
		fmt.Printf("Error getting friends for notifications: %v\n", err)
		return
	}
	mu.Lock()
	defer mu.Unlock()

	for _, f := range friends {
		if ch, ok := Store[f.ID]; ok {
			ch <- msg
		}
	}
}

func Hydrate() {
	ctx := context.Background()
	us := users.New()
	usersList, err := us.GetAll(ctx)
	if err != nil {
		fmt.Printf("Error hydrating notifications: %v\n", err)
		return
	}
	for _, user := range usersList {
		Register(user.ID)
		go ListenForNotifications(ctx, user.ID)
	}
}
