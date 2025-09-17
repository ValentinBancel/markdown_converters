#!/bin/bash

# Demo script for testing the Markdown Converters API
echo "=== Markdown Converters API Demo ==="

API_BASE_URL="http://localhost:8080/api/v1"

# Function to check if API is running
check_api() {
    echo -n "Checking if API is running... "
    if curl -s "${API_BASE_URL}/health" > /dev/null 2>&1; then
        echo "✅ API is running"
        return 0
    else
        echo "❌ API is not running"
        echo "Please start the services with: docker compose up --build"
        return 1
    fi
}

# Function to make API calls with error handling
api_call() {
    local method=$1
    local endpoint=$2
    local data=$3
    
    echo -e "\n--- $method $endpoint ---"
    
    if [ -n "$data" ]; then
        response=$(curl -s -w "\nHTTP_STATUS_CODE:%{http_code}" \
            -X $method \
            -H "Content-Type: application/json" \
            -d "$data" \
            "${API_BASE_URL}${endpoint}")
    else
        response=$(curl -s -w "\nHTTP_STATUS_CODE:%{http_code}" \
            -X $method \
            "${API_BASE_URL}${endpoint}")
    fi
    
    # Extract body and status code
    body=$(echo "$response" | sed '$d')
    status_code=$(echo "$response" | grep "HTTP_STATUS_CODE" | cut -d: -f2)
    
    echo "Status: $status_code"
    echo "Response: $body"
    
    # Pretty print JSON if it's valid
    if echo "$body" | python3 -m json.tool > /dev/null 2>&1; then
        echo "Formatted:"
        echo "$body" | python3 -m json.tool
    fi
}

# Main demo
if ! check_api; then
    exit 1
fi

echo -e "\n=== Starting API Demo ==="

# 1. Health check
api_call "GET" "/health"

# 2. List files (should be empty initially)
api_call "GET" "/files"

# 3. Create a markdown file
echo -e "\n💾 Creating a markdown file..."
api_call "POST" "/files" '{
    "name": "example.md",
    "content": "# Hello World\n\nThis is a **markdown** file created via API.\n\n## Features\n\n- Easy to use\n- RESTful API\n- PostgreSQL storage"
}'

# 4. Create another file
echo -e "\n💾 Creating another markdown file..."
api_call "POST" "/files" '{
    "name": "readme.md", 
    "content": "# Project README\n\nThis project demonstrates:\n\n1. Go API with Fiber\n2. PostgreSQL integration\n3. Docker containerization"
}'

# 5. List all files
echo -e "\n📋 Listing all files..."
api_call "GET" "/files"

# 6. Get specific file
echo -e "\n📄 Getting file with ID 1..."
api_call "GET" "/files/1"

# 7. Convert to HTML
echo -e "\n🔄 Converting file 1 to HTML..."
api_call "POST" "/files/1/convert-html"

# 8. Get the file again to see HTML content
echo -e "\n📄 Getting file 1 again (should now have HTML data)..."
api_call "GET" "/files/1"

# 9. Update a file
echo -e "\n✏️ Updating file 1..."
api_call "PUT" "/files/1" '{
    "id": 1,
    "name": "updated-example.md",
    "content": "# Updated Hello World\n\nThis file has been **updated** via the API.\n\n## New Section\n\nAdded some new content!"
}'

# 10. Convert updated file to HTML
echo -e "\n🔄 Converting updated file to HTML..."
api_call "POST" "/files/1/convert-html"

echo -e "\n=== Demo completed successfully! ==="
echo -e "\nNext steps:"
echo "- Access pgAdmin at http://localhost:5050 (admin@admin.com / admin)"
echo "- View database tables and data"
echo "- Try the API endpoints manually with curl or Postman"
echo "- Extend the API with more features"