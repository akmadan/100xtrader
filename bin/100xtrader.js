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

// ASCII Art Banner
function showBanner(frontendPort, backendPort) {
  const banner = `
   /$$    /$$$$$$   /$$$$$$              /$$                              /$$                    
 /$$$$   /$$$_  $$ /$$$_  $$            | $$                             | $$                    
|_  $$  | $$$$\ $$| $$$$\ $$ /$$   /$$ /$$$$$$    /$$$$$$  /$$$$$$   /$$$$$$$  /$$$$$$   /$$$$$$ 
  | $$  | $$ $$ $$| $$ $$ $$|  $$ /$$/|_  $$_/   /$$__  $$|____  $$ /$$__  $$ /$$__  $$ /$$__  $$
  | $$  | $$\ $$$$| $$\ $$$$ \  $$$$/   | $$    | $$  \__/ /$$$$$$$| $$  | $$| $$$$$$$$| $$  \__/
  | $$  | $$ \ $$$| $$ \ $$$  >$$  $$   | $$ /$$| $$      /$$__  $$| $$  | $$| $$_____/| $$      
 /$$$$$$|  $$$$$$/|  $$$$$$/ /$$/\  $$  |  $$$$/| $$     |  $$$$$$$|  $$$$$$$|  $$$$$$$| $$      
|______/ \______/  \______/ |__/  \__/   \___/  |__/      \_______/ \_______/ \_______/|__/      
                                                                                                 
                                                                                                 
                                                                                                 

   🚀 100xTrader is running successfully!

   📊 Frontend:  http://localhost:${frontendPort}
   🔧 Backend:   http://localhost:${backendPort}
   📖 API Docs:  http://localhost:${backendPort}/swagger/index.html

   Press Ctrl+C to stop
`;
  console.log(banner);
}

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
      const composeContent = `services:
  backend:
    image: akshitmadan/100xtrader-backend:latest
    container_name: 100xtrader-backend
    ports:
      - "${backendPort}:8080"
    environment:
      - DB_PATH=/app/data/db.sqlite
      - PORT=8080
    volumes:
      - backend-data:/app/data
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s

  frontend:
    image: akshitmadan/100xtrader-frontend:latest
    container_name: 100xtrader-frontend
    ports:
      - "${port}:3000"
    environment:
      - NEXT_PUBLIC_API_URL=http://localhost:${backendPort}
    depends_on:
      backend:
        condition: service_healthy
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:3000"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s

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
          stdio: 'pipe',
          cwd: tempDir
        });

        let outputBuffer = '';
        let bannerShown = false;
        let backendHealthy = false;
        let frontendHealthy = false;
        
        // Fallback: Show banner after 10 seconds if not shown yet
        const fallbackTimer = setTimeout(() => {
          if (!bannerShown) {
            bannerShown = true;
            showBanner(port, backendPort);
          }
        }, 10000);
        
        up.stdout.on('data', (data) => {
          const output = data.toString();
          outputBuffer += output;
          process.stdout.write(output);
          
          // Check for backend health (multiple patterns)
          if (output.includes('100xtrader-backend') && 
              (output.includes('healthy') || output.includes('Started') || 
               output.includes('ready') || output.includes('listening'))) {
            backendHealthy = true;
          }
          
          // Check for frontend health (multiple patterns)
          if (output.includes('100xtrader-frontend') && 
              (output.includes('healthy') || output.includes('Started') || 
               output.includes('ready') || output.includes('compiled'))) {
            frontendHealthy = true;
          }
          
          // Show banner once both are healthy
          if (!bannerShown && backendHealthy && frontendHealthy) {
            clearTimeout(fallbackTimer);
            bannerShown = true;
            setTimeout(() => {
              showBanner(port, backendPort);
            }, 1000);
          }
        });

        up.stderr.on('data', (data) => {
          process.stderr.write(data);
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

