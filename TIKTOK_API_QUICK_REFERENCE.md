# TikTok API Quick Reference

## 🎯 Base URL
```
https://api.lomilive.com/api/
```

## 📋 All Endpoints

### 1. POST /api/registerUser
**Purpose:** Social login/signup  
**Auth:** None (public)

**Request:**
```json
{
  "username": "johndoe",
  "first_name": "John",
  "last_name": "Doe",
  "email": "john@example.com",
  "phone": "+1234567890",
  "social_id": "google_12345",
  "social": "google",
  "device_token": "fcm_token_here"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": {
    "User": {
      "id": "uuid",
      "auth_token": "auth_...",
      "wallet": 0,
      "verified": 0
    }
  }
}
```

---

### 2. POST /api/showUserDetail
**Purpose:** Get user profile & wallet  
**Auth:** Required

**Request:**
```json
{
  "auth_token": "auth_...",
  "user_id": "uuid",
  "other_user_id": "uuid" // optional
}
```

**Response:**
```json
{
  "code": 200,
  "msg": {
    "User": {
      "id": "uuid",
      "username": "johndoe",
      "wallet": 1500,
      "verified": 1,
      "button": "following"
    }
  }
}
```

---

### 3. POST /api/showRelatedVideos
**Purpose:** Home feed (For You page)  
**Auth:** Required

**Request:**
```json
{
  "user_id": "uuid",
  "device_id": "device_123",
  "starting_point": 0,
  "lat": 9.0320,
  "long": 38.7469
}
```

**Response:**
```json
{
  "code": 200,
  "msg": [
    {
      "Video": {
        "id": 1,
        "description": "Amazing video!",
        "video": "https://...",
        "thumbnail": "https://...",
        "view": 15000,
        "like": 1200
      },
      "User": {
        "id": 1,
        "username": "creator",
        "verified": 1
      }
    }
  ]
}
```

---

### 4. POST /api/liveStream
**Purpose:** Start live streaming  
**Auth:** Required

**Request:**
```json
{
  "user_id": "uuid",
  "started_at": "2024-12-03 16:00:00"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": {
    "LiveStreaming": {
      "id": "stream_abc123",
      "channel_name": "stream_abc123",
      "rtmp_url": "rtmp://server:1935/live/stream_abc123",
      "playback_url": "http://server:8888/live/stream_abc123/index.m3u8"
    }
  }
}
```

---

### 5. POST /api/sendGift
**Purpose:** Send virtual gift  
**Auth:** Required

**Request:**
```json
{
  "sender_id": "uuid",
  "receiver_id": "uuid",
  "live_streaming_id": "stream_abc123",
  "gift_id": "uuid",
  "gift_count": 5
}
```

**Response:**
```json
{
  "code": 200,
  "msg": {
    "User": {
      "id": "uuid",
      "wallet": 950
    },
    "Gift": {
      "id": "uuid",
      "title": "Rose",
      "coin": 10
    }
  }
}
```

---

### 6. POST /api/purchaseCoin
**Purpose:** Buy coins  
**Auth:** Required

**Request:**
```json
{
  "user_id": "uuid",
  "coin": "500",
  "title": "Popular Pack",
  "price": "4.99",
  "transaction_id": "txn_abc123",
  "device": "ios"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": {
    "User": {
      "id": "uuid",
      "wallet": 1500
    },
    "Transaction": {
      "id": 999,
      "coin": 500,
      "price": 4.99
    }
  }
}
```

---

### 7. POST /api/showCoinWorth
**Purpose:** Get coin packages  
**Auth:** None

**Request:**
```json
{}
```

**Response:**
```json
{
  "code": 200,
  "msg": {
    "CoinPackages": [
      {
        "id": 1,
        "coins": 100,
        "price": 0.99,
        "title": "Starter Pack"
      }
    ]
  }
}
```

---

### 8. POST /api/showGifts
**Purpose:** Get gifts catalog  
**Auth:** None

**Request:**
```json
{}
```

**Response:**
```json
{
  "code": 200,
  "msg": [
    {
      "Gift": {
        "id": 1,
        "title": "Rose",
        "coin": 10,
        "image": "https://..."
      }
    }
  ]
}
```

---

## 🔐 Authentication

All protected endpoints accept `auth_token` in request body:

```json
{
  "auth_token": "auth_uuid_timestamp",
  ...
}
```

---

## 📊 Response Format

All responses follow this structure:

```json
{
  "code": 200,  // 200 = success, 201 = validation error, 500 = server error
  "msg": { ... } // Data object or error message
}
```

---

## 🧪 Testing with cURL

### Register User
```bash
curl -X POST http://localhost:8080/api/registerUser \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "first_name": "Test",
    "email": "test@example.com",
    "social_id": "google_123",
    "social": "google"
  }'
```

### Get Video Feed
```bash
curl -X POST http://localhost:8080/api/showRelatedVideos \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "1",
    "device_id": "device_123",
    "starting_point": 0
  }'
```

### Start Live Stream
```bash
curl -X POST http://localhost:8080/api/liveStream \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "uuid",
    "started_at": "2024-12-03 16:00:00"
  }'
```

### Send Gift
```bash
curl -X POST http://localhost:8080/api/sendGift \
  -H "Content-Type: application/json" \
  -d '{
    "sender_id": "uuid",
    "receiver_id": "uuid",
    "gift_id": "uuid",
    "gift_count": 1
  }'
```

---

## 🎬 Live Streaming URLs

### RTMP Publish (from app)
```
rtmp://api.lomilive.com:1935/live/{streaming_id}
```

### HLS Playback (for viewers)
```
http://api.lomilive.com:8888/live/{streaming_id}/index.m3u8
```

---

## 🔢 Error Codes

| Code | Meaning |
|------|---------|
| 200  | Success |
| 201  | Validation error (insufficient coins, etc.) |
| 400  | Bad request |
| 401  | Unauthorized |
| 404  | Not found |
| 500  | Server error |

---

## 💡 Quick Tips

1. **All endpoints use POST** (even for reads)
2. **Response format is always** `{"code": 200, "msg": {...}}`
3. **Auth token** is in request body, not headers
4. **Pagination** uses `starting_point` parameter
5. **Dummy videos** are returned for immediate testing

---

## 📱 App Configuration

### Android
```java
// ApiLinks.java
public static String API_BASE_URL = "https://api.lomilive.com/api/";
```

### iOS
```swift
// ProductEndPoint.swift
var baseURL: String {
    return "https://api.lomilive.com/api/"
}
```

---

## 🚀 Next Steps

1. ✅ Test all 6 endpoints with cURL
2. ✅ Update Android/iOS app base URL
3. ✅ Build and run apps
4. ✅ Test registration flow
5. ✅ Test video feed
6. ✅ Test live streaming
7. ✅ Test gift sending
8. ✅ Test coin purchase

---

**Happy Coding!** 🎉
