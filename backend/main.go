package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

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
		formats := []string{"html", "html-file", "pdf", "docx", "txt"}
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

// Simple markdown to HTML converter (enhanced implementation)
func convertMarkdown(content, format string) (string, string) {
	switch format {
	case "html":
		// Enhanced markdown to HTML conversion
		html := convertMarkdownToHTML(content)
		return html, ""
	case "html-file":
		// Generate complete HTML file
		html := convertMarkdownToHTML(content)
		completeHTML := generateCompleteHTML(html)
		return completeHTML, ""
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

// Enhanced markdown to HTML converter
func convertMarkdownToHTML(markdown string) string {
	html := markdown
	
	// Convert headers (h1-h6)
	html = regexp.MustCompile(`(?m)^#{6}\s+(.+)$`).ReplaceAllString(html, "<h6>$1</h6>")
	html = regexp.MustCompile(`(?m)^#{5}\s+(.+)$`).ReplaceAllString(html, "<h5>$1</h5>")
	html = regexp.MustCompile(`(?m)^#{4}\s+(.+)$`).ReplaceAllString(html, "<h4>$1</h4>")
	html = regexp.MustCompile(`(?m)^#{3}\s+(.+)$`).ReplaceAllString(html, "<h3>$1</h3>")
	html = regexp.MustCompile(`(?m)^#{2}\s+(.+)$`).ReplaceAllString(html, "<h2>$1</h2>")
	html = regexp.MustCompile(`(?m)^#{1}\s+(.+)$`).ReplaceAllString(html, "<h1>$1</h1>")
	
	// Convert bold text **text** to <strong>text</strong>
	html = regexp.MustCompile(`\*\*([^*]+)\*\*`).ReplaceAllString(html, "<strong>$1</strong>")
	
	// Convert italic text *text* to <em>text</em>
	html = regexp.MustCompile(`\*([^*]+)\*`).ReplaceAllString(html, "<em>$1</em>")
	
	// Convert inline code `code` to <code>code</code>
	html = regexp.MustCompile("`([^`]+)`").ReplaceAllString(html, "<code>$1</code>")
	
	// Convert links [text](url) to <a href="url">text</a>
	html = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`).ReplaceAllString(html, `<a href="$2">$1</a>`)
	
	// Convert images ![alt](url) to <img src="url" alt="alt">
	html = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`).ReplaceAllString(html, `<img src="$2" alt="$1">`)
	
	// Convert unordered lists
	html = convertUnorderedLists(html)
	
	// Convert ordered lists
	html = convertOrderedLists(html)
	
	// Convert code blocks ```code```
	html = convertCodeBlocks(html)
	
	// Convert line breaks and paragraphs
	html = convertParagraphs(html)
	
	return html
}

// Convert unordered lists (- item or * item)
func convertUnorderedLists(html string) string {
	lines := strings.Split(html, "\n")
	var result []string
	inList := false
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			if !inList {
				result = append(result, "<ul>")
				inList = true
			}
			content := strings.TrimSpace(trimmed[2:])
			result = append(result, "  <li>"+content+"</li>")
		} else {
			if inList {
				result = append(result, "</ul>")
				inList = false
			}
			result = append(result, line)
		}
	}
	
	if inList {
		result = append(result, "</ul>")
	}
	
	return strings.Join(result, "\n")
}

// Convert ordered lists (1. item)
func convertOrderedLists(html string) string {
	lines := strings.Split(html, "\n")
	var result []string
	inList := false
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if regexp.MustCompile(`^\d+\.\s+`).MatchString(trimmed) {
			if !inList {
				result = append(result, "<ol>")
				inList = true
			}
			content := regexp.MustCompile(`^\d+\.\s+`).ReplaceAllString(trimmed, "")
			result = append(result, "  <li>"+content+"</li>")
		} else {
			if inList {
				result = append(result, "</ol>")
				inList = false
			}
			result = append(result, line)
		}
	}
	
	if inList {
		result = append(result, "</ol>")
	}
	
	return strings.Join(result, "\n")
}

