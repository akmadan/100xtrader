#!/usr/bin/env node

/**
 * 100xTrader CLI
 * Run the entire 100xTrader application without cloning the repository
 * 
 * Usage:
 *   npx 100xtrader
 *   npx 100xtrader --port 3000
 *   npx 100xtrader --help
 */

const { spawn, exec } = require('child_process');
const fs = require('fs');
const path = require('path');
const os = require('os');

const DEFAULT_FRONTEND_PORT = 3000;
const DEFAULT_BACKEND_PORT = 8080;

// Parse command line arguments
const args = process.argv.slice(2);
const help = args.includes('--help') || args.includes('-h');
const port = args.find(arg => arg.startsWith('--port='))?.split('=')[1] || DEFAULT_FRONTEND_PORT;
const backendPort = args.find(arg => arg.startsWith('--backend-port='))?.split('=')[1] || DEFAULT_BACKEND_PORT;
const useDocker = !args.includes('--no-docker');

if (help) {
  console.log(`
🚀 100xTrader - Trading Journal Platform

Usage:
  npx 100xtrader [options]

Options:
  --port=<port>              Frontend port (default: 3000)
  --backend-port=<port>      Backend port (default: 8080)
  --no-docker               Run locally (requires Go & Node.js)
  --help, -h                Show this help message

Examples:
  npx 100xtrader
  npx 100xtrader --port 3001
  npx 100xtrader --no-docker

For more information, visit: https://github.com/yourusername/100xtrader
  `);
  process.exit(0);
}

console.log('🚀 Starting 100xTrader...\n');

if (useDocker) {
  // Check if Docker is available
  exec('docker --version', (error) => {
    if (error) {
      console.error('❌ Docker is not installed or not running.');
      console.error('Please install Docker from https://www.docker.com/get-started');
      console.error('\nOr run with --no-docker flag (requires Go & Node.js)');
      process.exit(1);
    }

    // Check if docker-compose is available
    exec('docker-compose --version', (error) => {
      if (error) {
        console.error('❌ docker-compose is not installed.');
        console.error('Please install docker-compose or use Docker Desktop');
        process.exit(1);
      }

      // Create docker-compose file in temp directory
      const tempDir = path.join(os.tmpdir(), '100xtrader');
      if (!fs.existsSync(tempDir)) {
        fs.mkdirSync(tempDir, { recursive: true });
      }

      const composeFile = path.join(tempDir, 'docker-compose.yml');
      const composeContent = `version: '3.8'

services:
  backend:
    image: 100xtrader/backend:latest
    container_name: 100xtrader-backend
    ports:
      - "${backendPort}:8080"
    environment:
      - DB_PATH=/app/data/temp_db.sqlite
      - PORT=8080
    volumes:
      - backend-data:/app/data
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  frontend:
    image: 100xtrader/frontend:latest
    container_name: 100xtrader-frontend
    ports:
      - "${port}:3000"
    environment:
      - NEXT_PUBLIC_API_URL=http://localhost:${backendPort}
    depends_on:
      backend:
        condition: service_healthy
    restart: unless-stopped

volumes:
  backend-data:
`;

      fs.writeFileSync(composeFile, composeContent);

      console.log('🐳 Pulling Docker images...');
      const pull = spawn('docker-compose', ['-f', composeFile, 'pull'], {
        stdio: 'inherit',
        cwd: tempDir
      });

      pull.on('close', (code) => {
        if (code !== 0) {
          console.error('❌ Failed to pull Docker images.');
          console.error('Make sure the images are published to Docker Hub.');
          process.exit(1);
        }

        console.log('\n✅ Starting containers...\n');
        const up = spawn('docker-compose', ['-f', composeFile, 'up'], {
          stdio: 'inherit',
          cwd: tempDir
        });

        up.on('close', (code) => {
          console.log('\n👋 100xTrader stopped.');
        });

        // Handle Ctrl+C
        process.on('SIGINT', () => {
          console.log('\n\n🛑 Stopping containers...');
          spawn('docker-compose', ['-f', composeFile, 'down'], {
            stdio: 'inherit',
            cwd: tempDir
          });
          process.exit(0);
        });
      });
    });
  });
} else {
  // Local execution (requires Go & Node.js)
  console.log('📦 Running locally (requires Go & Node.js)...\n');
  console.error('❌ Local execution not yet implemented.');
  console.error('Please use Docker or clone the repository.');
  process.exit(1);
}

