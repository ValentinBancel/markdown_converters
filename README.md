# Markdown Converters

Convert Markdown documents to HTML and PDF while preserving a history of every build.

## Prerequisites
- **Docker**
- **Docker Compose**
- **Go**
- **Node.js**

## Build and Run
1. Clone the repository and open it:
   ```bash
   git clone <repository-url>
   cd markdown_converters
   ```
2. Start all services:
   ```bash
   docker-compose up
   ```
   This launches:
   - the Go backend
   - a PostgreSQL database
   - the Angular frontend served by Nginx
3. When the containers are running, visit:
   - API: <http://localhost:8080/api>
   - Web app: <http://localhost:8081/>

## Rebuild Markdown to HTML/PDF
1. Send Markdown content to the backend to regenerate artifacts:
   ```bash
   curl -F file=@README.md http://localhost:8080/api/convert
   ```
   The service returns fresh HTML/PDF versions and stores the result in the database.
2. View conversion history through:
   - API: <http://localhost:8080/api/history>
   - Web app: <http://localhost:8081/history>