// Convert code blocks ```language\ncode\n```
func convertCodeBlocks(html string) string {
	// Convert fenced code blocks
	re := regexp.MustCompile("(?s)```([a-zA-Z]*)\n([\\s\\S]*?)```")
	html = re.ReplaceAllStringFunc(html, func(match string) string {
		parts := re.FindStringSubmatch(match)
		language := parts[1]
		code := strings.TrimSpace(parts[2])
		
		if language != "" {
			return fmt.Sprintf(`<pre><code class="language-%s">%s</code></pre>`, language, escapeHTML(code))
		}
		return fmt.Sprintf(`<pre><code>%s</code></pre>`, escapeHTML(code))
	})
	
	return html
}

// Convert paragraphs (handle line breaks)
func convertParagraphs(html string) string {
	lines := strings.Split(html, "\n")
	var result []string
	var currentParagraph []string
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		// Skip already processed HTML tags
		if strings.HasPrefix(trimmed, "<") && strings.HasSuffix(trimmed, ">") {
			// Finish current paragraph if any
			if len(currentParagraph) > 0 {
				result = append(result, "<p>"+strings.Join(currentParagraph, " ")+"</p>")
				currentParagraph = nil
			}
			result = append(result, line)
			continue
		}
		
		// Empty line ends paragraph
		if trimmed == "" {
			if len(currentParagraph) > 0 {
				result = append(result, "<p>"+strings.Join(currentParagraph, " ")+"</p>")
				currentParagraph = nil
			}
			continue
		}
		
		// Add to current paragraph
		currentParagraph = append(currentParagraph, trimmed)
	}
	
	// Finish last paragraph if any
	if len(currentParagraph) > 0 {
		result = append(result, "<p>"+strings.Join(currentParagraph, " ")+"</p>")
	}
	
	return strings.Join(result, "\n")
}

// Generate complete HTML document
func generateCompleteHTML(bodyContent string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Markdown Document</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            background-color: #fff;
        }
        h1, h2, h3, h4, h5, h6 {
            margin-top: 1.5em;
            margin-bottom: 0.5em;
            font-weight: 600;
        }
        h1 { font-size: 2.25em; border-bottom: 1px solid #eee; padding-bottom: 0.3em; }
        h2 { font-size: 1.75em; }
        h3 { font-size: 1.5em; }
        h4 { font-size: 1.25em; }
        h5 { font-size: 1.1em; }
        h6 { font-size: 1em; }
        p { margin-bottom: 1em; }
        code {
            background-color: #f6f8fa;
            border-radius: 3px;
            font-size: 85%%;
            margin: 0;
            padding: 0.2em 0.4em;
            font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
        }
        pre {
            background-color: #f6f8fa;
            border-radius: 6px;
            overflow: auto;
            padding: 16px;
            margin: 1em 0;
        }
        pre code {
            background-color: transparent;
            border: 0;
            display: inline;
            line-height: inherit;
            margin: 0;
            overflow: visible;
            padding: 0;
            word-wrap: normal;
        }
        ul, ol {
            margin-bottom: 1em;
            padding-left: 2em;
        }
        li {
            margin-bottom: 0.25em;
        }
        a {
            color: #0969da;
            text-decoration: none;
        }
        a:hover {
            text-decoration: underline;
        }
        img {
            max-width: 100%%;
            height: auto;
        }
        strong {
            font-weight: 600;
        }
        em {
            font-style: italic;
        }
    </style>
</head>
<body>
%s
</body>
</html>`, bodyContent)
}

// Escape HTML special characters
func escapeHTML(text string) string {
	text = strings.ReplaceAll(text, "&", "&amp;")
	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")
	text = strings.ReplaceAll(text, "\"", "&quot;")
	text = strings.ReplaceAll(text, "'", "&#39;")
	return text
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