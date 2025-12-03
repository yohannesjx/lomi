# TikTok Clone - Complete API Contract

**Reverse-engineered from Android (Kotlin) and iOS (Swift) source code**

---

## 📋 Table of Contents
1. [Authentication & Headers](#authentication--headers)
2. [Critical Endpoints (Implement First)](#critical-endpoints-implement-first)
3. [Complete API Reference](#complete-api-reference)
4. [Request/Response Examples](#requestresponse-examples)

---

## 🔐 Authentication & Headers

### Base URL
```
https://api.lomilive.com/api/
```

### Required Headers
All API requests must include these headers:

| Header | Value | Description |
|--------|-------|-------------|
| `Content-Type` | `application/json` | JSON content type |
| `API-KEY` | `{your_api_key}` | Static API key from Constants |
| `Auth-Token` | `{user_auth_token}` | User session token (after login) |

### Authentication Flow
1. **Social Login** → `registerUser` → Returns `User` object with `auth_token`
2. **Subsequent Requests** → Include `Auth-Token` header with the user's token
3. **Token Storage** → Apps store `auth_token` in SharedPreferences/UserDefaults

---

## 🎯 Critical Endpoints (Implement First)

These 6 endpoints are **essential** to make the apps functional:

| Priority | Endpoint | Method | Purpose |
|----------|----------|--------|---------|
| 1️⃣ | `/api/registerUser` | POST | User signup/login (social auth) |
| 2️⃣ | `/api/showUserDetail` | POST | Get user profile & wallet balance |
| 3️⃣ | `/api/showRelatedVideos` | POST | Home feed (For You page) |
| 4️⃣ | `/api/liveStream` | POST | Start live streaming session |
| 5️⃣ | `/api/sendGift` | POST | Send virtual gifts during live |
| 6️⃣ | `/api/purchaseCoin` | POST | Buy coins for gifting |

---

## 📚 Complete API Reference

### 1. Authentication & User Management

#### 1.1 Check Email Availability
**Endpoint:** `POST /api/checkEmail`

**Request:**
```json
{
  "email": "user@example.com"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "Email available"
}
```

---

#### 1.2 Verify Phone Number
**Endpoint:** `POST /api/verifyPhoneNo`

**Request:**
```json
{
  "phone": "+1234567890",
  "verify": 1,
  "code": "123456"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "Phone verified"
}
```

---

#### 1.3 Check Username Availability
**Endpoint:** `POST /api/checkUsername`

**Request:**
```json
{
  "username": "johndoe"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "Username available"
}
```

---

#### 1.4 Register User (Social Login)
**Endpoint:** `POST /api/registerUser`

**Request:**
```json
{
  "username": "johndoe",
  "dob": "1990-01-15",
  "first_name": "John",
  "last_name": "Doe",
  "email": "john@example.com",
  "phone": "+1234567890",
  "social_id": "google_12345",
  "auth_token": "temp_token",
  "device_token": "fcm_token_here",
  "social": "google"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": {
    "User": {
      "id": 123,
      "username": "johndoe",
      "first_name": "John",
      "last_name": "Doe",
      "email": "john@example.com",
      "phone": "+1234567890",
      "profile_pic": "https://...",
      "auth_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "wallet": 0,
      "verified": 0,
      "online": 1,
      "created": "2024-01-15 10:30:00"
    },
    "PushNotification": {
      "id": 1,
      "likes": 1,
      "comments": 1,
      "new_followers": 1,
      "mentions": 1,
      "direct_messages": 1,
      "video_updates": 1
    },
    "PrivacySetting": {
      "id": 1,
      "videos_download": 0,
      "direct_message": 0,
      "duet": 1,
      "liked_videos": 0,
      "video_comment": 1,
      "order_history": 0
    }
  }
}
```

**Triggered by:** Signup screen, social login buttons

---

#### 1.5 Show User Detail
**Endpoint:** `POST /api/showUserDetail`

**Request (Own Profile):**
```json
{
  "auth_token": "user_auth_token"
}
```

**Request (Other User):**
```json
{
  "user_id": "123",
  "other_user_id": "456"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": {
    "User": {
      "id": 123,
      "username": "johndoe",
      "first_name": "John",
      "wallet": 1500,
      "total_all_time_coins": 5000,
      "verified": 1,
      "button": "following"
    }
  }
}
```

**Triggered by:** Profile screen, app launch

---

#### 1.6 Edit Profile
**Endpoint:** `POST /api/editProfile`

**Request:**
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "username": "johndoe",
  "bio": "Content creator",
  "website": "https://johndoe.com",
  "email": "john@example.com",
  "phone": "+1234567890",
  "gender": "male",
  "profile_pic": "base64_image_data"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": {
    "User": { /* updated user object */ }
  }
}
```

**Triggered by:** Edit profile screen

---

#### 1.7 Delete User Account
**Endpoint:** `POST /api/deleteUserAccount`

**Request:**
```json
{
  "auth_token": "user_auth_token"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "Account deleted successfully"
}
```

**Triggered by:** Settings → Delete Account

---

### 2. Video Feed & Discovery

#### 2.1 Show Related Videos (For You Feed)
**Endpoint:** `POST /api/showRelatedVideos`

**Request:**
```json
{
  "user_id": "123",
  "device_id": "device_uuid",
  "starting_point": 0,
  "lat": 37.7749,
  "long": -122.4194,
  "tag_product": 0,
  "delivery_address_id": 0
}
```

**Response:**
```json
{
  "code": 200,
  "msg": [
    {
      "Video": {
        "id": 789,
        "user_id": 456,
        "description": "Amazing dance video!",
        "video": "https://cdn.example.com/videos/789.mp4",
        "thumbnail": "https://cdn.example.com/thumbs/789.jpg",
        "sound_id": 12,
        "view": 15000,
        "like": 1200,
        "comment_count": 45,
        "share": 30,
        "privacy_type": "public",
        "allow_comments": 1,
        "allow_duet": 1,
        "created": "2024-01-15 10:30:00"
      },
      "User": {
        "id": 456,
        "username": "dancer123",
        "profile_pic": "https://...",
        "verified": 1
      },
      "Sound": {
        "id": 12,
        "title": "Trending Beat",
        "sound": "https://cdn.example.com/sounds/12.mp3"
      }
    }
  ]
}
```

**Triggered by:** Home screen (For You tab)

---

#### 2.2 Show Following Videos
**Endpoint:** `POST /api/showFollowingVideos`

**Request:**
```json
{
  "user_id": "123",
  "device_id": "device_uuid",
  "starting_point": 0
}
```

**Response:** Same structure as `showRelatedVideos`

**Triggered by:** Home screen (Following tab)

---

#### 2.3 Show Nearby Videos
**Endpoint:** `POST /api/showNearbyVideos`

**Request:**
```json
{
  "user_id": "123",
  "device_id": "device_uuid",
  "lat": 37.7749,
  "long": -122.4194,
  "starting_point": 0
}
```

**Response:** Same structure as `showRelatedVideos`

**Triggered by:** Nearby videos screen

---

#### 2.4 Show Discovery Sections
**Endpoint:** `POST /api/showDiscoverySections`

**Request:**
```json
{
  "user_id": "123",
  "country_id": "1",
  "starting_point": 0
}
```

**Response:**
```json
{
  "code": 200,
  "msg": [
    {
      "DiscoverySection": {
        "id": 1,
        "title": "Trending Hashtags",
        "type": "hashtag"
      },
      "Videos": [ /* array of video objects */ ]
    }
  ]
}
```

**Triggered by:** Discover/Explore screen

---

#### 2.5 Search
**Endpoint:** `POST /api/search`

**Request:**
```json
{
  "user_id": "123",
  "type": "users",
  "keyword": "dance",
  "starting_point": 0
}
```

**Types:** `users`, `videos`, `hashtags`, `sounds`

**Response:**
```json
{
  "code": 200,
  "msg": {
    "Users": [ /* array of user objects */ ],
    "Videos": [ /* array of video objects */ ],
    "Hashtags": [ /* array of hashtag objects */ ],
    "Sounds": [ /* array of sound objects */ ]
  }
}
```

**Triggered by:** Search screen

---

### 3. Video Interactions

#### 3.1 Like Video
**Endpoint:** `POST /api/likeVideo`

**Request:**
```json
{
  "video_id": "789"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": {
    "Video": {
      "id": 789,
      "like": 1201,
      "liked": "1"
    }
  }
}
```

**Triggered by:** Double-tap or like button on video

---

#### 3.2 Add Video to Favorites
**Endpoint:** `POST /api/addVideoFavourite`

**Request:**
```json
{
  "video_id": "789"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "Video added to favorites"
}
```

**Triggered by:** Bookmark/favorite button

---

#### 3.3 Show Video Detail
**Endpoint:** `POST /api/showVideoDetail`

**Request:**
```json
{
  "video_id": "789"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": {
    "Video": { /* full video object */ },
    "User": { /* video owner */ },
    "Sound": { /* sound details */ }
  }
}
```

**Triggered by:** Video detail page

---

#### 3.4 Delete Video
**Endpoint:** `POST /api/deleteVideo`

**Request:**
```json
{
  "video_id": "789"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "Video deleted successfully"
}
```

**Triggered by:** Video options → Delete

---

#### 3.5 Repost Video
**Endpoint:** `POST /api/repostVideo`

**Request:**
```json
{
  "video_id": "789"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "Video reposted"
}
```

**Triggered by:** Share → Repost

---

#### 3.6 Share Video
**Endpoint:** `POST /api/shareVideo`

**Request:**
```json
{
  "video_id": "789"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "Share count updated"
}
```

**Triggered by:** Share button

---

#### 3.7 Report Video
**Endpoint:** `POST /api/reportVideo`

**Request:**
```json
{
  "video_id": "789",
  "reason_id": "3",
  "description": "Inappropriate content"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "Report submitted"
}
```

**Triggered by:** Video options → Report

---

### 4. Video Upload

#### 4.1 Post Video
**Endpoint:** `POST /api/postVideo` (Multipart)

**Request (Multipart Form Data):**
```
file: [video file]
privacy_type: "public"
user_id: "123"
sound_id: "12"
allow_comments: "1"
description: "Check out my new dance!"
allow_duet: "1"
users_json: [{"user_id": "456"}]
hashtags_json: [{"hashtag": "dance"}]
story: "0"
video_id: ""
location_string: "San Francisco, CA"
lat: "37.7749"
long: "-122.4194"
google_place_id: "ChIJIQBpAG2ahYAR_6128GcTUEo"
location_name: "San Francisco"
width: "1080"
height: "1920"
products: []
user_thumbnail: "base64_thumb"
default_thumbnail: "0"
user_selected_thum: "1"
duet: "0"
```

**Response:**
```json
{
  "code": 200,
  "msg": {
    "Video": {
      "id": 790,
      "user_id": 123,
      "video": "https://cdn.example.com/videos/790.mp4",
      "thumbnail": "https://cdn.example.com/thumbs/790.jpg"
    }
  }
}
```

**Triggered by:** Post video screen

---

#### 4.2 Edit Video
**Endpoint:** `POST /api/editVideo` (Multipart)

**Request:**
```
privacy_type: "public"
user_id: "123"
allow_comments: "1"
description: "Updated description"
allow_duet: "1"
users_json: [{"user_id": "456"}]
hashtags_json: [{"hashtag": "dance"}]
video_id: "790"
location_string: "San Francisco, CA"
lat: "37.7749"
long: "-122.4194"
google_place_id: "ChIJIQBpAG2ahYAR_6128GcTUEo"
location_name: "San Francisco"
products: []
tag_store_id: ""
```

**Response:**
```json
{
  "code": 200,
  "msg": "Video updated successfully"
}
```

**Triggered by:** Edit video screen

---

### 5. Comments

#### 5.1 Show Video Comments
**Endpoint:** `POST /api/showVideoComments`

**Request:**
```json
{
  "video_id": "789",
  "starting_point": 0
}
```

**Response:**
```json
{
  "code": 200,
  "msg": [
    {
      "VideoComment": {
        "id": 101,
        "video_id": 789,
        "user_id": 456,
        "comment": "Amazing video!",
        "like_count": 15,
        "like": "0",
        "owner_like": "0",
        "parent_id": 0,
        "created": "2024-01-15 11:00:00"
      },
      "User": {
        "id": 456,
        "username": "commenter",
        "profile_pic": "https://..."
      },
      "Replies": [ /* array of reply comments */ ]
    }
  ]
}
```

**Triggered by:** Comments screen

---

#### 5.2 Post Comment on Video
**Endpoint:** `POST /api/postCommentOnVideo`

**Request:**
```json
{
  "video_id": "789",
  "comment": "Great video!",
  "users_json": [
    {"user_id": "456"}
  ]
}
```

**Response:**
```json
{
  "code": 200,
  "msg": {
    "VideoComment": {
      "id": 102,
      "video_id": 789,
      "comment": "Great video!",
      "created": "2024-01-15 11:05:00"
    },
    "User": { /* commenter details */ },
    "Video": { /* video details */ }
  }
}
```

**Triggered by:** Comment input field

---

#### 5.3 Post Comment Reply
**Endpoint:** `POST /api/postCommentOnVideo`

**Request:**
```json
{
  "parent_id": "101",
  "user_id": "123",
  "comment": "Thanks!",
  "video_id": "789",
  "users_json": [
    {"user_id": "456"}
  ]
}
```

**Response:** Same as post comment

**Triggered by:** Reply to comment

---

#### 5.4 Like Comment
**Endpoint:** `POST /api/likeComment`

**Request:**
```json
{
  "comment_id": "101"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": {
    "VideoComment": {
      "id": 101,
      "like_count": 16,
      "like": "1"
    }
  }
}
```

**Triggered by:** Like button on comment

---

#### 5.5 Like Comment Reply
**Endpoint:** `POST /api/likeCommentReply`

**Request:**
```json
{
  "comment_reply_id": "103"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "Reply liked"
}
```

**Triggered by:** Like button on reply

---

#### 5.6 Delete Video Comment
**Endpoint:** `POST /api/deleteVideoComment`

**Request:**
```json
{
  "comment_id": "101"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "Comment deleted"
}
```

**Triggered by:** Delete comment option

---

#### 5.7 Delete Comment Reply
**Endpoint:** `POST /api/deleteVideoCommentReply`

**Request:**
```json
{
  "comment_reply_id": "103"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "Reply deleted"
}
```

**Triggered by:** Delete reply option

---

#### 5.8 Pin Comment
**Endpoint:** `POST /api/pinComment`

**Request:**
```json
{
  "comment_id": "101",
  "video_id": "789"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "Comment pinned"
}
```

**Triggered by:** Pin comment option (video owner)

---

### 6. Social Features

#### 6.1 Follow User
**Endpoint:** `POST /api/followUser`

**Request:**
```json
{
  "sender_id": "123",
  "receiver_id": "456"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": {
    "User": {
      "id": 456,
      "button": "following"
    }
  }
}
```

**Triggered by:** Follow/Unfollow button

---

#### 6.2 Show Followers
**Endpoint:** `POST /api/showFollowers`

**Request:**
```json
{
  "user_id": "123",
  "other_user_id": "456",
  "starting_point": 0
}
```

**Response:**
```json
{
  "code": 200,
  "msg": [
    {
      "User": {
        "id": 789,
        "username": "follower1",
        "profile_pic": "https://...",
        "button": "follow back"
      }
    }
  ]
}
```

**Triggered by:** Followers list

---

#### 6.3 Show Following
**Endpoint:** `POST /api/showFollowing`

**Request:**
```json
{
  "user_id": "123",
  "other_user_id": "456",
  "starting_point": 0
}
```

**Response:** Same structure as followers

**Triggered by:** Following list

---

#### 6.4 Show Suggested Users
**Endpoint:** `POST /api/showSuggestedUsers`

**Request:**
```json
{
  "user_id": "123",
  "other_user_id": "",
  "starting_point": 0
}
```

**Response:**
```json
{
  "code": 200,
  "msg": [
    {
      "User": {
        "id": 999,
        "username": "suggested_user",
        "profile_pic": "https://...",
        "button": "follow"
      }
    }
  ]
}
```

**Triggered by:** Suggested users section

---

#### 6.5 Show Profile Visitors
**Endpoint:** `POST /api/showProfileVisitors`

**Request:**
```json
{
  "user_id": "123",
  "starting_point": 0
}
```

**Response:**
```json
{
  "code": 200,
  "msg": [
    {
      "User": { /* visitor details */ },
      "visited_at": "2024-01-15 10:00:00"
    }
  ]
}
```

**Triggered by:** Profile visitors screen

---

#### 6.6 Block User
**Endpoint:** `POST /api/blockUser`

**Request:**
```json
{
  "user_id": "123",
  "block_user_id": "456"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "User blocked"
}
```

**Triggered by:** Block user option

---

#### 6.7 Show Blocked Users
**Endpoint:** `POST /api/showBlockedUsers`

**Request:**
```json
{
  "user_id": "123",
  "starting_point": 0
}
```

**Response:**
```json
{
  "code": 200,
  "msg": [
    {
      "User": {
        "id": 456,
        "username": "blocked_user",
        "profile_pic": "https://..."
      }
    }
  ]
}
```

**Triggered by:** Blocked users list

---

### 7. Live Streaming

#### 7.1 Start Live Stream
**Endpoint:** `POST /api/liveStream`

**Request:**
```json
{
  "user_id": "123",
  "started_at": "2024-01-15 12:00:00"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": {
    "LiveStreaming": {
      "id": "stream_123",
      "user_id": 123,
      "channel_name": "stream_123",
      "started_at": "2024-01-15 12:00:00",
      "status": "live"
    }
  }
}
```

**Triggered by:** Go Live button

---

#### 7.2 Show Gifts
**Endpoint:** `POST /api/showGifts`

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
        "image": "https://cdn.example.com/gifts/rose.png",
        "coin": 10,
        "type": "normal"
      }
    },
    {
      "Gift": {
        "id": 2,
        "title": "Diamond",
        "image": "https://cdn.example.com/gifts/diamond.png",
        "coin": 1000,
        "type": "premium"
      }
    }
  ]
}
```

