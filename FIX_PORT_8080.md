# Fix Port 8080 Error

## Quick Fix (Run on Server)

```bash
# Option 1: Stop all containers and free the port
docker-compose -f docker-compose.prod.yml down
sudo lsof -ti:8080 | xargs sudo kill -9 2>/dev/null || true

# Option 2: Use the fix script
chmod +x fix-port-8080.sh
./fix-port-8080.sh

# Then deploy again
./deploy.sh
```

## What's Happening

Port 8080 is already in use by:
- An old Docker container
- A previous backend instance
- Another service

## Solution

The `deploy.sh` script now automatically handles this, but if you get the error:

```bash
# Find what's using the port
sudo lsof -i :8080

# Kill it
sudo kill -9 $(sudo lsof -ti:8080)

# Or stop all containers
docker-compose -f docker-compose.prod.yml down

# Then deploy
./deploy.sh
```

