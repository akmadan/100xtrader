#!/bin/bash

echo "🚀 Starting 100xtrader API Server"
echo "=================================="

# Navigate to go-core directory
cd go-core

# Download dependencies
echo "📦 Downloading dependencies..."
go mod download

# Generate swagger documentation
echo "📚 Generating swagger documentation..."
$(go env GOPATH)/bin/swag init -g cmd/main.go

# Start API server
echo "🌐 Starting API server on http://localhost:8080"
echo "📖 Swagger UI available at http://localhost:8080/swagger/index.html"
echo "=================================="
go run cmd/main.go
