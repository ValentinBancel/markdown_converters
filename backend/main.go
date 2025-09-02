package main

import (
        "bytes"
        "errors"
        "fmt"
        "log"
        "os"
        "strings"
        "time"

        wkhtmltopdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
        "github.com/gofiber/fiber/v2"
        "github.com/golang-jwt/jwt/v4"
        "github.com/yuin/goldmark"
        "golang.org/x/crypto/bcrypt"
        "gorm.io/driver/postgres"
        "gorm.io/gorm"
)

// User represents a user in the system
type User struct {
        ID          uint `gorm:"primaryKey"`
        Email       string
        PasswordHash string
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

        jwtSecret := os.Getenv("JWT_SECRET")
        if jwtSecret == "" {
                jwtSecret = "secret"
        }

        app.Post("/auth/register", func(c *fiber.Ctx) error {
                var req struct {
                        Email    string `json:"email"`
                        Password string `json:"password"`
                }
                if err := c.BodyParser(&req); err != nil {
                        return fiber.NewError(fiber.StatusBadRequest, err.Error())
                }
                hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
                if err != nil {
                        return fiber.NewError(fiber.StatusInternalServerError, err.Error())
                }
                user := User{Email: req.Email, PasswordHash: string(hash)}
                if err := db.Create(&user).Error; err != nil {
                        return fiber.NewError(fiber.StatusInternalServerError, err.Error())
                }
                return c.SendStatus(fiber.StatusCreated)
        })

        app.Post("/auth/login", func(c *fiber.Ctx) error {
                var req struct {
                        Email    string `json:"email"`
                        Password string `json:"password"`
                }
                if err := c.BodyParser(&req); err != nil {
                        return fiber.NewError(fiber.StatusBadRequest, err.Error())
                }
                var user User
                if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
                        return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
                }
                if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
                        return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
                }
                token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
                        "sub": user.ID,
                        "exp": time.Now().Add(24 * time.Hour).Unix(),
                })
                t, err := token.SignedString([]byte(jwtSecret))
                if err != nil {
                        return fiber.NewError(fiber.StatusInternalServerError, err.Error())
                }
                return c.JSON(fiber.Map{"token": t})
        })

        app.Post("/convert/html", func(c *fiber.Ctx) error {
                userID, err := getUserID(c, jwtSecret)
                if err != nil {
                        return err
                }
                var req struct {
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
                        UserID:   userID,
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
                userID, err := getUserID(c, jwtSecret)
                if err != nil {
                        return err
                }
                var req struct {
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
                        UserID:   userID,
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
                userID, err := getUserID(c, jwtSecret)
                if err != nil {
                        return err
                }
                var convs []Conversion
                if err := db.Where("user_id = ?", userID).Order("created_at desc").Find(&convs).Error; err != nil {
                        return fiber.NewError(fiber.StatusInternalServerError, err.Error())
                }
                return c.JSON(convs)
        })

        app.Get("/history/:id/download", func(c *fiber.Ctx) error {
                userID, err := getUserID(c, jwtSecret)
                if err != nil {
                        return err
                }
                id := c.Params("id")
                var conv Conversion
                if err := db.Where("id = ? AND user_id = ?", id, userID).First(&conv).Error; err != nil {
                        return fiber.NewError(fiber.StatusNotFound, "record not found")
                }
                if conv.Type == "pdf" && conv.PDFPath != "" {
                        return c.Download(conv.PDFPath)
                }
                return c.SendString(conv.HTML)
        })

        log.Fatal(app.Listen(":8080"))
}

func getUserID(c *fiber.Ctx, secret string) (uint, error) {
        auth := c.Get("Authorization")
        if !strings.HasPrefix(auth, "Bearer ") {
                return 0, fiber.NewError(fiber.StatusUnauthorized, "missing token")
        }
        tokenStr := strings.TrimPrefix(auth, "Bearer ")
        token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
                if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
                        return nil, errors.New("invalid signing method")
                }
                return []byte(secret), nil
        })
        if err != nil || !token.Valid {
                return 0, fiber.NewError(fiber.StatusUnauthorized, "invalid token")
        }
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
                return 0, fiber.NewError(fiber.StatusUnauthorized, "invalid token")
        }
        sub, ok := claims["sub"].(float64)
        if !ok {
                return 0, fiber.NewError(fiber.StatusUnauthorized, "invalid token")
        }
        return uint(sub), nil
}