**Triggered by:** Gift panel in live stream

---

#### 7.3 Send Gift
**Endpoint:** `POST /api/sendGift`

**Request:**
```json
{
  "sender_id": "123",
  "receiver_id": "456",
  "video_id": "",
  "live_streaming_id": "stream_123",
  "gift_id": "1",
  "gift_count": "5"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": {
    "User": {
      "id": 123,
      "wallet": 950
    },
    "Gift": {
      "id": 1,
      "title": "Rose",
      "coin": 10
    }
  }
}
```

**Triggered by:** Send gift button in live stream

---

#### 7.4 Show Sent Gifts Against Video
**Endpoint:** `POST /api/showSentGiftsAgainstVideo`

**Request:**
```json
{
  "video_id": "789"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": [
    {
      "SentGift": {
        "id": 1,
        "sender_id": 123,
        "gift_id": 1,
        "gift_count": 5,
        "created": "2024-01-15 12:05:00"
      },
      "User": { /* sender details */ },
      "Gift": { /* gift details */ }
    }
  ]
}
```

**Triggered by:** Gift history on video

---

### 8. Coins & Wallet

#### 8.1 Show Coin Worth
**Endpoint:** `POST /api/showCoinWorth`

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
      },
      {
        "id": 2,
        "coins": 500,
        "price": 4.99,
        "title": "Popular Pack"
      },
      {
        "id": 3,
        "coins": 1000,
        "price": 9.99,
        "title": "Best Value"
      }
    ]
  }
}
```

**Triggered by:** Coin purchase screen

---

#### 8.2 Purchase Coin
**Endpoint:** `POST /api/purchaseCoin`

**Request:**
```json
{
  "user_id": "123",
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
      "id": 123,
      "wallet": 1500,
      "total_all_time_coins": 5500
    },
    "Transaction": {
      "id": 999,
      "coin": 500,
      "price": 4.99,
      "created": "2024-01-15 12:10:00"
    }
  }
}
```

**Triggered by:** In-app purchase completion

---

#### 8.3 Purchase Coins from Stripe
**Endpoint:** `POST /api/purchaseCoinsFromStripe`

**Request:**
```json
{
  "user_id": "123",
  "coin": "500",
  "title": "Popular Pack",
  "price": "4.99",
  "stripe_token": "tok_visa",
  "device": "android"
}
```

**Response:** Same as purchaseCoin

**Triggered by:** Stripe payment completion

---

#### 8.4 Withdraw Request
**Endpoint:** `POST /api/withdrawRequest`

**Request:**
```json
{
  "user_id": "123",
  "coin": "1000",
  "amount": "10.00",
  "email": "john@example.com"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": {
    "WithdrawalRequest": {
      "id": 50,
      "user_id": 123,
      "coin": 1000,
      "amount": 10.00,
      "status": "pending",
      "created": "2024-01-15 12:15:00"
    },
    "User": {
      "wallet": 500
    }
  }
}
```

**Triggered by:** Withdraw earnings screen

---

#### 8.5 Show Withdrawal History
**Endpoint:** `POST /api/showWithdrawalHistory`

**Request:**
```json
{
  "user_id": "123"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": [
    {
      "WithdrawalRequest": {
        "id": 50,
        "coin": 1000,
        "amount": 10.00,
        "status": "completed",
        "created": "2024-01-15 12:15:00"
      }
    }
  ]
}
```

**Triggered by:** Withdrawal history screen

---

#### 8.6 Add Payout Method
**Endpoint:** `POST /api/addPayout`

**Request:**
```json
{
  "user_id": "123",
  "value": "john@paypal.com",
  "type": "paypal"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "Payout method added"
}
```

**Triggered by:** Add payout method screen

---

#### 8.7 Show Payout Methods
**Endpoint:** `POST /api/showPayout`

**Request:**
```json
{
  "user_id": "123"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": [
    {
      "PayoutMethod": {
        "id": 1,
        "type": "paypal",
        "value": "john@paypal.com"
      }
    }
  ]
}
```

**Triggered by:** Payout methods screen

---

### 9. Notifications

#### 9.1 Show All Notifications
**Endpoint:** `POST /api/showAllNotifications`

**Request:**
```json
{
  "user_id": "123",
  "starting_point": 0
}
```

**Response:**
```json
{
  "code": 200,
  "msg": [
    {
      "Notification": {
        "id": 1,
        "sender_id": 456,
        "receiver_id": 123,
        "type": "like",
        "video_id": 789,
        "live_streaming_id": null,
        "room_id": null,
        "message": "liked your video",
        "read": 0,
        "created": "2024-01-15 11:00:00"
      },
      "User": {
        "id": 456,
        "username": "liker",
        "profile_pic": "https://..."
      }
    }
  ]
}
```

**Triggered by:** Notifications screen

---

#### 9.2 Show Unread Notifications
**Endpoint:** `POST /api/showUnReadNotifications`

**Request:**
```json
{
  "user_id": "123"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": {
    "unread_count": 5
  }
}
```

**Triggered by:** Notification badge

---

#### 9.3 Read Notification
**Endpoint:** `POST /api/readNotification`

**Request:**
```json
{
  "user_id": "123",
  "notification_id": "1"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "Notification marked as read"
}
```

**Triggered by:** Opening notification

---

#### 9.4 Send Notification
**Endpoint:** `POST /api/sendNotification`

**Request:**
```json
{
  "sender_id": "123",
  "receiver_id": "456",
  "title": "New Follower",
  "message": "johndoe started following you"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "Notification sent"
}
```

**Triggered by:** System events (follow, like, comment)

---

### 10. Sounds & Music

#### 10.1 Show Sounds
**Endpoint:** `POST /api/showSounds`

**Request:**
```json
{
  "user_id": "123",
  "starting_point": 0
}
```

**Response:**
```json
{
  "code": 200,
  "msg": [
    {
      "Sound": {
        "id": 12,
        "title": "Trending Beat",
        "sound": "https://cdn.example.com/sounds/12.mp3",
        "duration": 30,
        "artist": "DJ Cool",
        "cover": "https://cdn.example.com/covers/12.jpg"
      }
    }
  ]
}
```

**Triggered by:** Sounds library

---

#### 10.2 Show Favorite Sounds
**Endpoint:** `POST /api/showFavouriteSounds`

**Request:**
```json
{
  "user_id": "123",
  "starting_point": 0
}
```

**Response:** Same structure as showSounds

**Triggered by:** Favorite sounds list

---

#### 10.3 Add Sound to Favorites
**Endpoint:** `POST /api/addSoundFavourite`

**Request:**
```json
{
  "sound_id": "12"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "Sound added to favorites"
}
```

**Triggered by:** Favorite button on sound

---

#### 10.4 Show Videos Against Sound
**Endpoint:** `POST /api/showVideosAgainstSound`

**Request:**
```json
{
  "user_id": "123",
  "sound_id": "12",
  "starting_point": 0
}
```

**Response:**
```json
{
  "code": 200,
  "msg": [
    {
      "Video": { /* video using this sound */ },
      "User": { /* video creator */ }
    }
  ]
}
```

**Triggered by:** Sound detail page

---

### 11. Hashtags

#### 11.1 Show Videos Against Hashtag
**Endpoint:** `POST /api/showVideosAgainstHashtag`

**Request:**
```json
{
  "user_id": "123",
  "hashtag": "dance",
  "starting_point": 0
}
```

**Response:**
```json
{
  "code": 200,
  "msg": [
    {
      "Video": { /* video with this hashtag */ },
      "User": { /* video creator */ }
    }
  ]
}
```

**Triggered by:** Hashtag page

---

#### 11.2 Add Hashtag to Favorites
**Endpoint:** `POST /api/addHashtagFavourite`

**Request:**
```json
{
  "hashtag": "dance"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "Hashtag added to favorites"
}
```

**Triggered by:** Follow hashtag button

---

#### 11.3 Show Favorite Hashtags
**Endpoint:** `POST /api/showFavouriteHashtags`

**Request:**
```json
{
  "user_id": "123"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": [
    {
      "Hashtag": {
        "id": 1,
        "hashtag": "dance",
        "video_count": 15000
      }
    }
  ]
}
```

**Triggered by:** Favorite hashtags list

---

### 12. User Videos

#### 12.1 Show Videos Against User ID
**Endpoint:** `POST /api/showVideosAgainstUserID`

**Request:**
```json
{
  "user_id": "123",
  "other_user_id": "456",
  "starting_point": 0
}
```

**Response:**
```json
{
  "code": 200,
  "msg": [
    {
      "Video": { /* user's video */ }
    }
  ]
}
```

**Triggered by:** User profile (Videos tab)

---

#### 12.2 Show User Liked Videos
**Endpoint:** `POST /api/showUserLikedVideos`

**Request:**
```json
{
  "user_id": "123",
  "other_user_id": "456",
  "starting_point": 0
}
```

**Response:** Same structure as showVideosAgainstUserID

**Triggered by:** User profile (Liked tab)

---

#### 12.3 Show User Reposted Videos
**Endpoint:** `POST /api/showUserRepostedVideos`

**Request:**
```json
{
  "user_id": "123",
  "other_user_id": "456",
  "starting_point": 0
}
```

**Response:** Same structure as showVideosAgainstUserID

**Triggered by:** User profile (Reposts tab)

---

#### 12.4 Show Favorite Videos
**Endpoint:** `POST /api/showFavouriteVideos`

**Request:**
```json
{
  "user_id": "123",
  "starting_point": 0
}
```

**Response:** Same structure as showVideosAgainstUserID

**Triggered by:** Favorites screen

---

### 13. Location

#### 13.1 Show Videos Against Location
**Endpoint:** `POST /api/showVideosAgainstLocation`

**Request:**
```json
{
  "user_id": "123",
  "location_string": "San Francisco, CA",
  "starting_point": 0
}
```

**Response:**
```json
{
  "code": 200,
  "msg": [
    {
      "Video": { /* video from this location */ },
      "User": { /* video creator */ }
    }
  ]
}
```

**Triggered by:** Location page

---

### 14. Settings & Privacy

#### 14.1 Show Settings
**Endpoint:** `POST /api/showSettings`

**Request:**
```json
{
  "user_id": "123"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": {
    "PushNotification": {
      "likes": 1,
      "comments": 1,
      "new_followers": 1,
      "mentions": 1,
      "direct_messages": 1,
      "video_updates": 1
    },
    "PrivacySetting": {
      "videos_download": 0,
      "direct_message": 0,
      "duet": 1,
      "liked_videos": 0,
      "video_comment": 1,
      "order_history": 0
    }
  }
}
```

**Triggered by:** Settings screen

---

#### 14.2 Add Privacy Setting
**Endpoint:** `POST /api/addPrivacySetting`

**Request:**
```json
{
  "user_id": "123",
  "videos_download": 0,
  "direct_message": 0,
  "duet": 1,
  "liked_videos": 0,
  "video_comment": 1,
  "order_history": 0
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "Privacy settings updated"
}
```

**Triggered by:** Privacy settings screen

---

#### 14.3 Update Push Notification Settings
**Endpoint:** `POST /api/updatePushNotificationSettings`

**Request:**
```json
{
  "user_id": "123",
  "likes": 1,
  "comments": 1,
  "new_followers": 1,
  "mentions": 1,
  "direct_messages": 1,
  "video_updates": 1
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "Notification settings updated"
}
```

**Triggered by:** Notification settings screen

---

#### 14.4 User Verification Request
**Endpoint:** `POST /api/userVerificationRequest`

**Request:**
```json
{
  "user_id": "123",
  "name": "John Doe",
  "attachment": {
    "file_data": "base64_encoded_document"
  }
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "Verification request submitted"
}
```

**Triggered by:** Request verification screen

---

### 15. Reporting

#### 15.1 Show Report Reasons
**Endpoint:** `POST /api/showReportReasons`

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
      "ReportReason": {
        "id": 1,
        "reason": "Spam"
      }
    },
    {
      "ReportReason": {
        "id": 2,
        "reason": "Inappropriate content"
      }
    },
    {
      "ReportReason": {
        "id": 3,
        "reason": "Harassment"
      }
    }
  ]
}
```

