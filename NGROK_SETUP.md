# ngrok Setup for Telegram Mini App Development

## Why ngrok?

Telegram Mini Apps require:
- ✅ **HTTPS** (not HTTP)
- ✅ **Publicly accessible URL** (not localhost)

ngrok creates a secure HTTPS tunnel to your local server.

## Quick Start

### 1. Install ngrok

```bash
# macOS
brew install ngrok

# Or download from https://ngrok.com/download
# Sign up for free account at https://dashboard.ngrok.com
```

### 2. Authenticate ngrok (First time only)

```bash
ngrok config add-authtoken YOUR_AUTH_TOKEN
# Get token from: https://dashboard.ngrok.com/get-started/your-authtoken
```

### 3. Start your backend

```bash
# Option A: Docker (already running)
docker-compose up backend

# Option B: Local Go
cd backend
go run cmd/api/main.go
```

### 4. Start ngrok tunnel

```bash
ngrok http 8080
```

You'll see output like:
```
Forwarding  https://abc123.ngrok-free.app -> http://localhost:8080
```

### 5. Copy the HTTPS URL

Copy the `https://` URL (e.g., `https://abc123.ngrok-free.app`)

### 6. Update frontend API URL

**Option A: Environment Variable (Recommended)**

Create `frontend/.env`:
```bash
EXPO_PUBLIC_API_URL=https://abc123.ngrok-free.app/api/v1
```

**Option B: Update client.ts directly**

Update `frontend/src/api/client.ts`:
```typescript
const DEV_API_URL = 'https://abc123.ngrok-free.app/api/v1';
```

### 7. Restart frontend

```bash
cd frontend
npx expo start --clear
```

### 8. Create Telegram Mini App

1. Go to [@BotFather](https://t.me/BotFather)
2. Send `/newapp`
3. Select your bot: `lomi_social_bot`
4. When asked for Web App URL, use:
   ```
   https://abc123.ngrok-free.app
   ```
   (Use the ngrok URL, NOT `/api/v1` - that's just for the frontend)

### 8. Test

1. Open your bot in Telegram
2. Click menu → Your Mini App
3. The app should load and connect to your local backend via HTTPS

## Important Notes

### ngrok URLs Change

- **Free plan**: URL changes every time you restart ngrok
- **Solution**: Use environment variables to easily update
- **Paid plan**: Get static domain (e.g., `your-app.ngrok.io`)

### ngrok Dashboard

Access at: `http://localhost:4040`
- See all requests
- Inspect request/response
- Replay requests
- Very useful for debugging!

### Multiple Tunnels

If you need both frontend and backend:

```bash
# Terminal 1: Backend
ngrok http 8080

# Terminal 2: Frontend  
ngrok http 19000
```

Then use both URLs in Telegram Mini App config.

## Alternative: Cloudflare Tunnel (Free, Static Domain)

```bash
# Install
brew install cloudflared

# Create tunnel
cloudflared tunnel --url http://localhost:8080
```

## Production Setup

For production, deploy to a real server:
- Railway
- Render
- Fly.io
- DigitalOcean
- AWS/GCP/Azure

Then use the production HTTPS URL in Telegram Mini App.

## Troubleshooting

### "ngrok: command not found"
- Install ngrok: `brew install ngrok`
- Or add to PATH

### "Tunnel not found"
- Make sure backend is running on port 8080
- Check: `curl http://localhost:8080/api/v1/health`

### "CORS errors"
- Backend CORS is configured to allow all origins
- If issues persist, check `backend/cmd/api/main.go` CORS settings

### "Connection refused"
- Verify ngrok is forwarding to correct port
- Check backend is actually running
- Try accessing ngrok URL directly in browser

## Quick Script

Create `start-dev.sh`:

```bash
#!/bin/bash

# Start backend
cd backend && go run cmd/api/main.go &
BACKEND_PID=$!

# Wait for backend
sleep 3

# Start ngrok
echo "Starting ngrok..."
ngrok http 8080 &
NGROK_PID=$!

echo ""
echo "✅ Backend: http://localhost:8080"
echo "✅ ngrok dashboard: http://localhost:4040"
echo ""
echo "Copy the HTTPS URL from ngrok and:"
echo "1. Set EXPO_PUBLIC_API_URL in frontend/.env"
echo "2. Use the URL in Telegram Mini App config"
echo ""
echo "Press Ctrl+C to stop"

wait $BACKEND_PID $NGROK_PID
```

Make executable:
```bash
chmod +x start-dev.sh
```

