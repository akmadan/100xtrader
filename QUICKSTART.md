# Quick Start Guide

## 🚀 Fastest Way to Run (Docker)

```bash
# 1. Clone the repo
git clone https://github.com/akmadan/100xtrader.git
cd 100xtrader

# 2. Start everything
docker-compose up -d

# 3. Open in browser
# Frontend: http://localhost:3000
# API Docs: http://localhost:8080/swagger/index.html
```

That's it! 🎉

## 📋 Alternative Methods

### Using Make (if you have Go & Node.js installed)

```bash
make install    # Install dependencies
make dev        # Start both servers
```

### Manual Setup

```bash
# Backend
cd go-core && go mod download && go run cmd/main.go

# Frontend (new terminal)
cd web && npm install && npm run dev
```

## 🐳 Docker Commands Cheat Sheet

```bash
docker-compose up -d          # Start
docker-compose down           # Stop
docker-compose logs -f        # View logs
docker-compose restart        # Restart
docker-compose ps             # Check status
```

## ❓ Troubleshooting

**Port already in use?**
```bash
# Change ports in docker-compose.yml or kill existing process
lsof -i :8080
kill -9 <PID>
```

**Need to reset?**
```bash
docker-compose down -v  # Removes volumes too
docker-compose up -d
```

