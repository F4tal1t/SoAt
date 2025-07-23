package friendships

import (
	"SoAt/internals/database"
	"SoAt/internals/dto"
	"SoAt/models/users"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Friendships struct {
	gorm.Model
	UserID   uuid.UUID   `gorm:"uniqueIndex:idx_user_friend" json:"user_id"`
	FriendID uuid.UUID   `gorm:"uniqueIndex:idx_user_friend" json:"friend_id"`
	User     users.Users `gorm:"foreignKey:UserID;references:ID" json:"-"`
	Friend   users.Users `gorm:"foreignKey:FriendID;references:ID" json:"-"`

	Friends *dto.Friends `gorm:"-"`
}

func New() *Friendships {
	return &Friendships{}
}

func (f *Friendships) Create(ctx context.Context) error {
	// Check if friendship already exists
	var existing Friendships
	if err := database.Client().Where(
		"(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)",
		f.UserID, f.FriendID, f.FriendID, f.UserID,
	).First(&existing).Error; err == nil {
		return errors.New("friendship already exists")
	}

	if err := database.Client().Create(&f).Error; err != nil {
		fmt.Printf("Unable to create friendship: %v", err)
		return err
	}
	return nil
}

func (f *Friendships) Get(ctx context.Context) error {
	if err := database.Client().Where("id = ?", f.ID).First(&f).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Printf("Friendship not found: %v", err)
			return err
		}
		fmt.Printf("Error getting friendship: %v", err)
		return err
	}
	return nil
}

func (f *Friendships) GetAllFriends(ctx context.Context, userID uuid.UUID) ([]users.Users, error) {
	var AllFriends []users.Users

	// Get all friendships where user is either UserID or FriendID
	var friendships []Friendships
	if err := database.Client().Where(
		"user_id = ? OR friend_id = ?", userID, userID,
	).Preload("User").Preload("Friend").Find(&friendships).Error; err != nil {
		fmt.Printf("Error getting friendships: %v", err)
		return nil, err
	}

	// Extract the friends
	for _, friendship := range friendships {
		if friendship.UserID == userID {
			AllFriends = append(AllFriends, friendship.Friend)
		} else {
			AllFriends = append(AllFriends, friendship.User)
		}
	}

	return AllFriends, nil
}

func (f *Friendships) Delete(ctx context.Context) error {
	if err := database.Client().Where("id = ?", f.ID).Delete(&f).Error; err != nil {
		fmt.Printf("Error deleting friendship: %v", err)
		return err
	}
	return nil
}