**Triggered by:** Report screen

---

#### 15.2 Report User
**Endpoint:** `POST /api/reportUser`

**Request:**
```json
{
  "user_id": "123",
  "reported_user_id": "456",
  "reason_id": "3",
  "description": "User is harassing me"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "User reported successfully"
}
```

**Triggered by:** Report user option

---

### 16. Device & Analytics

#### 16.1 Add Device Data
**Endpoint:** `POST /api/addDeviceData`

**Request:**
```json
{
  "user_id": "123",
  "device": "android",
  "lat": "37.7749",
  "long": "-122.4194",
  "version": "1.0.0",
  "ip": "192.168.1.1",
  "device_token": "fcm_token_here"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "Device data updated"
}
```

**Triggered by:** App launch, location update

---

#### 16.2 Register Device
**Endpoint:** `POST /api/registerDevice`

**Request:**
```json
{
  "device_token": "fcm_token_here",
  "device_type": "ios"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "Device registered"
}
```

**Triggered by:** App installation, FCM token refresh

---

#### 16.3 Watch Video (Analytics)
**Endpoint:** `POST /api/watchVideo`

**Request:**
```json
{
  "user_id": "123",
  "watch_videos": [
    {
      "video_id": "789",
      "duration": 85
    },
    {
      "video_id": "790",
      "duration": 100
    }
  ]
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "Watch data recorded"
}
```

