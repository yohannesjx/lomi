#!/bin/bash

# Fix Port 8080 Already in Use Error
# Run this on your server

echo "🔍 Checking what's using port 8080..."

# Find process using port 8080
PID=$(sudo lsof -ti:8080 2>/dev/null || sudo fuser 8080/tcp 2>/dev/null | awk '{print $1}')

if [ -z "$PID" ]; then
    echo "⚠️  Could not find process (might be Docker container)"
    echo "Checking Docker containers..."
    
    # Check for old containers
    OLD_CONTAINER=$(docker ps -a | grep lomi_backend | awk '{print $1}')
    if [ ! -z "$OLD_CONTAINER" ]; then
        echo "Found old container: $OLD_CONTAINER"
        echo "Stopping and removing..."
        docker stop $OLD_CONTAINER 2>/dev/null || true
        docker rm $OLD_CONTAINER 2>/dev/null || true
    fi
else
    echo "Found process using port 8080: PID $PID"
    echo "Killing process..."
    sudo kill -9 $PID 2>/dev/null || true
    sleep 2
fi

# Also check for any Docker containers using the port
echo "Checking Docker containers..."
docker ps -a | grep 8080

# Stop all lomi containers
echo "Stopping all Lomi containers..."
docker-compose -f docker-compose.prod.yml down 2>/dev/null || true

# Clean up any orphaned containers
echo "Cleaning up..."
docker container prune -f

echo ""
echo "✅ Port 8080 should be free now!"
echo ""
echo "Try deploying again:"
echo "   ./deploy.sh"

