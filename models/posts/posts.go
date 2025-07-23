package posts

import (
	"SoAt/internals/database"
	"SoAt/models/users"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Posts struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Content   string    `json:"content"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User users.Users `gorm:"foreignKey:UserID;references:ID" json:"-"`
}

func (Posts) TableName() string {
	return "posts"
}

func New() *Posts {
	return &Posts{}
}

func (p *Posts) Create(ctx context.Context) error {
	// Generate UUID if not already set
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	
	if err := database.Client().Create(&p).Error; err != nil {
		fmt.Printf("Unable to create post: %v", err)
		return err
	}
	return nil
}

func (p *Posts) Get(ctx context.Context) error {
	if err := database.Client().Where("id = ?", p.ID).First(&p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Printf("Post not found: %v", err)
			return err
		}
		fmt.Printf("Error getting post: %v", err)
		return err
	}
	return nil
}

func (p *Posts) GetAllByUser(ctx context.Context, userID uuid.UUID) ([]Posts, error) {
	var posts []Posts
	if err := database.Client().Where("user_id = ?", userID).Find(&posts).Error; err != nil {
		fmt.Printf("Error getting posts: %v", err)
		return nil, err
	}
	return posts, nil
}

func (p *Posts) Delete(ctx context.Context) error {
	if err := database.Client().Where("id = ?", p.ID).Delete(&p).Error; err != nil {
		fmt.Printf("Error deleting post: %v", err)
		return err
	}
	return nil
}
