#!/bin/bash

# 100xTrader Setup Script
# This script helps you set up the project quickly

set -e

echo "🚀 100xTrader Setup"
echo "==================="
echo ""

# Check if Docker is available
if command -v docker &> /dev/null && command -v docker-compose &> /dev/null; then
    echo "✅ Docker detected"
    echo ""
    echo "Would you like to use Docker? (Recommended) [Y/n]"
    read -r response
    
    if [[ "$response" =~ ^([yY][eE][sS]|[yY]|"")$ ]]; then
        echo ""
        echo "🐳 Starting with Docker..."
        docker-compose up -d
        echo ""
        echo "✅ Setup complete!"
        echo ""
        echo "Access the application:"
        echo "  Frontend: http://localhost:3000"
        echo "  Backend:  http://localhost:8080"
        echo "  API Docs: http://localhost:8080/swagger/index.html"
        echo ""
        echo "View logs: docker-compose logs -f"
        echo "Stop:      docker-compose down"
        exit 0
    fi
fi

# Local setup
echo "📦 Setting up locally..."
echo ""

# Check prerequisites
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.24+ from https://go.dev/doc/install"
    exit 1
fi

if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed. Please install Node.js 20+ from https://nodejs.org/"
    exit 1
fi

echo "✅ Prerequisites check passed"
echo ""

# Install Go dependencies
echo "📦 Installing Go dependencies..."
cd go-core
go mod download
cd ..

# Install Node.js dependencies
echo "📦 Installing Node.js dependencies..."
cd web
npm install
cd ..

echo ""
echo "✅ Setup complete!"
echo ""
echo "To start development servers:"
echo "  make dev          # Start both servers"
echo "  npm run dev       # Or use npm scripts"
echo ""
echo "Or start manually:"
echo "  Terminal 1: cd go-core && go run cmd/main.go"
echo "  Terminal 2: cd web && npm run dev"
echo ""

