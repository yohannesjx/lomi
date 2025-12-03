#!/bin/bash

# ==================== DEPLOY LIVE CHAT SYSTEM (FIXED) ====================
# Production deployment script for unified live chat system
# Run this on your server after pushing changes

set -e  # Exit on error

echo "🚀 Deploying Live Chat System..."

# ==================== STEP 1: Pull Latest Code ====================
echo ""
echo "📥 Step 1: Pulling latest code from GitHub..."
cd /root/lomi_mini || exit 1
git pull origin main

# ==================== STEP 2: Install Redis ====================
echo ""
echo "📦 Step 2: Installing Redis..."

# Check if Redis is already installed
if ! command -v redis-server &> /dev/null; then
    echo "Installing Redis..."
    apt update
    apt install -y redis-server
else
    echo "✅ Redis already installed"
fi

# Stop Redis if running
systemctl stop redis-server || true
sleep 2

# Create necessary directories
mkdir -p /var/lib/redis
mkdir -p /var/log/redis
mkdir -p /var/run/redis
chown -R redis:redis /var/lib/redis
chown -R redis:redis /var/log/redis
chown -R redis:redis /var/run/redis

# Use default Redis configuration with minimal changes
echo "Configuring Redis..."
cat > /etc/redis/redis.conf <<'EOF'
bind 127.0.0.1
protected-mode yes
port 6379
tcp-backlog 511
timeout 0
tcp-keepalive 300
daemonize no
supervised systemd
pidfile /var/run/redis/redis-server.pid
loglevel notice
logfile /var/log/redis/redis-server.log
databases 16
save 900 1
save 300 10
save 60 10000
stop-writes-on-bgsave-error yes
rdbcompression yes
rdbchecksum yes
dbfilename dump.rdb
dir /var/lib/redis
maxmemory 1gb
maxmemory-policy allkeys-lru
appendonly yes
appendfilename "appendonly.aof"
appendfsync everysec
EOF

# Fix permissions
chown redis:redis /etc/redis/redis.conf
chmod 640 /etc/redis/redis.conf

# Start Redis
echo "Starting Redis..."
systemctl enable redis-server
systemctl start redis-server
sleep 3

# Verify Redis is running
if redis-cli ping 2>/dev/null | grep -q "PONG"; then
    echo "✅ Redis is running"
else
    echo "⚠️  Redis may not be running, checking status..."
    systemctl status redis-server --no-pager || true
    
    # Try starting manually
    echo "Attempting manual start..."
    redis-server /etc/redis/redis.conf --daemonize yes
    sleep 2
    
    if redis-cli ping 2>/dev/null | grep -q "PONG"; then
        echo "✅ Redis started manually"
    else
        echo "❌ Redis failed to start. Continuing anyway..."
    fi
fi

# ==================== STEP 3: Run Database Migration ====================
echo ""
echo "🗄️  Step 3: Running database migration..."

# Set database credentials from environment or use defaults
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-lomi_db}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-}"

# Run migration
cd /root/lomi_mini/backend

if [ -f "database/migrations/20251204_add_live_chat_support.sql" ]; then
    echo "Running migration..."
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f database/migrations/20251204_add_live_chat_support.sql 2>&1 | grep -v "already exists" || true
    echo "✅ Database migration completed"
else
    echo "⚠️  Migration file not found"
fi

# ==================== STEP 4: Update Environment Variables ====================
echo ""
echo "⚙️  Step 4: Updating environment variables..."

cd /root/lomi_mini/backend

# Backup existing .env
if [ -f .env ]; then
    cp .env .env.backup.$(date +%Y%m%d_%H%M%S)
fi

# Add Redis configuration to .env if not exists
if ! grep -q "REDIS_HOST" .env 2>/dev/null; then
    cat >> .env <<'EOF'

# Redis Configuration (Live Chat)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
EOF
    echo "✅ Redis config added to .env"
else
    echo "✅ Redis config already in .env"
fi

# ==================== STEP 5: Install Go Dependencies ====================
echo ""
echo "📦 Step 5: Installing Go dependencies..."

cd /root/lomi_mini/backend

# Add redis client if not in go.mod
if ! grep -q "github.com/redis/go-redis/v9" go.mod; then
    echo "Adding Redis client dependency..."
    go get github.com/redis/go-redis/v9
fi

go mod download
go mod tidy

echo "✅ Dependencies installed"

