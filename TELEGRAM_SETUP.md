# Telegram Bot Integration Setup

## ✅ Configuration Complete

Your Telegram bot has been configured:

- **Bot Username**: `lomi_social_bot`
- **Bot Token**: `8453633918:AAE6UxkHrplAxyKXXBLt56bQufhZpH-rVEM`
- **Backend**: Configured in `.env` and `docker-compose.yml`
- **Frontend**: Telegram WebApp integration ready

## How It Works

### 1. Telegram WebApp Authentication

When users open your app from Telegram:
1. Telegram provides `initData` containing user information
2. Frontend sends `initData` to backend `/api/v1/auth/telegram`
3. Backend validates the data using HMAC-SHA256 with bot token
4. Backend creates/finds user and returns JWT tokens
5. Frontend stores tokens and navigates to app

### 2. Backend Validation

The backend validates Telegram `initData` by:
- Extracting hash from query string
- Computing HMAC-SHA256 using bot token
- Comparing computed hash with provided hash
- Parsing user data if validation succeeds

### 3. User Flow

```
User opens app in Telegram
    ↓
WelcomeScreen loads
    ↓
User clicks "Continue with Telegram"
    ↓
Frontend gets initData from Telegram WebApp
    ↓
Sends to backend /auth/telegram
    ↓
Backend validates and creates/updates user
    ↓
Returns JWT tokens
    ↓
Frontend stores tokens and navigates
    ↓
- If profile complete → Main screen
- If profile incomplete → ProfileSetup screen
```

## Testing

### Option 1: Test in Telegram (Recommended)

1. Create a Telegram Mini App:
   - Go to [@BotFather](https://t.me/BotFather)
   - Send `/newapp`
   - Select your bot: `lomi_social_bot`
   - Provide app details:
     - Title: Lomi Social
     - Short name: lomi
     - Description: Find your Lomi in Ethiopia
     - Photo: Upload app icon
     - Web App URL: `https://your-domain.com` (or use ngrok for local testing)

2. Test the app:
   - Open your bot in Telegram
   - Click the menu button
   - Click your Mini App
   - The app should open with Telegram WebApp context

### Option 2: Local Development Testing

For local development without Telegram:

1. The app will detect if Telegram WebApp is not available
2. In development mode (`__DEV__`), it will skip authentication
3. You can test the UI flow directly

### Option 3: Using ngrok for Local Testing

1. Install ngrok: `brew install ngrok` (or download from ngrok.com)
2. Start your frontend: `cd frontend && npx expo start --web`
3. Expose with ngrok: `ngrok http 19000`
4. Use the ngrok URL in BotFather when creating the Mini App

## API Endpoints

### POST /api/v1/auth/telegram

**Request:**
```json
{
  "init_data": "query=string&hash=..."
}
```

**Response:**
```json
{
  "access_token": "jwt_token_here",
  "refresh_token": "refresh_token_here",
  "user": {
    "id": "uuid",
    "name": "User Name",
    "is_verified": false,
    "has_profile": false
  }
}
```

### POST /api/v1/auth/refresh

**Request:**
```json
{
  "refresh_token": "refresh_token_here"
}
```

**Response:**
```json
{
  "access_token": "new_jwt_token",
  "refresh_token": "new_refresh_token"
}
```

## Security Notes

1. **Bot Token**: Never commit the bot token to version control
   - ✅ Already in `.gitignore`
   - ✅ Stored in `.env` file

2. **InitData Validation**: The backend validates all Telegram initData
   - Uses HMAC-SHA256 with bot token
   - Prevents tampering with user data

3. **JWT Tokens**: 
   - Access tokens expire in 24 hours (configurable)
   - Refresh tokens expire in 7 days (configurable)
   - Stored securely using Expo SecureStore

## Troubleshooting

### "Invalid Telegram data" error

- Check that bot token is correct in `.env`
- Ensure `initData` is being sent correctly
- Verify the bot token matches the one in BotFather

### "Telegram WebApp not available"

- This is normal in development
- In production, ensure app is opened from Telegram
- Check that Telegram WebApp script is loaded

### Backend not validating

- Restart backend after updating `.env`
- Check logs for validation errors
- Verify bot token is loaded correctly

## Next Steps

1. **Set up Telegram Mini App**:
   - Use BotFather to create the Mini App
   - Point it to your deployed frontend URL

2. **Deploy Backend**:
   - Ensure bot token is set in production environment
   - Test the `/auth/telegram` endpoint

3. **Test Authentication**:
   - Open app from Telegram
   - Verify login flow works
   - Check user creation in database

## Files Modified

- `backend/.env` - Bot token added
- `backend/internal/handlers/auth.go` - Telegram login handler
- `backend/internal/utils/telegram.go` - Validation logic
- `frontend/src/utils/telegram.ts` - WebApp integration
- `frontend/src/screens/onboarding/WelcomeScreen.tsx` - Login flow
- `frontend/src/store/authStore.ts` - Token management
- `frontend/src/api/auth.ts` - API service

