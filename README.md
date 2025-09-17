# Markdown Converters

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