**Triggered by:** Background (batched every 3 videos)

---

### 17. Additional Endpoints

#### 17.1 Show App Slider
**Endpoint:** `POST /api/showAppSlider`

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
      "Slider": {
        "id": 1,
        "image": "https://cdn.example.com/sliders/1.jpg",
        "link": "https://example.com/promo"
      }
    }
  ]
}
```

**Triggered by:** Home screen banners

---

#### 17.2 Show Stickers
**Endpoint:** `POST /api/showStickers`

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
      "Sticker": {
        "id": 1,
        "image": "https://cdn.example.com/stickers/1.png",
        "category": "emoji"
      }
    }
  ]
}
```

**Triggered by:** Sticker picker in live stream

---

#### 17.3 Show Interest Sections
**Endpoint:** `POST /api/showInterestSection`

**Request:**
```json
{
  "auth_token": "user_auth_token"
}
```

**Response:**
```json
{
  "code": 200,
  "msg": [
    {
      "InterestSection": {
        "id": 1,
        "title": "Entertainment"
      },
      "Interest": [
        {
          "id": 1,
          "title": "Music",
          "icon": "https://..."
        },
        {
          "id": 2,
          "title": "Dance",
          "icon": "https://..."
        }
      ]
    }
  ]
}
```

**Triggered by:** Onboarding (interest selection)

