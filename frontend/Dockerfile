# Multi-stage build for Angular app
FROM node:20-alpine AS builder

# Set working directory
WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies - use npm install instead of npm ci for better error handling
RUN npm config set strict-ssl false && npm install

# Copy source code
COPY . .

# Build the Angular app for production
RUN npm run build

# Production stage - Use nginx to serve the built app
FROM nginx:alpine

# Copy built app from builder stage
COPY --from=builder /app/dist/frontend/browser /usr/share/nginx/html

# Copy nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Expose port 80
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"]