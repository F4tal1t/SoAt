package friendships

import (
	"SoAt/internals/cache"
	"SoAt/internals/dto"
	"SoAt/internals/validator"
	"SoAt/services/friendships"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Add(c *fiber.Ctx) error {
	ctx := c.UserContext()

	var friend dto.FriendCreate

	if err := c.BodyParser(&friend); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("incorrect input body")
	}
	if err := validator.Friendships(friend); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	fs := friendships.New()
	fs.Friendship = &dto.Friends{}
	fs.Friendship.UserID = friend.UserID
	fs.Friendship.FriendID = friend.FriendID

	if err := fs.Create(ctx); err != nil {
		if err.Error() == "friendship already exists" {
			return c.Status(fiber.StatusConflict).JSON("Friendship already exists")
		}
		return c.Status(fiber.StatusInternalServerError).JSON("Failed to create friendship")
	}

	if err := cache.Client().Del(ctx, friend.UserID.String()).Err; err != nil {
		fmt.Printf("Error Invalidating Cache: %v", err)
	}

	return c.Status(fiber.StatusCreated).JSON(fs.Friendship)
}

func Get(c *fiber.Ctx) error {
	ctx := c.UserContext()
	id := c.Params("id")
	friendshipId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("incorrect friendship id")
	}

	fs := friendships.New()
	fs.Friendship = &dto.Friends{}
	fs.Friendship.ID = friendshipId

	if err := fs.Get(ctx); err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON("Friendship not found")
		}
		return c.Status(fiber.StatusInternalServerError).JSON("Internal Server Error")
	}

	return c.Status(fiber.StatusOK).JSON(fs.Friendship)
}

func GetAllFriends(c *fiber.Ctx) error {
	ctx := c.UserContext()
	id := c.Params("userId")
	userId, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("incorrect user id")
	}

	fs := friendships.New()
	friends, err := fs.GetAllFriends(ctx, userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("Failed to fetch friends")
	}

	return c.Status(fiber.StatusOK).JSON(friends)
}

func Delete(c *fiber.Ctx) error {
	ctx := c.UserContext()
	id := c.Params("id")
	friendshipId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("incorrect friendship id")
	}

	// First check if friendship exists
	fs := friendships.New()
	fs.Friendship = &dto.Friends{}
	fs.Friendship.ID = friendshipId

	if err := fs.Get(ctx); err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON("Friendship not found")
		}
		return c.Status(fiber.StatusInternalServerError).JSON("Internal Server Error")
	}

	// Delete the friendship
	if err := fs.Delete(ctx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("Failed to delete friendship")
	}

	if err := cache.Client().Del(ctx, fs.Friendship.UserID.String()).Err(); err != nil {
		fmt.Printf("Error Invalidating Cache: %v", err)
	}

	return c.Status(fiber.StatusOK).JSON("Friendship deleted successfully")
}
