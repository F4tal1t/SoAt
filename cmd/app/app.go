package app

import (
	"SoAt/internals/cache"
	"SoAt/internals/database"
	"SoAt/internals/notifications"
	"SoAt/internals/server"
	"SoAt/models/friendships"
	"SoAt/models/posts"
	"SoAt/models/users"
	"fmt"
	"log"
)

func Setup() {
	database.Connect()
	cache.Connect()

	fmt.Println("Running migrations...")
	err := database.DB.AutoMigrate(&users.Users{}, &friendships.Friendships{}, &posts.Posts{})
	if err != nil {
		panic("Failed to run migrations")
	}
	fmt.Println("Migrations completed successfully")

	notifications.InitNotificationsSystem()
	notifications.Hydrate()

	server.Setup()
	app := server.New()
	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
