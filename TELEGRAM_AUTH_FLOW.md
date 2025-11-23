# Complete Telegram Mini App Authentication Flow

## Overview

This document describes the complete, secure, passwordless authentication flow for the Lomi Social Telegram Mini App.

## Architecture

### Frontend Flow
1. **AuthGuard Component** - Automatically authenticates on app load
2. **WelcomeScreen** - Fallback for manual login (if needed)
3. **Onboarding Flow** - For new users
4. **Main App** - For authenticated users who completed onboarding

### Backend Flow
1. **Receive initData** - Via `Authorization: tma <initData>` header
2. **Validate initData** - HMAC-SHA256 hash verification using bot token
3. **Extract User Data** - Parse Telegram user information
4. **Find/Create User** - Database lookup by Telegram ID
5. **Generate JWT Tokens** - Access token + Refresh token
6. **Return Response** - Tokens + user profile data

---

## Step-by-Step Flow

### 1. App Opens in Telegram

When a user opens the Mini App from Telegram:

```
User clicks Mini App button in Telegram
    ↓
Telegram injects initData into the page
    ↓
AuthGuard component detects Telegram environment
    ↓
Automatically extracts initData
```

**initData Format:**
```
user=%7B%22id%22%3A123456789%2C%22first_name%22%3A%22John%22%7D&auth_date=1234567890&hash=abc123...
```

### 2. Automatic Authentication

**Frontend (AuthGuard.tsx):**
```typescript
1. Check if app is in Telegram: isTelegramWebApp()
2. Wait for initData (with retries): waitForInitData()
3. Send to backend: POST /api/v1/auth/telegram
   Headers: { Authorization: "tma <initData>" }
4. Store tokens and user data
5. Route based on onboarding status
```

**Backend (auth.go):**
```go
1. Extract initData from Authorization header
2. Validate using initdata.Validate(initData, botToken, time.Hour)
3. Parse user data: initdata.Parse(initData)
4. Find or create user in database
5. Generate JWT tokens
6. Return: { access_token, refresh_token, user }
```

### 3. Backend Validation (Security Critical)

The backend **MUST** validate initData before trusting any user data:

```go
// 1. Validate hash using HMAC-SHA256
err := initdata.Validate(initData, botToken, time.Hour)
if err != nil {
    return 401, "Invalid Telegram data"
}

// 2. Parse user data
parsedData, err := initdata.Parse(initData)
if err != nil {
    return 401, "Failed to parse Telegram data"
}

// 3. Extract user ID (this is the only trusted identifier)
userID := parsedData.User.ID
```

**Security Requirements:**
- ✅ Never trust client-side data without server validation
- ✅ Always verify HMAC-SHA256 hash
- ✅ Check expiration (auth_date must be within 1 hour)
- ✅ Use Telegram user ID as unique identifier
- ✅ Generate secure JWT tokens server-side

### 4. User Creation/Retrieval

**New User:**
```go
user = models.User{
    TelegramID: tgUser.ID,
    Name: firstName,
    Age: 18,                    // Default, updated in onboarding
    Gender: GenderOther,        // Default, updated in onboarding
    City: "Not Set",            // Default, updated in onboarding
    Religion: ReligionNone,     // Default
    RelationshipGoal: GoalDating,
    VerificationStatus: VerificationPending,
    IsActive: true,
}
database.DB.Create(&user)
```

**Existing User:**
```go
database.DB.Where("telegram_id = ?", tgUser.ID).First(&user)
// Update Telegram info if changed
database.DB.Model(&user).Updates(updates)
```

### 5. Token Generation

```go
tokens, err := utils.CreateToken(user.ID, jwtSecret)
// Returns:
// - access_token (24h expiry)
// - refresh_token (7 days expiry)
```

### 6. Response Format

**Success (200):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "uuid-here",
    "name": "John",
    "is_verified": false,
    "has_profile": false  // false = needs onboarding
  }
}
```

**Error (401):**
```json
{
  "error": "Invalid Telegram data",
  "details": "sign is missing"
}
```

### 7. Frontend Token Storage

```typescript
// Store securely
await storage.setItem('lomi_access_token', accessToken);
await storage.setItem('lomi_refresh_token', refreshToken);
await storage.setItem('lomi_user', JSON.stringify(user));