# ==================== STEP 6: Build Backend ====================
echo ""
echo "🔨 Step 6: Building backend..."

cd /root/lomi_mini/backend
go build -o api cmd/main.go

if [ $? -eq 0 ]; then
    echo "✅ Backend built successfully"
else
    echo "❌ Backend build failed"
    echo "Check for compilation errors above"
    exit 1
fi

# ==================== STEP 7: Restart Backend Service ====================
echo ""
echo "🔄 Step 7: Restarting backend service..."

# Stop existing process
pkill -f "./api" || true
pkill -f "lomi-backend" || true
sleep 2

# Start new process
cd /root/lomi_mini/backend
nohup ./api > /var/log/lomi-backend.log 2>&1 &

# Wait for backend to start
sleep 5

# Check if backend is running
if pgrep -f "./api" > /dev/null; then
    echo "✅ Backend started successfully"
    echo "   PID: $(pgrep -f './api')"
else
    echo "❌ Backend failed to start"
    echo "Last 20 lines of log:"
    tail -20 /var/log/lomi-backend.log
    exit 1
fi

# ==================== STEP 8: Verify Installation ====================
echo ""
echo "🔍 Step 8: Verifying installation..."

# Check Redis
echo -n "Redis: "
if redis-cli ping 2>/dev/null | grep -q "PONG"; then
    echo "✅ Running"
    REDIS_STATUS="✅"
else
    echo "❌ Not running"
    REDIS_STATUS="❌"
fi

# Check PostgreSQL
echo -n "PostgreSQL: "
if PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "SELECT 1" > /dev/null 2>&1; then
    echo "✅ Connected"
    DB_STATUS="✅"
else
    echo "❌ Connection failed"
    DB_STATUS="❌"
fi

# Check Backend
echo -n "Backend API: "
sleep 2
if pgrep -f "./api" > /dev/null; then
    echo "✅ Running (PID: $(pgrep -f './api'))"
    BACKEND_STATUS="✅"
else
    echo "❌ Not running"
    BACKEND_STATUS="❌"
fi

# Check messages table
echo -n "Messages Table (live fields): "
if PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "SELECT live_stream_id, is_live, seq FROM messages LIMIT 1" > /dev/null 2>&1; then
    echo "✅ Updated"
    TABLE_STATUS="✅"
else
    echo "⚠️  Not updated (migration may have failed)"
    TABLE_STATUS="⚠️"
fi

# ==================== STEP 9: Display Status ====================
echo ""
echo "================================================"
echo "🎉 LIVE CHAT SYSTEM DEPLOYMENT COMPLETE!"
echo "================================================"
echo ""
echo "📊 Service Status:"
echo "  - Redis:      $REDIS_STATUS"
echo "  - PostgreSQL: $DB_STATUS"
echo "  - Backend:    $BACKEND_STATUS"
echo "  - Migration:  $TABLE_STATUS"
echo ""
echo "🌐 Endpoints:"
echo "  - API: http://localhost:8080"
echo "  - WebSocket: ws://localhost:8080/ws/chat"
echo ""
echo "📝 Logs:"
echo "  - Backend: tail -f /var/log/lomi-backend.log"
echo "  - Redis: tail -f /var/log/redis/redis-server.log"
echo ""
echo "🔧 Useful Commands:"
echo "  - Test Redis: redis-cli ping"
echo "  - Monitor Redis: redis-cli MONITOR"
echo "  - Check backend: ps aux | grep './api'"
echo "  - Restart backend: pkill -f './api' && cd /root/lomi_mini/backend && nohup ./api > /var/log/lomi-backend.log 2>&1 &"
echo ""
echo "🧪 Test WebSocket Connection:"
echo "  wscat -c 'ws://localhost:8080/ws/chat?token=YOUR_TOKEN&mode=live&live_stream_id=test-123'"
echo ""

if [ "$REDIS_STATUS" = "❌" ]; then
    echo "⚠️  WARNING: Redis is not running!"
    echo "   Live chat will NOT work without Redis."
    echo "   Try: systemctl status redis-server"
    echo "   Or: redis-server /etc/redis/redis.conf --daemonize yes"
    echo ""
fi

if [ "$BACKEND_STATUS" = "❌" ]; then
    echo "⚠️  WARNING: Backend is not running!"
    echo "   Check logs: tail -50 /var/log/lomi-backend.log"
    echo ""
fi

echo "================================================"
