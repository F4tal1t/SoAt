package users

import (
	"SoAt/internals/database"
	"SoAt/internals/dto"
	"context"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Users struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User *dto.User `gorm:"-"`
}

func New() *Users {
	return &Users{}
}

func (u *Users) Create(ctx context.Context) {
	// Generate UUID if not already set
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}

	if err := database.Client().Create(&u).Error; err != nil {
		fmt.Printf("Unable to create user: %v", err)
	}
}

func (u *Users) Get(ctx context.Context) error {
	if err := database.Client().Where("id = ?", u.ID).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Printf("User not found: %v", err)
			return err
		}
		fmt.Printf("Error getting user: %v", err)
		return err
	}
	return nil
}

func (u *Users) GetAll(ctx context.Context) ([]Users, error) {
	var users []Users
	if err := database.Client().Find(&users).Error; err != nil {
		fmt.Printf("Error getting all users: %v", err)
		return nil, err
	}
	return users, nil
}

func (u *Users) Delete(ctx context.Context) error {
	if err := database.Client().Where("id = ?", u.ID).Delete(&u).Error; err != nil {
		fmt.Printf("Error deleting user: %v", err)
		return err
	}
	return nil
}
