#!/bin/bash

# Script to build and publish Docker images to Docker Hub

set -e

DOCKER_USERNAME="${DOCKER_USERNAME:-yourusername}"
VERSION="${1:-latest}"

echo "🐳 Building and publishing 100xtrader Docker images"
echo "=================================================="
echo "Docker Hub Username: $DOCKER_USERNAME"
echo "Version: $VERSION"
echo ""

# Check if logged in to Docker Hub
if ! docker info | grep -q "Username"; then
    echo "⚠️  Not logged in to Docker Hub"
    echo "Please run: docker login"
    exit 1
fi

# Build backend
echo "📦 Building backend image..."
docker build -f Dockerfile.backend \
    -t "$DOCKER_USERNAME/100xtrader-backend:$VERSION" \
    -t "$DOCKER_USERNAME/100xtrader-backend:latest" \
    .

# Build frontend
echo "📦 Building frontend image..."
docker build -f Dockerfile.frontend \
    -t "$DOCKER_USERNAME/100xtrader-frontend:$VERSION" \
    -t "$DOCKER_USERNAME/100xtrader-frontend:latest" \
    .

# Push backend
echo "🚀 Pushing backend image..."
docker push "$DOCKER_USERNAME/100xtrader-backend:$VERSION"
docker push "$DOCKER_USERNAME/100xtrader-backend:latest"

# Push frontend
echo "🚀 Pushing frontend image..."
docker push "$DOCKER_USERNAME/100xtrader-frontend:$VERSION"
docker push "$DOCKER_USERNAME/100xtrader-frontend:latest"

echo ""
echo "✅ Images published successfully!"
echo ""
echo "Users can now run:"
echo "  docker-compose -f docker-compose.public.yml up"
echo ""
echo "Or update docker-compose.public.yml with:"
echo "  image: $DOCKER_USERNAME/100xtrader-backend:latest"
echo "  image: $DOCKER_USERNAME/100xtrader-frontend:latest"

