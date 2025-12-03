# 🚀 TikTok Clone - Deployment Checklist

## ✅ Pre-Deployment Checklist

### 1. Backend Setup
- [ ] Go backend code compiled successfully
- [ ] All dependencies installed (`go mod tidy`)
- [ ] Database connection tested
- [ ] Environment variables configured
- [ ] Test script passes (`./test-tiktok-api.sh`)

### 2. Database
- [ ] PostgreSQL running
- [ ] All required tables exist (users, gifts, transactions)
- [ ] Database backups configured
- [ ] Connection string in `.env` file

### 3. Live Streaming
- [ ] MediaMTX downloaded and installed
- [ ] MediaMTX running on ports 1935 (RTMP) and 8888 (HLS)
- [ ] Test stream successful

### 4. Mobile Apps
- [ ] Android app base URL updated
- [ ] iOS app base URL updated
- [ ] Apps built successfully
- [ ] Test on real devices

---

## 🔧 Local Testing (Development)

### Step 1: Start Backend
```bash
cd backend
go run cmd/api/main.go
```
✅ Server running on http://localhost:8080

### Step 2: Test Endpoints
```bash
./test-tiktok-api.sh
```
✅ All 6 tests pass

### Step 3: Start MediaMTX (Optional)
```bash
cd mediamtx
./mediamtx
```
✅ RTMP on port 1935, HLS on port 8888

### Step 4: Update Apps
**Android:** `ApiLinks.java`
```java
public static String API_BASE_URL = "http://YOUR_LOCAL_IP:8080/api/";
```

**iOS:** `ProductEndPoint.swift`
```swift
var baseURL: String {
    return "http://YOUR_LOCAL_IP:8080/api/"
}
```

### Step 5: Build & Test Apps
```bash
# Android
cd android
./gradlew assembleDebug

# iOS
cd ios
xcodebuild -workspace VideoSmash.xcworkspace -scheme VideoSmash build
```

✅ Apps connect to backend  
✅ Registration works  
✅ Video feed loads  
✅ Live streaming works  
✅ Gifts can be sent  
✅ Coins can be purchased  

---

## 🌐 Production Deployment

### Step 1: Server Setup
```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Go (if not installed)
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Install PostgreSQL (if not installed)
sudo apt install postgresql postgresql-contrib

# Install Nginx
sudo apt install nginx

# Install Certbot for SSL
sudo apt install certbot python3-certbot-nginx
```

### Step 2: Deploy Backend
```bash
# Clone/upload your code
cd /opt
git clone YOUR_REPO lomi-backend
cd lomi-backend/backend

# Build binary
go build -o lomi-backend cmd/api/main.go

# Create systemd service
sudo nano /etc/systemd/system/lomi-backend.service
```

**Service file:**
```ini
[Unit]
Description=Lomi Backend API
After=network.target postgresql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/lomi-backend/backend
ExecStart=/opt/lomi-backend/backend/lomi-backend
Restart=always
RestartSec=5
Environment="APP_ENV=production"
EnvironmentFile=/opt/lomi-backend/backend/.env.production

[Install]
WantedBy=multi-user.target
```

```bash
# Start service
sudo systemctl daemon-reload
sudo systemctl enable lomi-backend
sudo systemctl start lomi-backend
sudo systemctl status lomi-backend
```

✅ Backend running as systemd service

### Step 3: Deploy MediaMTX
```bash
# Download MediaMTX
cd /opt
wget https://github.com/bluenviron/mediamtx/releases/download/v1.3.0/mediamtx_v1.3.0_linux_amd64.tar.gz
tar -xzf mediamtx_v1.3.0_linux_amd64.tar.gz
mv mediamtx_v1.3.0_linux_amd64 mediamtx

# Create systemd service
sudo nano /etc/systemd/system/mediamtx.service
```

**Service file:**
```ini
[Unit]
Description=MediaMTX RTMP/HLS Server
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/mediamtx
ExecStart=/opt/mediamtx/mediamtx
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

```bash
# Start service
sudo systemctl daemon-reload
sudo systemctl enable mediamtx
sudo systemctl start mediamtx
sudo systemctl status mediamtx
```

✅ MediaMTX running as systemd service

### Step 4: Configure Nginx
```bash
sudo nano /etc/nginx/sites-available/lomi-api
```

**Nginx config:**
```nginx
server {
    listen 80;
    server_name api.lomilive.com;

    # API endpoints
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }

    # HLS streaming
    location /live/ {
        proxy_pass http://localhost:8888;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        add_header Cache-Control no-cache;
        add_header Access-Control-Allow-Origin *;
    }
}
```

```bash
# Enable site
sudo ln -s /etc/nginx/sites-available/lomi-api /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

✅ Nginx configured

### Step 5: SSL Certificate
```bash
sudo certbot --nginx -d api.lomilive.com
```

✅ SSL certificate installed

### Step 6: Firewall
```bash
# Allow necessary ports
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 1935/tcp  # RTMP
sudo ufw allow 8888/tcp  # HLS
sudo ufw enable
```

✅ Firewall configured

### Step 7: Update Mobile Apps
**Android:** `ApiLinks.java`
```java
public static String API_BASE_URL = "https://api.lomilive.com/api/";
```

**iOS:** `ProductEndPoint.swift`
```swift
var baseURL: String {
    return "https://api.lomilive.com/api/"
}
```

