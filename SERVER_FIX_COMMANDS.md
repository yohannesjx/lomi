# Commands to Run on Your Server (SSH Session)

## Fix Port 8080 Error

Copy and paste these commands one by one in your SSH session:

```bash
# 1. Stop all containers
docker-compose -f docker-compose.prod.yml down

# 2. Kill any process using port 8080
sudo kill -9 $(sudo lsof -ti:8080) 2>/dev/null || true

# 3. Wait a moment
sleep 2

# 4. Check if port is free
sudo lsof -i :8080
# Should return nothing (port is free)

# 5. Deploy again
./deploy.sh
```

## Alternative: One-Line Fix

```bash
docker-compose -f docker-compose.prod.yml down && sudo kill -9 $(sudo lsof -ti:8080) 2>/dev/null || true && sleep 2 && ./deploy.sh
```

## If Still Not Working

```bash
# Check what's using the port
sudo netstat -tulpn | grep 8080
# or
sudo ss -tulpn | grep 8080

# Find and kill the process
sudo fuser -k 8080/tcp

# Remove any orphaned containers
docker container prune -f

# Try deploy again
./deploy.sh
```

## Check Docker Containers

```bash
# See all containers
docker ps -a

# Remove specific container if needed
docker rm -f lomi_backend

# Then deploy
./deploy.sh
```

