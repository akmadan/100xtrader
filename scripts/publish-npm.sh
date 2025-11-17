#!/bin/bash

# Script to publish npm package for npx distribution

set -e

echo "📦 Preparing npm package for 100xTrader"
echo "========================================"
echo ""

# Create temporary directory
TEMP_DIR=$(mktemp -d)
echo "📁 Using temporary directory: $TEMP_DIR"

# Copy necessary files
echo "📋 Copying files..."
cp bin/100xtrader.js "$TEMP_DIR/"
cp docker-compose.public.yml "$TEMP_DIR/"
cp package.json.public "$TEMP_DIR/package.json"
cp README.md "$TEMP_DIR/"

# Make bin executable
chmod +x "$TEMP_DIR/bin/100xtrader.js"

# Check if logged in to npm
if ! npm whoami &> /dev/null; then
    echo "⚠️  Not logged in to npm"
    echo "Please run: npm login"
    exit 1
fi

# Publish
echo "🚀 Publishing to npm..."
cd "$TEMP_DIR"
npm publish --access public

echo ""
echo "✅ Package published successfully!"
echo ""
echo "Users can now run:"
echo "  npx 100xtrader"
echo ""
echo "Or install globally:"
echo "  npm install -g 100xtrader"
echo "  100xtrader"

# Cleanup
cd - > /dev/null
rm -rf "$TEMP_DIR"

