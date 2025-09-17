# Golang Fiber PostgreSQL App

A simple web application built with Go Fiber framework and PostgreSQL database connection.

## Features

- RESTful API with CRUD operations
- PostgreSQL database integration
- Environment-based configuration
- CORS and logging middleware
- Health check endpoint
- User management system

## Prerequisites

- Go 1.24.7 or higher
- PostgreSQL database

## Installation

1. Clone the repository:
```bash
git clone https://github.com/ValentinBancel/markdown_converters.git
cd markdown_converters
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your database configuration
```

4. Create a PostgreSQL database:
```sql
CREATE DATABASE fiber_app;
```

## Configuration

Copy `.env.example` to `.env` and configure the following variables:

```env
# Server Configuration
PORT=3000

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=fiber_app
DB_SSLMODE=disable
```

## Running the Application

### Option 1: Demo Version (In-Memory Storage)
Perfect for testing without setting up PostgreSQL:
```bash
make run-demo
# or
go run demo.go
```

### Option 2: PostgreSQL Version
```bash
# First, set up your database configuration
cp .env.example .env
# Edit .env with your database settings

# Run the application
make run-postgres
# or
go run main.go
```

### Option 3: Using Docker Compose
The easiest way to run with PostgreSQL:
```bash
docker-compose up
```

The server will start on `http://localhost:3000` (or the port specified in your .env file).

## API Endpoints

### Health Check
- `GET /health` - Check if the server is running

### Users API
- `GET /api/v1/users` - Get all users
- `POST /api/v1/users` - Create a new user
- `GET /api/v1/users/:id` - Get a user by ID
- `PUT /api/v1/users/:id` - Update a user by ID
- `DELETE /api/v1/users/:id` - Delete a user by ID

### Example Usage

#### Create a user:
```bash
curl -X POST http://localhost:3000/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com"}'
```

#### Get all users:
```bash
curl http://localhost:3000/api/v1/users
```

#### Get a specific user:
```bash
curl http://localhost:3000/api/v1/users/1
```

#### Update a user:
```bash
curl -X PUT http://localhost:3000/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Jane Doe", "email": "jane@example.com"}'
```

#### Delete a user:
```bash
curl -X DELETE http://localhost:3000/api/v1/users/1
```

## Database Schema

The application automatically creates the following table:

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Project Structure

```
.
├── main.go              # Main application with PostgreSQL
├── demo.go              # Demo version with in-memory storage
├── go.mod               # Go module file
├── go.sum               # Go dependencies
├── Makefile             # Build and development commands
├── Dockerfile           # Docker configuration
├── docker-compose.yml   # Docker Compose for easy setup
├── .env.example         # Environment variables template
├── .gitignore           # Git ignore file
└── README.md            # This file
```

## Development Commands

The project includes a Makefile for common development tasks:

```bash
make build         # Build both applications
make run-demo      # Run demo version
make run-postgres  # Run PostgreSQL version
make clean         # Clean build artifacts
make deps          # Install dependencies
make fmt           # Format code
make help          # Show all available commands
```

## Dependencies

- [Fiber](https://github.com/gofiber/fiber) - Web framework
- [pq](https://github.com/lib/pq) - PostgreSQL driver
- [godotenv](https://github.com/joho/godotenv) - Environment variables loader

## Two Application Versions

This project provides two versions:

1. **main.go** - Full PostgreSQL integration for production use
2. **demo.go** - In-memory storage for quick testing and demonstration

Both versions share the same API endpoints and functionality, making it easy to test the application without setting up a database.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.