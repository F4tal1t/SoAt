package posts

import (
	"SoAt/internals/dto"
	"SoAt/internals/notifications"
	"SoAt/internals/validator"
	"SoAt/models/friendships"
	"SoAt/services/posts"
	"SoAt/services/users"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Add(c *fiber.Ctx) error {
	ctx := c.UserContext()

	var post dto.PostCreate

	id := c.Params("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("incorrect user id")
	}
	if err := c.BodyParser(&post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("Incorrect Input Body")
	}
	if err := validator.Posts(post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	us := users.New()
	us.User = &dto.User{}
	us.User.ID = userID
	if err := us.Get(ctx); err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON("user not found")
		}
		return c.Status(fiber.StatusInternalServerError).JSON("Internal Server Error")
	}
	ps := posts.New()
	ps.Post = &dto.Post{}
	ps.Post.UserID = userID
	ps.Post.Content = post.Content

	if err := ps.Create(ctx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("Failed to create post")
	}
	// Notify each friend with their name in the message
	fs := friendships.New()
	friends, err := fs.GetAllFriends(ctx, us.User.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("Failed to fetch friends for notification")
	}
	for _, friend := range friends {
		msg := fmt.Sprintf("Hello %v, your friend %v has created a new post", friend.Name, us.User.Name)
		notifications.NotifyUsers(ctx, friend.ID, msg)
	}
	return c.Status(fiber.StatusCreated).JSON(ps.Post)
}

func Get(c *fiber.Ctx) error {
	ctx := c.UserContext()

	// Get user ID from URL parameter
	userIdStr := c.Params("id")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("incorrect user id")
	}

	ps := posts.New()
	postsList, err := ps.GetAllByUser(ctx, userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("Failed to fetch posts")
	}

	return c.Status(fiber.StatusOK).JSON(postsList)
}

func Delete(c *fiber.Ctx) error {
	ctx := c.UserContext()

	// Get post ID from URL parameter
	postIdStr := c.Params("post_id")
	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("incorrect post id")
	}

	// First check if post exists
	ps := posts.New()
	ps.Post = &dto.Post{}
	ps.Post.ID = postId

	if err := ps.Get(ctx); err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON("Post not found")
		}
		return c.Status(fiber.StatusInternalServerError).JSON("Internal Server Error")
	}

	// Delete the post
	if err := ps.Delete(ctx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("Failed to delete post")
	}

	return c.Status(fiber.StatusOK).JSON("Post deleted successfully")
}