---

#### 17.4 Add User Interest
**Endpoint:** `POST /api/addUserInterest`

**Request:**
```json
{
  "user_id": "123",
  "interests": [1, 2, 5, 8]
}
```

**Response:**
```json
{
  "code": 200,
  "msg": "Interests saved"
}
```

**Triggered by:** Interest selection screen

---

## 📝 Request/Response Examples

### Example 1: Complete Login Flow

**Step 1: Check Email**
```bash
POST /api/checkEmail
{
  "email": "john@example.com"
}
→ Response: {"code": 200, "msg": "Email available"}
```

**Step 2: Verify Phone**
```bash
POST /api/verifyPhoneNo
{
  "phone": "+1234567890",
  "verify": 1,
  "code": "123456"
}
→ Response: {"code": 200, "msg": "Phone verified"}
```

**Step 3: Check Username**
```bash
POST /api/checkUsername
{
  "username": "johndoe"
}
→ Response: {"code": 200, "msg": "Username available"}
```

**Step 4: Register User**
```bash
POST /api/registerUser
{
  "username": "johndoe",
  "first_name": "John",
  "last_name": "Doe",
  "email": "john@example.com",
  "phone": "+1234567890",
  "social_id": "google_12345",
  "social": "google",
  "device_token": "fcm_token",
  "auth_token": "temp"
}
→ Response: {
  "code": 200,
  "msg": {
    "User": {
      "id": 123,
      "auth_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "wallet": 0
    }
  }
}
```

