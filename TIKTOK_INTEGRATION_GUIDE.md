# TikTok Clone Integration Guide

## 🎯 Overview

This guide explains how to integrate the TikTok-style Android/iOS apps with your existing Lomi Go backend.

## ✅ What's Been Implemented

### 6 Critical Endpoints (Ready to Use)

All endpoints match the exact API contract from the reverse-engineered TikTok apps:

1. **POST /api/registerUser** - Social login/signup
2. **POST /api/showUserDetail** - Get user profile & wallet balance
3. **POST /api/showRelatedVideos** - Home feed (For You page) with 5 dummy videos
4. **POST /api/liveStream** - Start live streaming session
5. **POST /api/sendGift** - Send virtual gifts
6. **POST /api/purchaseCoin** - Buy coins (in-app purchase)

### Bonus Endpoints

7. **POST /api/showCoinWorth** - Get coin packages
8. **POST /api/showGifts** - Get gifts catalog

---

## 📁 Files Created

```
backend/
├── internal/
│   ├── handlers/
│   │   └── streaming.go          # All 6 endpoint handlers
│   └── routes/
│       └── streaming_routes.go   # Route definitions
└── cmd/api/main.go                # Updated to include streaming routes
```

---

## 🚀 Quick Start

### 1. Update Android/iOS App Configuration

Update the base URL in both apps to point to your backend:

**Android:** `app/src/main/java/com/qboxus/tictic/apiclasses/ApiLinks.java`
```java
public static String API_BASE_URL = "https://api.lomilive.com/api/";
```

**iOS:** `VideoSmash/API/EndPoint/ProductEndPoint.swift`
```swift
var baseURL: String {
    return "https://api.lomilive.com/api/"
}
```

### 2. Build and Run Backend

```bash
cd backend
go mod tidy
go run cmd/api/main.go
```

The server will start on port 8080 (or your configured port).

### 3. Test the Endpoints

```bash
# Test registration
curl -X POST http://localhost:8080/api/registerUser \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "first_name": "Test",
    "last_name": "User",
    "email": "test@example.com",
    "social_id": "google_12345",
    "social": "google",
    "device_token": "fcm_token_here"
  }'

# Test video feed
curl -X POST http://localhost:8080/api/showRelatedVideos \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "1",
    "device_id": "device_123",
    "starting_point": 0,
    "lat": 9.0320,
    "long": 38.7469
  }'
```

---

## 🔐 Authentication Flow

### How It Works

1. **App Opens** → Calls `/api/registerUser` with social login credentials
2. **Backend** → Creates/finds user, returns `auth_token` in response
3. **App** → Stores `auth_token` in SharedPreferences/UserDefaults
4. **Subsequent Requests** → App includes `auth_token` in request body

### Current Implementation

- The endpoints currently accept `auth_token` in the request body (TikTok style)
- User ID is extracted from the token or request
- For production, you should validate the token properly

### Recommended Enhancement

Add proper JWT validation:

```go
// In streaming.go, add this helper:
func getUserIDFromToken(authToken string) (uuid.UUID, error) {
    // Parse and validate JWT token
    // Return user ID
}
```

---

## 💰 Monetization Flow

### Coin System Integration

The implementation reuses your existing coin/wallet system:

1. **User buys coins** → `/api/purchaseCoin`
   - Creates `CoinTransaction` with type `purchase`
   - Updates `User.CoinBalance`

2. **User sends gift** → `/api/sendGift`
   - Creates `GiftTransaction`
   - Deducts from sender's `CoinBalance`
   - Adds to receiver's `CoinBalance` and `TotalEarned`
   - Creates two `CoinTransaction` records (sent/received)

3. **User withdraws** → Use existing `/api/v1/payouts/request`

### Gift Catalog

The dummy gifts endpoint returns 4 gifts. To use real gifts from your database:

```go
// In streaming_routes.go, replace the dummy data:
api.Post("/showGifts", func(c *fiber.Ctx) error {
    var gifts []models.Gift
    database.DB.Where("is_active = ?", true).
        Order("display_order ASC").
        Find(&gifts)
    
    // Transform to TikTok format...
})
```

---

## 📹 Live Streaming Setup

### Option 1: MediaMTX (Recommended)

MediaMTX is a modern, lightweight RTMP/HLS server.

#### Install MediaMTX

```bash
# Download latest release
wget https://github.com/bluenviron/mediamtx/releases/download/v1.3.0/mediamtx_v1.3.0_linux_amd64.tar.gz
tar -xzf mediamtx_v1.3.0_linux_amd64.tar.gz
cd mediamtx

# Run MediaMTX
./mediamtx
```

#### Configure MediaMTX

Edit `mediamtx.yml`:

```yaml
# RTMP server
rtmpAddress: :1935

# HLS server
hlsAddress: :8888
hlsAlwaysRemux: yes

# Paths (streams)
paths:
  live:
    source: publisher
    runOnPublish: curl http://localhost:8080/api/liveStreamStarted
    runOnUnPublish: curl http://localhost:8080/api/liveStreamEnded
```

