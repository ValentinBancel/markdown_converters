package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MarkdownFile struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Content  string `json:"content"`
	HTMLData string `json:"html_data,omitempty"`
}

var db *gorm.DB

func initDatabase() {
	var err error
	
	// Get database configuration from environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "markdown_converters")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&MarkdownFile{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database connected and migrated successfully")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// Initialize database
	initDatabase()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Markdown Converters API",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Routes
	api := app.Group("/api/v1")
	
	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"message": "Markdown Converters API is running",
		})
	})

	// Markdown files routes
	api.Get("/files", getMarkdownFiles)
	api.Post("/files", createMarkdownFile)
	api.Get("/files/:id", getMarkdownFile)
	api.Put("/files/:id", updateMarkdownFile)
	api.Delete("/files/:id", deleteMarkdownFile)
	api.Post("/files/:id/convert-html", convertToHTML)

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

func getMarkdownFiles(c *fiber.Ctx) error {
	var files []MarkdownFile
	result := db.Find(&files)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch files"})
	}
	return c.JSON(files)
}

func createMarkdownFile(c *fiber.Ctx) error {
	var file MarkdownFile
	if err := c.BodyParser(&file); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	result := db.Create(&file)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create file"})
	}

	return c.Status(201).JSON(file)
}

func getMarkdownFile(c *fiber.Ctx) error {
	id := c.Params("id")
	var file MarkdownFile
	
	result := db.First(&file, id)
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "File not found"})
	}

	return c.JSON(file)
}

func updateMarkdownFile(c *fiber.Ctx) error {
	id := c.Params("id")
	var file MarkdownFile
	
	// Check if file exists
	result := db.First(&file, id)
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "File not found"})
	}

	// Parse request body
	if err := c.BodyParser(&file); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Update file
	result = db.Save(&file)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update file"})
	}

	return c.JSON(file)
}

func deleteMarkdownFile(c *fiber.Ctx) error {
	id := c.Params("id")
	
	result := db.Delete(&MarkdownFile{}, id)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete file"})
	}

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "File not found"})
	}

	return c.JSON(fiber.Map{"message": "File deleted successfully"})
}

func convertToHTML(c *fiber.Ctx) error {
	id := c.Params("id")
	var file MarkdownFile
	
	// Get the markdown file
	result := db.First(&file, id)
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "File not found"})
	}

	// Simple markdown to HTML conversion (basic implementation)
	// In a real application, you'd use a proper markdown parser like blackfriday
	htmlContent := simpleMarkdownToHTML(file.Content)
	
	// Update the file with HTML data
	file.HTMLData = htmlContent
	db.Save(&file)

	return c.JSON(fiber.Map{
		"message": "Converted to HTML successfully",
		"html": htmlContent,
	})
}

// Simple markdown to HTML converter (basic implementation)
func simpleMarkdownToHTML(markdown string) string {
	// This is a very basic implementation for demonstration
	// In production, use a proper markdown library like blackfriday
	html := markdown
	
	// Convert headers
	html = replaceSimple(html, "# ", "<h1>", "</h1>")
	html = replaceSimple(html, "## ", "<h2>", "</h2>")
	html = replaceSimple(html, "### ", "<h3>", "</h3>")
	
	// Convert paragraphs (very basic)
	if html != "" && html[0] != '<' {
		html = "<p>" + html + "</p>"
	}
	
	return html
}

func replaceSimple(text, prefix, openTag, closeTag string) string {
	if len(text) > len(prefix) && text[:len(prefix)] == prefix {
		content := text[len(prefix):]
		return openTag + content + closeTag
	}
	return text
}