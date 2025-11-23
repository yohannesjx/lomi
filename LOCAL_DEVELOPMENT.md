# Local Development Setup for Telegram Mini App

## Problem

Telegram Mini Apps require:
- ✅ HTTPS (not HTTP)
- ✅ Publicly accessible URL (not localhost)

## Solutions

### Option 1: Use ngrok (Recommended for Local Development)

ngrok creates a secure tunnel to your local server.

#### Installation

```bash
# macOS
brew install ngrok

# Or download from https://ngrok.com/download
```

#### Setup

1. **Start your backend:**
   ```bash
   cd backend
   go run cmd/api/main.go
   # Or use Docker: docker-compose up backend
   ```

2. **Start ngrok tunnel:**
   ```bash
   ngrok http 8080
   ```

3. **Copy the HTTPS URL** (e.g., `https://abc123.ngrok.io`)

4. **Update frontend API client:**
   - Update `frontend/src/api/client.ts` to use the ngrok URL
   - Or set environment variable

5. **Create Telegram Mini App:**
   - Go to [@BotFather](https://t.me/BotFather)
   - Send `/newapp`
   - Use the ngrok HTTPS URL for the Web App URL

#### Note: ngrok URLs change each time
- Free ngrok URLs change on restart
- Consider ngrok paid plan for static domain
- Or use environment variables to easily update

### Option 2: Use Cloudflare Tunnel (Free, Static Domain)

1. Install `cloudflared`:
   ```bash
   brew install cloudflared
   ```

2. Create tunnel:
   ```bash
   cloudflared tunnel --url http://localhost:8080
   ```

3. Use the provided HTTPS URL

### Option 3: Deploy to Production/Staging

Deploy your backend to a server with HTTPS:
- Railway
- Render
- Fly.io
- DigitalOcean
- AWS/GCP/Azure

### Option 4: Test Backend Separately (Development Only)

For testing the backend API without Telegram:

1. Use Postman/curl to test `/auth/telegram` endpoint
2. Manually create test `initData` for validation
3. Test the full flow once deployed

## Quick Setup Script

Create a script to start ngrok automatically:

```bash
#!/bin/bash
# start-dev.sh

# Start backend
cd backend && go run cmd/api/main.go &
BACKEND_PID=$!

# Wait for backend to start
sleep 3

# Start ngrok
ngrok http 8080 &
NGROK_PID=$!

echo "Backend running on http://localhost:8080"
echo "ngrok tunnel starting..."
echo "Check ngrok dashboard: http://localhost:4040"
echo ""
echo "Press Ctrl+C to stop"

# Wait for user interrupt
wait $BACKEND_PID $NGROK_PID
```

## Environment Variables

Update your frontend to use environment-based URLs:

```typescript
// frontend/src/api/client.ts
const DEV_API_URL = __DEV__
  ? process.env.EXPO_PUBLIC_API_URL || 'http://localhost:8080/api/v1'
  : 'https://api.lomi.social/api/v1';
```

Set in `.env`:
```
EXPO_PUBLIC_API_URL=https://your-ngrok-url.ngrok.io/api/v1
```

## Testing Without Telegram

You can test the backend API directly:

```bash
# Test health endpoint
curl http://localhost:8080/api/v1/health

# Test Telegram login (with mock data)
curl -X POST http://localhost:8080/api/v1/auth/telegram \
  -H "Content-Type: application/json" \
  -d '{"init_data":"user=%7B%22id%22%3A123%7D&hash=test"}'
```

## Production Setup

For production:
1. Deploy backend with HTTPS
2. Update frontend API URL
3. Create Telegram Mini App with production URL
4. Test authentication flow