### Step 8: Build Production Apps
```bash
# Android (Release)
cd android
./gradlew assembleRelease
# APK at: app/build/outputs/apk/release/app-release.apk

# iOS (Archive)
cd ios
xcodebuild -workspace VideoSmash.xcworkspace \
  -scheme VideoSmash \
  -configuration Release \
  archive -archivePath build/VideoSmash.xcarchive
```

✅ Production apps built

---

## 🧪 Production Testing

### Test 1: Health Check
```bash
curl https://api.lomilive.com/api/v1/health
```
Expected: `{"status":"ok"}`

### Test 2: Register User
```bash
curl -X POST https://api.lomilive.com/api/registerUser \
  -H "Content-Type: application/json" \
  -d '{
    "username":"testuser",
    "email":"test@example.com",
    "social_id":"google_123",
    "social":"google"
  }'
```
Expected: `{"code":200,"msg":{...}}`

### Test 3: Video Feed
```bash
curl -X POST https://api.lomilive.com/api/showRelatedVideos \
  -H "Content-Type: application/json" \
  -d '{
    "user_id":"1",
    "device_id":"device_123",
    "starting_point":0
  }'
```
Expected: Array of 5 videos

### Test 4: Live Streaming
```bash
# Start stream
curl -X POST https://api.lomilive.com/api/liveStream \
  -H "Content-Type: application/json" \
  -d '{
    "user_id":"uuid",
    "started_at":"2024-12-03 16:00:00"
  }'

# Test RTMP (with FFmpeg)
ffmpeg -re -i test.mp4 -c copy -f flv rtmp://api.lomilive.com:1935/live/test_stream

# Test HLS playback
curl https://api.lomilive.com:8888/live/test_stream/index.m3u8
```

✅ All production tests pass

---

## 📊 Monitoring

### Check Service Status
```bash
# Backend
sudo systemctl status lomi-backend

# MediaMTX
sudo systemctl status mediamtx

# Nginx
sudo systemctl status nginx

# PostgreSQL
sudo systemctl status postgresql
```

### View Logs
```bash
# Backend logs
sudo journalctl -u lomi-backend -f

# MediaMTX logs
sudo journalctl -u mediamtx -f

# Nginx logs
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log
```

### Database Monitoring
```bash
# Connect to database
psql -U lomi -d lomi_db

# Check user count
SELECT COUNT(*) FROM users;

# Check recent transactions
SELECT * FROM coin_transactions ORDER BY created_at DESC LIMIT 10;

# Check gift transactions
SELECT * FROM gift_transactions ORDER BY created_at DESC LIMIT 10;
```

---

## 🔒 Security Checklist

- [ ] SSL certificate installed and auto-renewing
- [ ] Firewall configured (UFW)
- [ ] Database password strong and secure
- [ ] JWT secret is random and secure
- [ ] API keys not exposed in client apps
- [ ] Rate limiting enabled
- [ ] CORS configured properly
- [ ] Database backups automated
- [ ] Server logs monitored
- [ ] Fail2ban installed (optional)

---

## 📱 App Store Submission

### Android (Google Play)
- [ ] App signed with release key
- [ ] Version code incremented
- [ ] Privacy policy URL added
- [ ] Screenshots prepared
- [ ] App description written
- [ ] APK uploaded to Play Console

### iOS (App Store)
- [ ] App signed with distribution certificate
- [ ] Version number incremented
- [ ] Privacy policy URL added
- [ ] Screenshots prepared
- [ ] App description written
- [ ] IPA uploaded to App Store Connect

---

## 🎉 Launch Checklist

### Pre-Launch
- [ ] All endpoints tested in production
- [ ] Mobile apps tested on real devices
- [ ] Live streaming tested end-to-end
- [ ] Gift sending tested
- [ ] Coin purchase tested
- [ ] Database backups verified
- [ ] Monitoring set up
- [ ] Support email configured

### Launch Day
- [ ] Backend running smoothly
- [ ] MediaMTX running smoothly
- [ ] Apps submitted to stores
- [ ] Marketing materials ready
- [ ] Support team ready
- [ ] Analytics tracking enabled

### Post-Launch
- [ ] Monitor server resources
- [ ] Check error logs
- [ ] Track user registrations
- [ ] Monitor transaction volume
- [ ] Gather user feedback
- [ ] Plan next features

---

## 📞 Emergency Contacts

### Server Issues
```bash
# Restart backend
sudo systemctl restart lomi-backend

# Restart MediaMTX
sudo systemctl restart mediamtx

# Restart Nginx
sudo systemctl restart nginx

# Check disk space
df -h

# Check memory
free -h

# Check CPU
top
```

### Database Issues
```bash
# Restart PostgreSQL
sudo systemctl restart postgresql

# Check connections
psql -U lomi -d lomi_db -c "SELECT count(*) FROM pg_stat_activity;"

# Vacuum database
psql -U lomi -d lomi_db -c "VACUUM ANALYZE;"
```

---

## ✅ Final Verification

Before going live, verify:

- [ ] ✅ Backend responds to all 6 endpoints
- [ ] ✅ SSL certificate valid
- [ ] ✅ Apps connect successfully
- [ ] ✅ User registration works
- [ ] ✅ Video feed loads
- [ ] ✅ Live streaming works
- [ ] ✅ Gifts can be sent
- [ ] ✅ Coins can be purchased
- [ ] ✅ Database transactions recorded
- [ ] ✅ Logs are clean (no errors)

---

**You're ready to launch!** 🚀

Good luck with your TikTok clone! 🎬
