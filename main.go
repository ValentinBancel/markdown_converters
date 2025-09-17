package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var db *sql.DB

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database connection
	if err := initDB(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Create tables if they don't exist
	if err := createTables(); err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Golang Fiber PostgreSQL App",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Routes
	setupRoutes(app)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

func initDB() error {
	var err error
	
	// Get database connection string from environment
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "fiber_app")
	sslMode := getEnv("DB_SSLMODE", "disable")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, sslMode)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return err
	}

	log.Println("Connected to PostgreSQL database successfully")
	return nil
}

func createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	log.Println("Database tables created successfully")
	return nil
}

func setupRoutes(app *fiber.App) {
	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is running",
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
	rows, err := db.Query("SELECT id, name, email FROM users ORDER BY id")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch users"})
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to scan user"})
		}
		users = append(users, user)
	}

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

	query := "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id"
	err := db.QueryRow(query, user.Name, user.Email).Scan(&user.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
	}

	return c.Status(201).JSON(user)
}

func getUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user User

	query := "SELECT id, name, email FROM users WHERE id = $1"
	err := db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"error": "User not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch user"})
	}

	return c.JSON(user)
}

func updateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if user.Name == "" || user.Email == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Name and email are required"})
	}

	query := "UPDATE users SET name = $1, email = $2 WHERE id = $3"
	result, err := db.Exec(query, user.Name, user.Email, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update user"})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	user.ID = int(rowsAffected)
	return c.JSON(user)
}

func deleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	query := "DELETE FROM users WHERE id = $1"
	result, err := db.Exec(query, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete user"})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(fiber.Map{"message": "User deleted successfully"})
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}