# Get JWT Token - Curl Commands

## Method 1: Using Telegram initData (Recommended)

When you open the Mini App in Telegram, you get `initData` automatically. Use it to authenticate:

```bash
# Replace <YOUR_INIT_DATA> with the actual initData from Telegram
curl -X POST http://localhost/api/v1/auth/telegram \
  -H "Authorization: tma <YOUR_INIT_DATA>" \
  -v
```

**Example:**
```bash
curl -X POST http://localhost/api/v1/auth/telegram \
  -H "Authorization: tma query_id=AAHdF6IQAAAAAN0XohDhrOrc&user=%7B%22id%22%3A279058397%2C%22first_name%22%3A%22Vladislav%22%2C%22last_name%22%3A%22Kibenko%22%2C%22username%22%3A%22vdkfrost%22%2C%22language_code%22%3A%22ru%22%7D&auth_date=1662771648&hash=c501b71e775f74ce10e377dea85a7ea24ecd640b223ea86dfe453e0eaed2e2b2" \
  -v
```

**Response will contain:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "uuid",
    "telegram_id": 123456789,
    "name": "John Doe",
    ...
  }
}
```

## Method 2: Extract Token from Response

```bash
# Get token and save it
TOKEN=$(curl -s -X POST http://localhost/api/v1/auth/telegram \
  -H "Authorization: tma <YOUR_INIT_DATA>" | jq -r '.access_token')

echo "Your JWT token: $TOKEN"
```

## Method 3: One-liner to Get and Use Token

```bash
# Get token and use it immediately
TOKEN=$(curl -s -X POST http://localhost/api/v1/auth/telegram \
  -H "Authorization: tma <YOUR_INIT_DATA>" | jq -r '.access_token')

# Now use it for other requests
curl -X GET "http://localhost/api/v1/users?telegram_only=true" \
  -H "Authorization: Bearer $TOKEN" | jq '.'
```

## How to Get initData from Telegram

### Option A: From Browser Console (Web)
1. Open your Mini App in Telegram Web
2. Press F12 to open DevTools
3. Go to Console tab
4. Type: `window.Telegram.WebApp.initData`
5. Copy the result

### Option B: From React Native/Expo
1. Add this to your code temporarily:
```javascript
import * as Linking from 'expo-linking';

// Get initData
const initData = Linking.getInitialURL()?.split('tgWebAppData=')[1];
console.log('initData:', initData);
```

### Option C: From Telegram Mini App URL
The initData is in the URL when opening the Mini App:
```
https://your-domain.com/?tgWebAppData=query_id=...&user=...&hash=...
```

Extract everything after `tgWebAppData=`

## Quick Test Script

```bash
#!/bin/bash

# Get initData from user
read -p "Enter Telegram initData: " INIT_DATA

if [ -z "$INIT_DATA" ]; then
  echo "❌ initData is required"
  exit 1
fi

# Authenticate
echo "🔐 Authenticating..."
RESPONSE=$(curl -s -X POST http://localhost/api/v1/auth/telegram \
  -H "Authorization: tma $INIT_DATA")

# Check if successful
TOKEN=$(echo $RESPONSE | jq -r '.access_token' 2>/dev/null)

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
  echo "❌ Authentication failed"
  echo "Response: $RESPONSE"
  exit 1
fi

echo "✅ Authentication successful!"
echo ""
echo "JWT Token:"
echo "$TOKEN"
echo ""
echo "Save it for future use:"
echo "export TOKEN=\"$TOKEN\""
```

## Using the Token

Once you have the token, use it in all protected endpoints:

```bash
# Set token as variable
export TOKEN="your-jwt-token-here"

# Get your profile
curl -X GET "http://localhost/api/v1/users/me" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

# List all Telegram users
curl -X GET "http://localhost/api/v1/users?telegram_only=true" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

# Get upload URL
curl -X GET "http://localhost/api/v1/users/media/upload-url?media_type=photo" \
  -H "Authorization: Bearer $TOKEN" | jq '.'
```

## Troubleshooting

### Error: "sign is missing" or "Invalid Telegram data"
- Make sure you're using the full initData string
- Check that `TELEGRAM_BOT_TOKEN` is set in backend environment

### Error: 401 Unauthorized
- Token might be expired (default: 24 hours)
- Get a new token using initData

### Error: "Could not create user"
- Check backend logs for database errors
- Verify database connection is working

