# Markdown Converters
A full-stack application for converting markdown content to various formats, built with Angular frontend and Go Fiber API backend.

![Markdown Converters App](https://github.com/user-attachments/assets/da8784ce-2f63-4612-af5b-143c3fad99d0)

## Features

### Frontend (Angular + Custom CSS)
- 🎨 Modern, responsive UI with gradient background
- 📝 Markdown input textarea with format selection
- 🔄 Real-time API status indicator
- 👀 Live preview of converted content
- ⚡ Multiple output format support (HTML, PDF, DOCX, TXT)
- 🚨 Error handling and loading states
- 📱 Mobile-friendly responsive design

### Backend (Go Fiber API)
- 🚀 Fast RESTful API with CORS support
- 💚 Health check endpoint (`/api/health`)
- 📋 Available formats endpoint (`/api/formats`)
- 🔄 Markdown conversion endpoint (`/api/convert`)
- 🛠️ Extensible conversion system
- 📊 JSON API responses

## Project Structure

```
markdown_converters/
├── frontend/                 # Angular application
│   ├── src/
│   │   ├── app/
│   │   │   ├── services/     # API services
│   │   │   ├── app.ts        # Main component
│   │   │   └── app.html      # UI template
│   │   └── styles.css        # Custom CSS (TailwindCSS-like)
│   ├── package.json
│   └── angular.json
├── backend/                  # Go Fiber API
│   ├── main.go              # API server
│   ├── go.mod
│   └── go.sum
└── README.md
```

## Getting Started

### Prerequisites
- Node.js 16+ and npm
- Go 1.18+
- Git

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/ValentinBancel/markdown_converters.git
   cd markdown_converters
   ```

2. **Setup Backend (Go API):**
   ```bash
   cd backend
   go mod download
   go build -o server main.go
   ```

3. **Setup Frontend (Angular):**
   ```bash
   cd ../frontend
   npm install
   ```

### Running the Application

1. **Start the API server:**
   ```bash
   cd backend
   ./server
   ```
   The API will be available at `http://localhost:8080`

2. **Start the frontend development server:**
   ```bash
   cd frontend
   npm start
   ```
   The web application will be available at `http://localhost:4200`

### API Endpoints

- **GET** `/api/health` - Check API health status
- **GET** `/api/formats` - Get available conversion formats
- **POST** `/api/convert` - Convert markdown content

#### Convert Request Example:
```json
{
  "content": "# Hello World\nThis is **bold** text",
  "format": "html"
}
```

#### Convert Response Example:
```json
{
  "convertedContent": "<h1>Hello World\nThis is **bold** text</h1>",
  "format": "html",
  "success": true,
  "message": "Conversion successful"
}
```

## Usage

1. **Enter Markdown:** Type or paste your markdown content into the input area
2. **Choose Format:** Select the desired output format from the dropdown
3. **Convert & Preview:** Click the convert button to see the formatted output

## Development

### Building for Production

**Frontend:**
```bash
cd frontend
npm run build
```

**Backend:**
```bash
cd backend
go build -o server main.go
```

### Features to Add

- [ ] Implement proper markdown parser (e.g., using a library like Blackfriday)
- [ ] Add PDF generation support
- [ ] Add DOCX generation support
- [ ] File upload/download functionality
- [ ] Syntax highlighting for code blocks
- [ ] Live preview as you type
- [ ] Export/save functionality
- [ ] User preferences and themes

## Technology Stack

- **Frontend:** Angular 20, TypeScript, Custom CSS (TailwindCSS-inspired)
- **Backend:** Go, Fiber (web framework)
- **Development:** Node.js, npm, Go modules

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin feature-name`
5. Create a Pull Request

## License

This project is open source and available under the [MIT License](LICENSE).
A Go API using Fiber framework for converting markdown files to various formats, with PostgreSQL database integration.

## Quick Start with Docker Compose

### Prerequisites
- Docker
- Docker Compose

### Quick Start
```bash
# Validate setup
./validate.sh

# Start all services
make start
# OR manually:
docker compose up --build
```

### Using Make Commands
```bash
make help          # Show all available commands
make validate       # Run validation checks
make build          # Build all services
make up             # Start services
make demo           # Run API demo
make down           # Stop services
make clean          # Clean up everything
```

3. The services will be available at:
   - **API**: http://localhost:8080
   - **Database**: localhost:5432
   - **pgAdmin**: http://localhost:5050 (admin@admin.com / admin)

### API Endpoints

- `GET /api/v1/health` - Health check
- `GET /api/v1/files` - Get all markdown files
- `POST /api/v1/files` - Create a new markdown file
- `GET /api/v1/files/:id` - Get a specific markdown file
- `PUT /api/v1/files/:id` - Update a markdown file
- `DELETE /api/v1/files/:id` - Delete a markdown file
- `POST /api/v1/files/:id/convert-html` - Convert markdown to HTML

### Example Usage

#### Using the Demo Script
```bash
# Make sure services are running first
docker compose up -d

# Run the interactive demo
./demo.sh
```

#### Manual API Testing

1. **Health Check**:
   ```bash
   curl http://localhost:8080/api/v1/health
   ```

2. **Create a markdown file**:
   ```bash
   curl -X POST http://localhost:8080/api/v1/files \
     -H "Content-Type: application/json" \
     -d '{"name": "test.md", "content": "# Hello World\nThis is a test markdown file."}'
   ```

3. **Convert to HTML**:
   ```bash
   curl -X POST http://localhost:8080/api/v1/files/1/convert-html
   ```

### Environment Variables

Copy `.env.example` to `.env` and modify as needed:

```bash
cp .env.example .env
```

### Development

To stop the services:
```bash
docker-compose down
```

To rebuild after code changes:
```bash
docker-compose up --build
```

To view logs:
```bash
docker-compose logs api
docker-compose logs postgres
```

### Database Management

Access pgAdmin at http://localhost:5050:
- Email: admin@admin.com
- Password: admin

Add a new server with:
- Host: postgres
- Port: 5432
- Database: markdown_converters
- Username: postgres
- Password: postgres

## Architecture

- **Go API**: Fiber framework with GORM for database operations
- **PostgreSQL**: Database for storing markdown files and metadata
- **Docker**: Containerized deployment
- **pgAdmin**: Web-based database administration tool

## Features

- REST API for markdown file management
- PostgreSQL database integration
- Basic markdown to HTML conversion
- Docker Compose for easy deployment
- Health check endpoints
- CORS enabled for frontend integration
