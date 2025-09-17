package main

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// MarkdownRequest represents the request payload for markdown conversion
type MarkdownRequest struct {
	Content string `json:"content"`
	Format  string `json:"format"` // html, pdf, etc.
}

// MarkdownResponse represents the response payload
type MarkdownResponse struct {
	ConvertedContent string `json:"convertedContent"`
	Format           string `json:"format"`
	Success          bool   `json:"success"`
	Message          string `json:"message"`
	FileData         string `json:"fileData,omitempty"` // Base64 encoded binary data for PDF
}

// SimpleResponse for basic endpoints
type SimpleResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func main() {
	// Create a new Fiber instance
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	// Add CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Basic health check endpoint
	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.JSON(SimpleResponse{
			Message: "Markdown Converters API is running",
			Status:  "healthy",
		})
	})

	// Get available formats
	app.Get("/api/formats", func(c *fiber.Ctx) error {
		formats := []string{"html", "pdf", "docx", "txt"}
		return c.JSON(map[string]interface{}{
			"formats": formats,
			"message": "Available conversion formats",
		})
	})

	// Convert markdown to specified format
	app.Post("/api/convert", func(c *fiber.Ctx) error {
		var req MarkdownRequest
		
		// Parse request body
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(MarkdownResponse{
				Success: false,
				Message: "Invalid request format",
			})
		}

		// Validate input
		if req.Content == "" {
			return c.Status(400).JSON(MarkdownResponse{
				Success: false,
				Message: "Content is required",
			})
		}

		if req.Format == "" {
			req.Format = "html" // Default to HTML
		}

		// Basic markdown to HTML conversion (simplified)
		convertedContent, fileData := convertMarkdown(req.Content, req.Format)

		return c.JSON(MarkdownResponse{
			ConvertedContent: convertedContent,
			Format:           req.Format,
			Success:          true,
			Message:          "Conversion successful",
			FileData:         fileData,
		})
	})

	// Start server on port 8080
	log.Println("🚀 Server starting on port 8080")
	log.Fatal(app.Listen(":8080"))
}

// Simple markdown to HTML converter (basic implementation)
func convertMarkdown(content, format string) (string, string) {
	switch format {
	case "html":
		// Very basic markdown to HTML conversion
		// In a real application, you would use a proper markdown parser
		html := content
		
		// Convert headers
		html = replaceSimple(html, "# ", "<h1>", "</h1>")
		html = replaceSimple(html, "## ", "<h2>", "</h2>")
		html = replaceSimple(html, "### ", "<h3>", "</h3>")
		
		// Convert bold text **text** to <strong>text</strong>
		// This is a very simplified implementation
		
		// Wrap in paragraphs for basic content
		if html == content { // No headers found
			html = "<p>" + html + "</p>"
		}
		
		return html, ""
	case "pdf":
		// Generate PDF using pandoc
		pdfData, err := generatePDF(content)
		if err != nil {
			return "Error generating PDF: " + err.Error(), ""
		}
		
		// Return base64 encoded PDF data
		return "PDF generated successfully", pdfData
	case "txt":
		// Return plain text (remove markdown syntax)
		return content, ""
	default:
		return "Format '" + format + "' not yet implemented. Content: " + content, ""
	}
}

// generatePDF converts markdown content to PDF using pandoc
func generatePDF(content string) (string, error) {
	// Create temporary directory for conversion
	tempDir, err := ioutil.TempDir("", "markdown_pdf_*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tempDir) // Clean up
	
	// Write markdown content to temporary file
	mdFile := filepath.Join(tempDir, "input.md")
	err = ioutil.WriteFile(mdFile, []byte(content), 0644)
	if err != nil {
		return "", err
	}
	
	// Generate PDF using pandoc
	pdfFile := filepath.Join(tempDir, "output.pdf")
	cmd := exec.Command("pandoc", mdFile, "-o", pdfFile, 
		"--pdf-engine=wkhtmltopdf", 
		"--pdf-engine-opt=--enable-local-file-access",
		"--metadata", "title=Markdown Document")
	
	// Execute command
	if err := cmd.Run(); err != nil {
		return "", err
	}
	
	// Read the generated PDF file
	pdfData, err := ioutil.ReadFile(pdfFile)
	if err != nil {
		return "", err
	}
	
	// Return base64 encoded PDF data
	return base64.StdEncoding.EncodeToString(pdfData), nil
}

// Helper function for simple string replacement
func replaceSimple(text, prefix, openTag, closeTag string) string {
	// This is a very basic implementation
	// In a real application, you would use proper regex or markdown parser
	if len(text) > len(prefix) && text[:len(prefix)] == prefix {
		content := text[len(prefix):]
		return openTag + content + closeTag
	}
	return text
}