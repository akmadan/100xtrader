# Testing Locally with Docker

## ✅ Quick Start (Production Build)

```bash
# 1. Build and start all services
docker-compose up -d --build

# 2. Wait for services to start (15-20 seconds)
sleep 15

# 3. Check container status
docker-compose ps

# 4. Test the backend health endpoint
curl http://localhost:8080/health

# 5. Test the frontend (open in browser)
open http://localhost:3000

# 6. View logs
docker-compose logs -f

# 7. Stop everything
docker-compose down
```

## 🚀 One-Line Test Command

```bash
docker-compose up -d --build && sleep 15 && docker-compose ps && curl http://localhost:8080/health && echo "" && echo "✅ Backend: http://localhost:8080" && echo "✅ Frontend: http://localhost:3000"
```

## Development Mode (Hot Reload)

```bash
# 1. Build and start in development mode
docker-compose -f docker-compose.dev.yml up --build

# 2. Check logs
docker-compose -f docker-compose.dev.yml logs -f

# 3. Stop
docker-compose -f docker-compose.dev.yml down
```

## Step-by-Step Commands

### Build Images Only
```bash
# Build backend
docker build -f Dockerfile.backend -t 100xtrader-backend:local .

# Build frontend
docker build -f Dockerfile.frontend -t 100xtrader-frontend:local .
```

### Run with Docker Compose
```bash
# Start services (detached mode)
docker-compose up -d --build

# View logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f backend
docker-compose logs -f frontend

# Check service status
docker-compose ps

# Stop services
docker-compose down

# Stop and remove volumes (clean slate)
docker-compose down -v
```

### Test Endpoints
```bash
# Health check
curl http://localhost:8080/health

# API docs
open http://localhost:8080/swagger/index.html

# Frontend
open http://localhost:3000
```

### Debugging
```bash
# Enter backend container
docker exec -it 100xtrader-backend sh

# Enter frontend container
docker exec -it 100xtrader-frontend sh

# View backend database
docker exec -it 100xtrader-backend sqlite3 /app/data/db.sqlite ".tables"
```

## Clean Up
```bash
# Remove all containers, networks, and volumes
docker-compose down -v

# Remove images
docker rmi 100xtrader-backend 100xtrader-frontend

# Full cleanup (removes everything)
docker system prune -a --volumes
```

