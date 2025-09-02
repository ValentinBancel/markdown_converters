package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	wkhtmltopdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/gofiber/fiber/v2"
	"github.com/yuin/goldmark"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID          uint `gorm:"primaryKey"`
	Email       string
	Conversions []Conversion
}

// Conversion stores markdown conversion metadata
type Conversion struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	Markdown  string
	HTML      string
	PDFPath   string
	Type      string
	CreatedAt time.Time
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=markdown_converters port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(&User{}, &Conversion{}); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}

	app := fiber.New()

	app.Post("/convert/html", func(c *fiber.Ctx) error {
		var req struct {
			UserID   uint   `json:"user_id"`
			Markdown string `json:"markdown"`
		}
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		var buf bytes.Buffer
		if err := goldmark.Convert([]byte(req.Markdown), &buf); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		conv := Conversion{
			UserID:   req.UserID,
			Markdown: req.Markdown,
			HTML:     buf.String(),
			Type:     "html",
		}
		if err := db.Create(&conv).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.SendString(conv.HTML)
	})

	app.Post("/convert/pdf", func(c *fiber.Ctx) error {
		var req struct {
			UserID   uint   `json:"user_id"`
			Markdown string `json:"markdown"`
		}
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		var buf bytes.Buffer
		if err := goldmark.Convert([]byte(req.Markdown), &buf); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		pdfg, err := wkhtmltopdf.NewPDFGenerator()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		pdfg.AddPage(wkhtmltopdf.NewPageReader(strings.NewReader(buf.String())))
		if err := pdfg.Create(); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		fileName := fmt.Sprintf("conversion_%d.pdf", time.Now().UnixNano())
		if err := pdfg.WriteFile(fileName); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		conv := Conversion{
			UserID:   req.UserID,
			Markdown: req.Markdown,
			HTML:     buf.String(),
			PDFPath:  fileName,
			Type:     "pdf",
		}
		if err := db.Create(&conv).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.Download(fileName)
	})

	app.Get("/history", func(c *fiber.Ctx) error {
		var convs []Conversion
		if err := db.Order("created_at desc").Find(&convs).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(convs)
	})

	app.Get("/history/:id/download", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var conv Conversion
		if err := db.First(&conv, id).Error; err != nil {
			return fiber.NewError(fiber.StatusNotFound, "record not found")
		}
		if conv.Type == "pdf" && conv.PDFPath != "" {
			return c.Download(conv.PDFPath)
		}
		return c.SendString(conv.HTML)
	})

	log.Fatal(app.Listen(":8080"))
}
