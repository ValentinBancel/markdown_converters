package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// In-memory storage for demo purposes
var users []User
var nextID = 1

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using defaults")
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Golang Fiber Demo App",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Initialize sample data
	initSampleData()

	// Routes
	setupRoutes(app)

	// Get port from environment or use default
	port := getEnv("PORT", "3000")

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("📖 Visit http://localhost:%s/health for health check", port)
	log.Printf("👥 Visit http://localhost:%s/api/v1/users for users API", port)
	log.Fatal(app.Listen(":" + port))
}

func initSampleData() {
	users = []User{
		{ID: 1, Name: "John Doe", Email: "john@example.com"},
		{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
	}
	nextID = 3
	log.Println("✅ Sample data initialized")
}

func setupRoutes(app *fiber.App) {
	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Golang Fiber app is running successfully!",
			"app":     "Golang Fiber PostgreSQL Demo",
		})
	})

	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to Golang Fiber PostgreSQL App!",
			"endpoints": fiber.Map{
				"health": "/health",
				"users":  "/api/v1/users",
			},
		})
	})

	// API routes
	api := app.Group("/api/v1")

	// User routes
	api.Get("/users", getUsers)
	api.Post("/users", createUser)
	api.Get("/users/:id", getUser)
	api.Put("/users/:id", updateUser)
	api.Delete("/users/:id", deleteUser)
}

func getUsers(c *fiber.Ctx) error {
	return c.JSON(users)
}

func createUser(c *fiber.Ctx) error {
	var user User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if user.Name == "" || user.Email == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Name and email are required"})
	}

	// Check if email already exists
	for _, existingUser := range users {
		if existingUser.Email == user.Email {
			return c.Status(400).JSON(fiber.Map{"error": "Email already exists"})
		}
	}

	user.ID = nextID
	nextID++
	users = append(users, user)

	return c.Status(201).JSON(user)
}

func getUser(c *fiber.Ctx) error {
	id := c.Params("id")
	
	for _, user := range users {
		if fmt.Sprintf("%d", user.ID) == id {
			return c.JSON(user)
		}
	}

	return c.Status(404).JSON(fiber.Map{"error": "User not found"})
}

func updateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var updatedUser User
	if err := c.BodyParser(&updatedUser); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if updatedUser.Name == "" || updatedUser.Email == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Name and email are required"})
	}

	for i, user := range users {
		if fmt.Sprintf("%d", user.ID) == id {
			// Check if email already exists (exclude current user)
			for j, existingUser := range users {
				if j != i && existingUser.Email == updatedUser.Email {
					return c.Status(400).JSON(fiber.Map{"error": "Email already exists"})
				}
			}
			
			users[i].Name = updatedUser.Name
			users[i].Email = updatedUser.Email
			return c.JSON(users[i])
		}
	}

	return c.Status(404).JSON(fiber.Map{"error": "User not found"})
}

func deleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	for i, user := range users {
		if fmt.Sprintf("%d", user.ID) == id {
			users = append(users[:i], users[i+1:]...)
			return c.JSON(fiber.Map{"message": "User deleted successfully"})
		}
	}

	return c.Status(404).JSON(fiber.Map{"error": "User not found"})
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}