#### How It Works

1. **App calls** `/api/liveStream` → Backend returns streaming ID
2. **App streams to** `rtmp://your-server:1935/live/{streaming_id}`
3. **Viewers watch at** `http://your-server:8888/live/{streaming_id}/index.m3u8`

### Option 2: Nginx-RTMP

```bash
# Install nginx with RTMP module
sudo apt install nginx libnginx-mod-rtmp

# Configure /etc/nginx/nginx.conf
rtmp {
    server {
        listen 1935;
        chunk_size 4096;

        application live {
            live on;
            record off;
            
            # HLS
            hls on;
            hls_path /tmp/hls;
            hls_fragment 3;
            hls_playlist_length 60;
        }
    }
}

# Restart nginx
sudo systemctl restart nginx
```

---

## 🎥 Video Feed Implementation

### Current: Dummy Videos

The `/api/showRelatedVideos` endpoint currently returns 5 dummy videos with:
- Sample video URLs from Google Cloud Storage
- Random thumbnails from Picsum
- Fake engagement metrics

### Production: Real Videos

To implement real video feed:

1. **Create Video Model**

```go
// backend/internal/models/video.go
type Video struct {
    ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
    UserID      uuid.UUID `gorm:"type:uuid;not null"`
    Description string    `gorm:"type:text"`
    VideoURL    string    `gorm:"type:text;not null"`
    ThumbnailURL string   `gorm:"type:text"`
    SoundID     *uuid.UUID
    ViewCount   int       `gorm:"default:0"`
    LikeCount   int       `gorm:"default:0"`
    CommentCount int      `gorm:"default:0"`
    ShareCount  int       `gorm:"default:0"`
    PrivacyType string    `gorm:"default:'public'"`
    AllowComments bool    `gorm:"default:true"`
    AllowDuet   bool      `gorm:"default:true"`
    CreatedAt   time.Time
}
```

2. **Update Handler**

```go
func (h *StreamingHandler) ShowRelatedVideos(c *fiber.Ctx) error {
    var videos []models.Video
    
    // Get videos with pagination
    database.DB.
        Preload("User").
        Preload("Sound").
        Where("privacy_type = ?", "public").
        Order("created_at DESC").
        Offset(req.StartingPoint).
        Limit(10).
        Find(&videos)
    
    // Transform to TikTok format...
}
```

3. **Add Video Upload Endpoint**

```go
// POST /api/postVideo (multipart)
func (h *StreamingHandler) PostVideo(c *fiber.Ctx) error {
    // Handle multipart file upload
    // Save to S3/R2
    // Create Video record
    // Return TikTok-style response
}
```

---

## 🔧 Database Migrations

### Required Tables

Most tables already exist in your database. You may need to add:

```sql
-- Videos table
CREATE TABLE videos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    description TEXT,
    video_url TEXT NOT NULL,
    thumbnail_url TEXT,
    sound_id UUID,
    view_count INTEGER DEFAULT 0,
    like_count INTEGER DEFAULT 0,
    comment_count INTEGER DEFAULT 0,
    share_count INTEGER DEFAULT 0,
    privacy_type VARCHAR(20) DEFAULT 'public',
    allow_comments BOOLEAN DEFAULT true,
    allow_duet BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Live streaming sessions
CREATE TABLE live_streaming_sessions (
    id VARCHAR(255) PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    channel_name VARCHAR(255) NOT NULL,
    started_at TIMESTAMPTZ NOT NULL,
    ended_at TIMESTAMPTZ,
    status VARCHAR(20) DEFAULT 'live',
    viewer_count INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Video likes
CREATE TABLE video_likes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    video_id UUID NOT NULL REFERENCES videos(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, video_id)
);

-- Video comments
CREATE TABLE video_comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    video_id UUID NOT NULL REFERENCES videos(id),
    parent_id UUID REFERENCES video_comments(id),
    comment TEXT NOT NULL,
    like_count INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

---

## 📱 App Configuration

### Android App

1. **Update API Base URL**
   - File: `app/src/main/java/com/qboxus/tictic/apiclasses/ApiLinks.java`
   - Change: `API_BASE_URL = "https://api.lomilive.com/api/"`

2. **Update Constants**
   - File: `app/src/main/java/com/qboxus/tictic/Constants.java`
   - Set your API key if needed

3. **Build APK**
   ```bash
   cd android
   ./gradlew assembleRelease
   ```

### iOS App

1. **Update API Base URL**
   - File: `VideoSmash/API/EndPoint/ProductEndPoint.swift`
   - Change: `baseURL = "https://api.lomilive.com/api/"`

2. **Update Info.plist**
   - Add your domain to allowed network requests

3. **Build IPA**
   ```bash
   cd ios
   xcodebuild -workspace VideoSmash.xcworkspace -scheme VideoSmash archive
   ```

---

## 🌐 Deployment

### 1. Backend Deployment

```bash
# Build binary
cd backend
go build -o lomi-backend cmd/api/main.go

