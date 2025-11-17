# 100xTrader

A comprehensive trading journal platform for tracking trades, analyzing performance, and automating trading strategies.

## 🚀 Quick Start

### Option 1: Run Without Cloning (Easiest!)

Run 100xTrader instantly without cloning the repository:

```bash
# Using npx (requires npm)
npx 100xtrader

# Or using Docker directly
curl -o docker-compose.yml https://raw.githubusercontent.com/yourusername/100xtrader/main/docker-compose.public.yml
docker-compose up -d
```

**Access the application:**
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- API Docs (Swagger): http://localhost:8080/swagger/index.html

**Stop:**
```bash
# If using npx, press Ctrl+C
# If using docker-compose
docker-compose down
```

### Option 2: Clone and Run Locally

If you want to contribute or customize:

```bash
# Clone the repository
git clone https://github.com/yourusername/100xtrader.git
cd 100xtrader

# Start everything with Docker Compose
docker-compose up -d

# View logs
docker-compose logs -f
```

**Development mode with hot reload:**
```bash
docker-compose -f docker-compose.dev.yml up
```

### Option 3: Automated Setup Script

Run our interactive setup script:

```bash
# Clone the repository
git clone https://github.com/yourusername/100xtrader.git
cd 100xtrader

# Run setup script
chmod +x setup.sh
./setup.sh
```

The script will:
- Detect Docker and offer Docker setup (recommended)
- Install all dependencies automatically
- Guide you through the process

### Option 4: Using Make Commands

If you have Go and Node.js installed:

```bash
# Install all dependencies
make install

# Start development servers (both backend and frontend)
make dev

# Build for production
make build

# Docker commands
make docker-up      # Start containers
make docker-down    # Stop containers
make docker-logs    # View logs
```

See [Makefile](Makefile) for all available commands.

### Option 5: Using npm Scripts

We also provide npm scripts at the root level:

```bash
# Install all dependencies
npm run install:all

# Start development servers (both backend and frontend)
npm run dev

# Build for production
npm run build

# Docker commands
npm run docker:up
npm run docker:down
npm run docker:logs
```

### Option 6: Manual Setup

#### Prerequisites

- **Go** 1.24+ ([Install](https://go.dev/doc/install))
- **Node.js** 20+ and npm ([Install](https://nodejs.org/))
- **SQLite** (usually pre-installed on macOS/Linux)

#### Setup Steps

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/100xtrader.git
   cd 100xtrader
   ```

2. **Install dependencies**
   ```bash
   # Install Go dependencies
   cd go-core
   go mod download
   
   # Install Node.js dependencies
   cd ../web
   npm install
   ```

3. **Configure environment** (Optional)
   ```bash
   # Copy example env file
   cp .env.example .env
   
   # Edit .env with your settings
   nano .env
   ```

4. **Start development servers**
   ```bash
   # Terminal 1 - Backend
   cd go-core
   go run cmd/main.go
   
   # Terminal 2 - Frontend
   cd web
   npm run dev
   ```

5. **Access the application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - API Docs: http://localhost:8080/swagger/index.html

## 📁 Project Structure

```
100xtrader/
├── go-core/              # Backend (Go)
│   ├── cmd/              # Application entry point
│   ├── internal/         # Internal packages
│   │   ├── api/          # API handlers, DTOs, middleware
│   │   ├── data/         # Database models and repositories
│   │   ├── services/     # Business logic (brokers, etc.)
│   │   └── utils/        # Utility functions
│   └── migrations/       # Database migrations
├── web/                  # Frontend (Next.js/React)
│   ├── src/
│   │   ├── app/          # Next.js app directory
│   │   ├── components/   # React components
│   │   ├── services/     # API service layer
│   │   └── types/        # TypeScript types
│   └── public/           # Static assets
├── scripts/              # Helper scripts
├── docker-compose.yml    # Docker Compose configuration
├── Dockerfile.backend    # Backend Docker image
├── Dockerfile.frontend   # Frontend Docker image
├── Makefile              # Make commands
└── README.md             # This file
```

## 🛠️ Development

### Backend Development

```bash
cd go-core

# Run server
go run cmd/main.go

# Run with hot reload (install air: go install github.com/cosmtrek/air@latest)
air

# Generate Swagger docs
swag init -g cmd/main.go
```

### Frontend Development

```bash
cd web

# Start dev server
npm run dev

# Build for production
npm run build

# Start production server
npm start

# Lint code
npm run lint
```

## 🐳 Docker Details

### Building Images

```bash
# Build all images
docker-compose build

# Build specific service
docker-compose build backend
docker-compose build frontend
```

### Running Containers

```bash
# Start in detached mode
docker-compose up -d

# Start with logs
docker-compose up

# Stop containers
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

### Environment Variables

Create a `.env` file in the root directory:

```env
# Backend
DB_PATH=./db.sqlite
PORT=8080

# Dhan API (Optional)
DHAN_PROD_API=https://api.dhan.co/v2
DHAN_SANDBOX_API=https://sandbox.dhan.co/v2

# Frontend
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## 📚 API Documentation

Once the backend is running, visit:
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **API Base URL**: http://localhost:8080/api/v1

## 🧪 Testing

```bash
# Backend tests
cd go-core
go test ./...

# Frontend tests (when added)
cd web
npm test
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Troubleshooting

### Port Already in Use

```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>

# Or change port in .env
PORT=8081
```

### Database Issues

```bash
# Remove existing database
rm db.sqlite

# Restart server (will recreate database)
make dev
```

### Docker Issues

```bash
# Clean everything and rebuild
make docker-clean
docker-compose build --no-cache
docker-compose up
```

### Node Modules Issues

```bash
cd web
rm -rf node_modules package-lock.json
npm install
```

### Go Dependencies Issues

```bash
cd go-core
go mod tidy
go mod download
```

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/yourusername/100xtrader/issues)
- **Discussions**: [GitHub Discussions](https://github.com/yourusername/100xtrader/discussions)

## 🎯 Features

- ✅ Trade Journal & Tracking
- ✅ Performance Analytics
- ✅ Strategy Management
- ✅ Rule & Mistake Tracking
- ✅ Calendar View
- ✅ Broker Integration (Dhan, Zerodha)
- ✅ Algorithm Builder (Canvas & Code Views)
- ✅ Automated Trading Algorithms

---

**Made with ❤️ for traders**