// Update auth store
set({
  accessToken,
  refreshToken,
  user,
  isAuthenticated: true,
});
```

### 8. Routing Logic

**After Authentication:**
```typescript
if (user.has_profile === true) {
    // User completed onboarding → Main App
    navigation.reset({ routes: [{ name: 'Main' }] });
} else {
    // New user → Onboarding
    navigation.reset({ routes: [{ name: 'ProfileSetup' }] });
}
```

### 9. Onboarding Flow

**New users go through:**
1. ProfileSetup → Name, Age, Gender, City
2. PhotoUpload → Profile photos
3. Interests → Select interests
4. GenderPreference → Who they're looking for
5. Main App → Start matching

**After onboarding:**
- `has_profile` is set to `true`
- User can access main app features

### 10. Subsequent Requests

**All API requests include token:**
```typescript
api.interceptors.request.use(async (config) => {
    const token = await storage.getItem('lomi_access_token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});
```

---

## Error Handling

### Not in Telegram
**Frontend shows:**
```
📱 Please open this app from inside Telegram
[Open in Telegram Button]
```

### initData Missing
**Frontend:**
- Retries up to 10 times with increasing delays
- Shows error if still not found

**Backend:**
- Returns 401: "Missing Authorization header"

### Invalid initData
**Backend:**
- Returns 401: "Invalid Telegram data" + details
- Logs error for debugging

### Expired initData
**Backend:**
- Returns 401: "init data is expired"
- User must reopen app from Telegram

### Database Error
**Backend:**
- Returns 500 with detailed error
- Logs full error for debugging

---

## Security Checklist

- ✅ **Server-side validation** - All initData validated on backend
- ✅ **HMAC-SHA256 verification** - Hash verified using bot token
- ✅ **Expiration check** - initData must be < 1 hour old
- ✅ **HTTPS only** - All API calls use HTTPS
- ✅ **JWT tokens** - Secure session tokens
- ✅ **Token refresh** - Automatic token refresh on expiry
- ✅ **No password storage** - Passwordless authentication
- ✅ **Telegram ID as unique key** - Prevents duplicate accounts

---

## API Endpoints

### POST /api/v1/auth/telegram
**Request:**
```
Headers:
  Authorization: tma <initData>

Body: (empty)
```

**Response (200):**
```json
{
  "access_token": "...",
  "refresh_token": "...",
  "user": { ... }
}
```

**Response (401):**
```json
{
  "error": "Invalid Telegram data",
  "details": "..."
}
```

### POST /api/v1/auth/refresh
**Request:**
```json
{
  "refresh_token": "..."
}
```

**Response (200):**
```json
{
  "access_token": "...",
  "refresh_token": "..."
}
```

---

## Frontend Components

### AuthGuard
- Wraps entire app
- Auto-authenticates on load
- Handles routing based on auth state
- Shows loading/error screens

### WelcomeScreen
- Fallback for manual login
- Shows "Open in Telegram" if not in Telegram
- Manual retry button

### Onboarding Screens
- ProfileSetup
- PhotoUpload
- Interests
- GenderPreference

---

## Testing

### Test with Invalid Data
```bash
curl -X POST https://api.lomi.social/api/v1/auth/telegram \
  -H "Authorization: tma test"
# Expected: 401
```

### Test with Valid initData
```bash
# Get real initData from Telegram WebApp
curl -X POST https://api.lomi.social/api/v1/auth/telegram \
  -H "Authorization: tma <real-initData>"
# Expected: 200 with tokens
```

---

## Troubleshooting

### "initData not found"
- Check if app is opened from Telegram
- Wait for Telegram WebApp to initialize
- Check browser console for errors

### "Invalid Telegram data"
- initData expired (reopen app)
- Hash verification failed
- Bot token mismatch

### "Could not create user"
- Database connection issue
- Missing required fields
- Check backend logs

---

## Flow Diagram

```
User Opens App in Telegram
    ↓
AuthGuard Detects Telegram
    ↓
Extract initData
    ↓
POST /api/v1/auth/telegram
    ↓
Backend Validates initData
    ↓
Find/Create User
    ↓
Generate JWT Tokens
    ↓
Store Tokens (Frontend)
    ↓
Check has_profile
    ↓
┌─────────────┴─────────────┐
│                           │
has_profile = false    has_profile = true
│                           │
↓                           ↓
Onboarding              Main App
```

---

## Implementation Files

**Frontend:**
- `frontend/src/components/AuthGuard.tsx` - Auto-authentication
- `frontend/src/store/authStore.ts` - Auth state management
- `frontend/src/api/auth.ts` - Auth API calls
- `frontend/src/utils/telegram.ts` - Telegram utilities

**Backend:**
- `backend/internal/handlers/auth.go` - Auth endpoint
- `backend/internal/utils/jwt.go` - Token generation
- `backend/internal/models/user.go` - User model

---

## Next Steps

1. ✅ Auto-authentication on app load
2. ✅ Secure initData validation
3. ✅ User creation with minimal fields
4. ✅ Onboarding flow routing
5. ✅ Token management
6. ✅ Error handling
7. ✅ "Open in Telegram" fallback

**Complete!** 🎉