# Run with systemd
sudo systemctl start lomi-backend
```

### 2. MediaMTX Deployment

```bash
# Run as systemd service
sudo systemctl start mediamtx
```

### 3. Nginx Configuration

```nginx
# /etc/nginx/sites-available/lomi-api
server {
    listen 80;
    server_name api.lomilive.com;

    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location /live/ {
        proxy_pass http://localhost:8888;
    }
}
```

### 4. SSL Certificate

```bash
sudo certbot --nginx -d api.lomilive.com
```

---

## 🧪 Testing

### Test Registration

```bash
curl -X POST https://api.lomilive.com/api/registerUser \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "first_name": "Test",
    "last_name": "User",
    "email": "test@example.com",
    "social_id": "google_12345",
    "social": "google"
  }'
```

Expected response:
```json
{
  "code": 200,
  "msg": {
    "User": {
      "id": "...",
      "auth_token": "auth_...",
      "wallet": 0
    }
  }
}
```

### Test Video Feed

```bash
curl -X POST https://api.lomilive.com/api/showRelatedVideos \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "1",
    "device_id": "device_123",
    "starting_point": 0
  }'
```

### Test Live Streaming

```bash
# Start stream
curl -X POST https://api.lomilive.com/api/liveStream \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "...",
    "started_at": "2024-12-03 16:00:00"
  }'

# Stream with OBS/FFmpeg
ffmpeg -re -i input.mp4 -c copy -f flv rtmp://api.lomilive.com:1935/live/{streaming_id}

# Watch stream
vlc http://api.lomilive.com:8888/live/{streaming_id}/index.m3u8
```

---

## 📊 Monitoring

### Check Logs

```bash
# Backend logs
tail -f /var/log/lomi-backend.log

# MediaMTX logs
tail -f /var/log/mediamtx.log

# Nginx logs
tail -f /var/log/nginx/access.log
```

### Database Queries

```sql
-- Check user registrations
SELECT COUNT(*) FROM users WHERE created_at > NOW() - INTERVAL '1 day';

-- Check coin purchases
SELECT SUM(coin_amount) FROM coin_transactions 
WHERE transaction_type = 'purchase' 
AND created_at > NOW() - INTERVAL '1 day';

-- Check gifts sent
SELECT COUNT(*), SUM(coin_amount) FROM gift_transactions 
WHERE created_at > NOW() - INTERVAL '1 day';
```

---

## 🔒 Security Recommendations

### 1. Implement Proper JWT Validation

```go
// Add to streaming.go
func validateAuthToken(token string) (uuid.UUID, error) {
    claims := jwt.MapClaims{}
    _, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
        return []byte(config.Cfg.JWTSecret), nil
    })
    if err != nil {
        return uuid.Nil, err
    }
    
    userIDStr := claims["user_id"].(string)
    return uuid.Parse(userIDStr)
}
```

### 2. Add Rate Limiting

```go
// Use existing middleware
protected := api.Group("", middleware.AuthMiddleware)
protected.Post("/sendGift", middleware.GiftRateLimit(), streamingHandler.SendGift)
```

### 3. Validate In-App Purchases

```go
// For iOS
func verifyAppleReceipt(receiptData string) error {
    // Call Apple's verification API
}

// For Android
func verifyGooglePlayPurchase(purchaseToken string) error {
    // Call Google Play verification API
}
```

---

## 🚀 Next Steps

### Phase 1: Core Features (Week 1)
- ✅ User registration/login
- ✅ User profile & wallet
- ✅ Dummy video feed
- ✅ Live streaming setup
- ✅ Gift sending
- ✅ Coin purchases

### Phase 2: Video System (Week 2)
- [ ] Video upload endpoint
- [ ] Video storage (S3/R2)
- [ ] Video transcoding
- [ ] Real video feed algorithm
- [ ] Video likes/comments

### Phase 3: Social Features (Week 3)
- [ ] Follow/unfollow
- [ ] Notifications
- [ ] Comments system
- [ ] Hashtags
- [ ] Sounds library

### Phase 4: Advanced Features (Week 4)
- [ ] Video analytics
- [ ] Recommendation algorithm
- [ ] Search functionality
- [ ] Direct messaging
- [ ] Push notifications

---

## 📞 Support

If you encounter issues:

1. Check logs: `tail -f /var/log/lomi-backend.log`
2. Verify database connection: `psql -U lomi -d lomi_db`
3. Test endpoints with curl
4. Check MediaMTX status: `systemctl status mediamtx`

---

## 🎉 Success Checklist

- [ ] Backend running on port 8080
- [ ] MediaMTX running on ports 1935 (RTMP) and 8888 (HLS)
- [ ] Android app connects and shows video feed
- [ ] iOS app connects and shows video feed
- [ ] User can register/login
- [ ] User can view dummy videos
- [ ] User can start live stream
- [ ] User can send gifts
- [ ] User can buy coins

---

**You're all set!** 🚀

The Android and iOS apps should now work with your Go backend. Start with the dummy video feed, then gradually implement real video upload and streaming features.
