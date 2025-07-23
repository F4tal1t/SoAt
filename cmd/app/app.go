package app

import (
	"SoAt/internals/cache"
	"SoAt/internals/database"
	"SoAt/internals/notifications"
	"SoAt/internals/server"
	"SoAt/models/friendships"
	"SoAt/models/posts"
	"SoAt/models/users"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Setup() {
	// Connect to database
	database.Connect()
	
	// Try to connect to Redis, but continue if it fails
	cache.Connect()

	// Add role column if it doesn't exist
	if err := database.DB.Exec("ALTER TABLE IF EXISTS users ADD COLUMN IF NOT EXISTS role VARCHAR(10) DEFAULT 'user'").Error; err != nil {
		log.Printf("Warning: Failed to add role column: %v", err)
	}

	fmt.Println("Running migrations...")
	err := database.DB.AutoMigrate(&users.Users{}, &friendships.Friendships{}, &posts.Posts{})
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	fmt.Println("Migrations completed successfully")

	// Initialize notification system
	notifications.InitNotificationsSystem()
	
	// Try to hydrate notifications, but continue if it fails
	notifications.Hydrate()

	// Setup and start the server
	server.Setup()
	app := server.New()
	fmt.Println("Server starting on :3000")
	
	// Create channel for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	
	// Start server in a goroutine
	go func() {
		if err := app.Listen(":3000"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	
	// Wait for interrupt signal
	<-c
	fmt.Println("\nGracefully shutting down...")
	
	// Close database connections
	if sqlDB, err := database.DB.DB(); err == nil {
		sqlDB.Close()
	}
	
	// Shutdown server with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	app.ShutdownWithContext(ctx)
}
