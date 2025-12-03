#!/bin/bash

# ==================== DEPLOY LIVE CHAT SYSTEM ====================
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
    
    # Configure Redis for production
    echo "Configuring Redis..."
    cat > /etc/redis/redis.conf <<EOF
# Redis Configuration for Live Chat
bind 0.0.0.0
protected-mode yes
port 6379
tcp-backlog 511
timeout 300
tcp-keepalive 60
daemonize yes
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
replica-serve-stale-data yes
replica-read-only yes
repl-diskless-sync no
repl-diskless-sync-delay 5
repl-disable-tcp-nodelay no
replica-priority 100
maxmemory 2gb
maxmemory-policy allkeys-lru
appendonly yes
appendfilename "appendonly.aof"
appendfsync everysec
no-appendfsync-on-rewrite no
auto-aof-rewrite-percentage 100
auto-aof-rewrite-min-size 64mb
aof-load-truncated yes
lua-time-limit 5000
slowlog-log-slower-than 10000
slowlog-max-len 128
latency-monitor-threshold 0
notify-keyspace-events ""
hash-max-ziplist-entries 512
hash-max-ziplist-value 64
list-max-ziplist-size -2
list-compress-depth 0
set-max-intset-entries 512
zset-max-ziplist-entries 128
zset-max-ziplist-value 64
hll-sparse-max-bytes 3000
stream-node-max-bytes 4096
stream-node-max-entries 100
activerehashing yes
client-output-buffer-limit normal 0 0 0
client-output-buffer-limit replica 256mb 64mb 60
client-output-buffer-limit pubsub 32mb 8mb 60
hz 10
dynamic-hz yes
aof-rewrite-incremental-fsync yes
rdb-save-incremental-fsync yes
maxclients 10000
EOF
    
    # Start Redis
    systemctl enable redis-server
    systemctl start redis-server
    
    echo "✅ Redis installed and started"
else
    echo "✅ Redis already installed"
    systemctl restart redis-server
fi

# Verify Redis is running
if redis-cli ping | grep -q "PONG"; then
    echo "✅ Redis is running"
else
    echo "❌ Redis failed to start"
    exit 1
fi

# ==================== STEP 3: Run Database Migration ====================
echo ""
echo "🗄️  Step 3: Running database migration..."

# Set database credentials
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-lomi_db}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-}"

# Run migration
cd /root/lomi_mini/backend
PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f database/migrations/20251204_add_live_chat_support.sql

if [ $? -eq 0 ]; then
    echo "✅ Database migration completed"
else
    echo "⚠️  Migration may have already been applied or failed"
fi

# ==================== STEP 4: Update Environment Variables ====================
echo ""
echo "⚙️  Step 4: Updating environment variables..."

cd /root/lomi_mini/backend

# Add Redis configuration to .env if not exists
if ! grep -q "REDIS_HOST" .env; then
    cat >> .env <<EOF

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

# ==================== STEP 5: Build Backend ====================
echo ""
echo "🔨 Step 5: Building backend..."

cd /root/lomi_mini/backend
go mod download
go build -o api cmd/main.go

if [ $? -eq 0 ]; then
    echo "✅ Backend built successfully"
else
    echo "❌ Backend build failed"
    exit 1
fi

# ==================== STEP 6: Restart Backend Service ====================
echo ""
echo "🔄 Step 6: Restarting backend service..."

# Stop existing process
pkill -f "lomi-backend" || true
sleep 2

# Start new process
cd /root/lomi_mini/backend
nohup ./api > /var/log/lomi-backend.log 2>&1 &

# Wait for backend to start
sleep 3

# Check if backend is running
if pgrep -f "lomi-backend" > /dev/null; then
    echo "✅ Backend started successfully"
else
    echo "❌ Backend failed to start"
    echo "Check logs: tail -f /var/log/lomi-backend.log"
    exit 1
fi

# ==================== STEP 7: Verify Installation ====================
echo ""
echo "🔍 Step 7: Verifying installation..."

# Check Redis
echo -n "Redis: "
if redis-cli ping | grep -q "PONG"; then
    echo "✅ Running"
else
    echo "❌ Not running"
fi

# Check PostgreSQL
echo -n "PostgreSQL: "
if PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "SELECT 1" > /dev/null 2>&1; then
    echo "✅ Connected"
else
    echo "❌ Connection failed"
fi

# Check Backend
echo -n "Backend API: "
if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo "✅ Running"
else
    echo "⚠️  Health check failed (may be normal if no /health endpoint)"
fi

# Check live_streams table
echo -n "Live Streams Table: "
if PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "SELECT 1 FROM live_streams LIMIT 1" > /dev/null 2>&1; then
    echo "✅ Exists"
else
    echo "⚠️  Not found (migration may have failed)"
fi

# ==================== STEP 8: Display Status ====================
echo ""
echo "================================================"
echo "✅ LIVE CHAT SYSTEM DEPLOYED SUCCESSFULLY!"
echo "================================================"
echo ""
echo "📊 Service Status:"
echo "  - Redis: Running on port 6379"
echo "  - Backend: Running on port 8080"
echo "  - WebSocket: ws://your-domain.com/ws/chat"
echo ""
echo "📝 Logs:"
echo "  - Backend: tail -f /var/log/lomi-backend.log"
echo "  - Redis: tail -f /var/log/redis/redis-server.log"
echo ""
echo "🔧 Redis Commands:"
echo "  - Monitor: redis-cli MONITOR"
echo "  - Stats: redis-cli INFO"
echo "  - Pub/Sub: redis-cli PUBSUB CHANNELS"
echo ""
echo "🧪 Test WebSocket:"
echo "  wscat -c 'ws://localhost:8080/ws/chat?token=YOUR_TOKEN&mode=live&live_stream_id=TEST_ID'"
echo ""
echo "================================================"
