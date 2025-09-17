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