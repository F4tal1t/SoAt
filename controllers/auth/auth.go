package auth

import (
	"SoAt/internals/auth"
	"SoAt/internals/dto"
	"SoAt/internals/notifications"
	"SoAt/internals/validator"
	"SoAt/models/users"
	usersService "SoAt/services/users"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register handles user registration
func Register(c *fiber.Ctx) error {
	ctx := c.UserContext()

	var user dto.UserCreate

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Incorrect input body",
		})
	}

	if err := validator.Users(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	us := usersService.New()
	us.User = &dto.User{}
	us.User.Name = user.Name
	us.User.Email = user.Email
	us.User.Password = hashedPassword
	us.User.Role = "user" // Default role is 'user'

	us.Create(ctx)

	// Generate JWT token
	token, err := auth.GenerateToken(us.User.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	notifications.Register(us.User.ID)
	go notifications.ListenForNotifications(ctx, us.User.ID)
	go notifications.NotifyUsers(ctx, us.User.ID, "New User Joined")

	// Return user data and token
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user":  us.User,
		"token": token,
	})
}

// Login handles user login
func Login(c *fiber.Ctx) error {
	ctx := c.UserContext()

	var loginRequest LoginRequest
	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Incorrect input body",
		})
	}

	// Find user by email
	user := users.New()
	usersList, err := user.GetAll(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}

	var foundUser *users.Users
	for _, u := range usersList {
		if u.Email == loginRequest.Email {
			foundUser = &u
			break
		}
	}

	if foundUser == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	// Check password
	if !auth.CheckPasswordHash(loginRequest.Password, foundUser.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	// Generate JWT token
	fmt.Printf("Generating token for user ID: %v\n", foundUser.ID)
	token, err := auth.GenerateToken(foundUser.ID)
	if err != nil {
		fmt.Printf("Error generating token: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}
	fmt.Printf("Generated token: %s\n", token)

	// Return user data and token
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": fiber.Map{
			"id":    foundUser.ID,
			"name":  foundUser.Name,
			"email": foundUser.Email,
			"role":  foundUser.Role,
		},
		"token": token,
	})
}