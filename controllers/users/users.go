package users

import (
	"SoAt/internals/dto"
	"SoAt/internals/notifications"
	"SoAt/internals/validator"
	"SoAt/services/users"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Add(c *fiber.Ctx) error {

	ctx := c.UserContext()

	var user dto.UserCreate

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("incorrect input body")
	}
	if err := validator.Users(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	us := users.New()
	us.User = &dto.User{}
	us.User.Name = user.Name
	us.User.Email = user.Email
	us.User.Password = user.Password

	us.Create(ctx)

	notifications.Register(us.User.ID)
	go notifications.ListenForNotifications(ctx, us.User.ID)
	go notifications.NotifyUsers(ctx, us.User.ID, "New User Joined")
	return c.Status(fiber.StatusCreated).JSON(us.User)
}
func Get(c *fiber.Ctx) error {
	ctx := c.UserContext()
	id := c.Params("id")
	userId, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("incorrect user id")
	}
	us := users.New()
	us.User = &dto.User{}

	us.User.ID = userId
	if err := us.Get(ctx); err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON("User not found")
		}
		return c.Status(fiber.StatusInternalServerError).JSON("Internal Server Error")
	}
	return c.Status(fiber.StatusOK).JSON(us.User)
}

func GetAll(c *fiber.Ctx) error {
	ctx := c.UserContext()
	us := users.New()

	usersList, err := us.GetAll(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("Failed to fetch users")
	}

	return c.Status(fiber.StatusOK).JSON(usersList)
}

func Delete(c *fiber.Ctx) error {
	ctx := c.UserContext()
	id := c.Params("id")
	userId, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("incorrect user id")
	}

	// First check if user exists
	us := users.New()
	us.User = &dto.User{}
	us.User.ID = userId

	if err := us.Get(ctx); err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON("User not found")
		}
		return c.Status(fiber.StatusInternalServerError).JSON("Internal Server Error")
	}

	// Delete the user
	if err := us.Delete(ctx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("Failed to delete user")
	}

	return c.Status(fiber.StatusOK).JSON("User deleted successfully")
}
