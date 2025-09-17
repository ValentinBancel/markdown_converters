#!/bin/bash

# Validation script for the markdown converters Docker setup
echo "=== Markdown Converters Docker Setup Validation ==="

# Function to check if a command succeeded
check_result() {
    if [ $? -eq 0 ]; then
        echo "✅ $1"
    else
        echo "❌ $1"
        exit 1
    fi
}

# 1. Check Docker and Docker Compose are available
echo -e "\n1. Checking Docker environment..."
docker --version > /dev/null 2>&1
check_result "Docker is available"

docker compose version > /dev/null 2>&1
check_result "Docker Compose is available"

# 2. Validate docker-compose.yml
echo -e "\n2. Validating docker-compose.yml..."
docker compose config --quiet
check_result "docker-compose.yml syntax is valid"

# 3. Check Go application compiles
echo -e "\n3. Checking Go application..."
go version > /dev/null 2>&1
check_result "Go is available"

go build -o /tmp/test-build . > /dev/null 2>&1
check_result "Go application compiles successfully"
rm -f /tmp/test-build

# 4. Check required files exist
echo -e "\n4. Checking required files..."
required_files=(
    "docker-compose.yml"
    "Dockerfile" 
    "main.go"
    "go.mod"
    "go.sum"
    "init.sql"
    "README.md"
    ".env.example"
    ".gitignore"
)

for file in "${required_files[@]}"; do
    if [ -f "$file" ]; then
        echo "✅ $file exists"
    else
        echo "❌ $file missing"
        exit 1
    fi
done

# 5. Validate docker-compose services configuration
echo -e "\n5. Validating service configuration..."

# Check if services are properly defined
if docker compose config | grep -q "api:"; then
    echo "✅ API service is defined"
else
    echo "❌ API service missing"
    exit 1
fi

if docker compose config | grep -q "postgres:"; then
    echo "✅ PostgreSQL service is defined"
else
    echo "❌ PostgreSQL service missing"
    exit 1
fi

if docker compose config | grep -q "pgadmin:"; then
    echo "✅ pgAdmin service is defined"
else
    echo "❌ pgAdmin service missing"
    exit 1
fi

# 6. Check environment variables
echo -e "\n6. Checking environment configuration..."
if docker compose config | grep -q "DB_HOST"; then
    echo "✅ Database environment variables configured"
else
    echo "❌ Database environment variables missing"
    exit 1
fi

# 7. Check network configuration
echo -e "\n7. Checking network configuration..."
if docker compose config | grep -q "markdown_network"; then
    echo "✅ Custom network is configured"
else
    echo "❌ Custom network missing"
    exit 1
fi

# 8. Check volume configuration
echo -e "\n8. Checking volume configuration..."
if docker compose config | grep -q "postgres_data"; then
    echo "✅ PostgreSQL data volume is configured"
else
    echo "❌ PostgreSQL data volume missing"
    exit 1
fi

# 9. Check health check configuration
echo -e "\n9. Checking health check configuration..."
if docker compose config | grep -q "healthcheck"; then
    echo "✅ Health check is configured for PostgreSQL"
else
    echo "❌ Health check missing"
    exit 1
fi

# 10. Test API endpoints structure (basic syntax check)
echo -e "\n10. Checking API endpoints in code..."
if grep -q "/health" main.go; then
    echo "✅ Health check endpoint defined"
else
    echo "❌ Health check endpoint missing"
    exit 1
fi

if grep -q "/files" main.go; then
    echo "✅ File creation endpoint defined"
else
    echo "❌ File creation endpoint missing"
    exit 1
fi

if grep -q "convert-html" main.go; then
    echo "✅ HTML conversion endpoint defined"
else
    echo "❌ HTML conversion endpoint missing"
    exit 1
fi

echo -e "\n=== ✅ All validations passed! ==="
echo -e "\nTo start the services:"
echo "docker compose up --build"
echo ""
echo "Services will be available at:"
echo "- API: http://localhost:8080"
echo "- PostgreSQL: localhost:5432"
echo "- pgAdmin: http://localhost:5050"
echo ""
echo "Test the API:"
echo "curl http://localhost:8080/api/v1/health"