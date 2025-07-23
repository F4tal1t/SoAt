package users

import (
	"SoAt/internals/dto"
	"SoAt/models/users"
	"context"
)

type User struct {
	User *dto.User
}

func New() *User {
	return &User{}
}

func (u *User) Create(ctx context.Context) {
	m := users.New()

	m.Name = u.User.Name
	m.Email = u.User.Email
	m.Password = u.User.Password
	m.Create(ctx)
	
	// Copy generated values back to the DTO
	u.User.ID = m.ID
	u.User.CreatedAt = &m.CreatedAt
	u.User.UpdatedAt = &m.UpdatedAt
}

func (u *User) Get(ctx context.Context) error {
	m := users.New()
	m.ID = u.User.ID
	m.User = u.User

	if err := m.Get(ctx); err != nil {
		return err
	}

	// Copy the retrieved data back to the service layer
	u.User.Name = m.Name
	u.User.Email = m.Email
	u.User.CreatedAt = &m.CreatedAt
	u.User.UpdatedAt = &m.UpdatedAt

	return nil
}

func (u *User) GetAll(ctx context.Context) ([]dto.User, error) {
	m := users.New()
	usersList, err := m.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Convert model users to DTOs
	var dtoUsers []dto.User
	for _, user := range usersList {
		dtoUsers = append(dtoUsers, dto.User{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: &user.CreatedAt,
			UpdatedAt: &user.UpdatedAt,
		})
	}

	return dtoUsers, nil
}

func (u *User) Delete(ctx context.Context) error {
	m := users.New()
	m.ID = u.User.ID

	if err := m.Delete(ctx); err != nil {
		return err
	}

	return nil
}