**Step 5: Store auth_token and use in all subsequent requests**

---

### Example 2: Send Gift During Live Stream

**Step 1: Start Live Stream**
```bash
POST /api/liveStream
{
  "user_id": "123",
  "started_at": "2024-01-15 12:00:00"
}
→ Response: {
  "code": 200,
  "msg": {
    "LiveStreaming": {
      "id": "stream_123"
    }
  }
}
```

**Step 2: Viewer Sends Gift**
```bash
POST /api/sendGift
{
  "sender_id": "456",
  "receiver_id": "123",
  "live_streaming_id": "stream_123",
  "gift_id": "1",
  "gift_count": "5"
}
→ Response: {
  "code": 200,
  "msg": {
    "User": {
      "wallet": 950
    }
  }
}
```

---

### Example 3: Post a Video

**Step 1: Upload Video**
```bash
POST /api/postVideo (multipart/form-data)
file: [video.mp4]
user_id: "123"
description: "My awesome dance!"
privacy_type: "public"
allow_comments: "1"
allow_duet: "1"
hashtags_json: [{"hashtag": "dance"}, {"hashtag": "viral"}]
lat: "37.7749"
long: "-122.4194"
→ Response: {
  "code": 200,
  "msg": {
    "Video": {
      "id": 790,
      "video": "https://cdn.example.com/videos/790.mp4"
    }
  }
}
```

---

## 🎯 Summary

### Total Endpoints: **100+**

### Categories:
- **Authentication:** 7 endpoints
- **Video Feed:** 5 endpoints
- **Video Interactions:** 7 endpoints
- **Video Upload:** 2 endpoints
- **Comments:** 8 endpoints
- **Social:** 7 endpoints
- **Live Streaming:** 4 endpoints
- **Coins & Wallet:** 7 endpoints
- **Notifications:** 4 endpoints
- **Sounds:** 4 endpoints
- **Hashtags:** 3 endpoints
- **User Videos:** 4 endpoints
- **Settings:** 4 endpoints
- **Reporting:** 2 endpoints
- **Device & Analytics:** 3 endpoints
- **Miscellaneous:** 4 endpoints

### Priority Implementation Order:
1. ✅ `registerUser` - User authentication
2. ✅ `showUserDetail` - Profile & wallet
3. ✅ `showRelatedVideos` - Home feed
4. ✅ `liveStream` - Start streaming
5. ✅ `sendGift` - Virtual gifting
6. ✅ `purchaseCoin` - Coin purchases

---

**Generated:** 2024-12-03  
**Source:** Android (Kotlin) + iOS (Swift) reverse engineering  
**Backend Target:** Go (lomi_mini)
