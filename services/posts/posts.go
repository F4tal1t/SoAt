package posts

import (
	"SoAt/internals/dto"
	"SoAt/models/posts"
	"context"
	"github.com/google/uuid"
)

type Post struct {
	Post *dto.Post
}

func New() *Post {
	return &Post{}
}

func (p *Post) Create(ctx context.Context) error {
	m := posts.New()
	
	m.Content = p.Post.Content
	m.UserID = p.Post.UserID
	
	if err := m.Create(ctx); err != nil {
		return err
	}
	
	// Copy generated values back to the DTO
	p.Post.ID = m.ID
	p.Post.CreatedAt = &m.CreatedAt
	p.Post.UpdatedAt = &m.UpdatedAt
	
	return nil
}

func (p *Post) Get(ctx context.Context) error {
	m := posts.New()
	m.ID = p.Post.ID
	
	if err := m.Get(ctx); err != nil {
		return err
	}
	
	// Copy the retrieved data back to the service layer
	p.Post.Content = m.Content
	p.Post.UserID = m.UserID
	p.Post.CreatedAt = &m.CreatedAt
	p.Post.UpdatedAt = &m.UpdatedAt
	
	return nil
}

func (p *Post) GetAllByUser(ctx context.Context, userID uuid.UUID) ([]dto.Post, error) {
	m := posts.New()
	postsList, err := m.GetAllByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	// Convert model posts to DTOs
	var dtoPosts []dto.Post
	for _, post := range postsList {
		dtoPosts = append(dtoPosts, dto.Post{
			ID:        post.ID,
			Content:   post.Content,
			UserID:    post.UserID,
			CreatedAt: &post.CreatedAt,
			UpdatedAt: &post.UpdatedAt,
		})
	}
	
	return dtoPosts, nil
}

func (p *Post) Delete(ctx context.Context) error {
	m := posts.New()
	m.ID = p.Post.ID
	
	if err := m.Delete(ctx); err != nil {
		return err
	}
	
	return nil
}
