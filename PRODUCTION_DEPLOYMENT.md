# Production Deployment Checklist

## Backend Deployment

### Environment Variables Required

Set these in your production environment:

```bash
# Application
APP_ENV=production
APP_PORT=8080
APP_NAME=Lomi Social API

# Database (PostgreSQL)
DB_HOST=your-db-host
DB_PORT=5432
DB_USER=your-db-user
DB_PASSWORD=your-db-password
DB_NAME=lomi_db
DB_SSL_MODE=require

# Redis
REDIS_HOST=your-redis-host
REDIS_PORT=6379
REDIS_PASSWORD=your-redis-password
REDIS_DB=0

# Cloudflare R2 Storage
S3_ENDPOINT=https://a53cdfc7c678dac2a028159bcd178da2.r2.cloudflarestorage.com
S3_ACCESS_KEY=c7df163be474aae5317aa530bd8448bf
S3_SECRET_KEY=bef8caaaf0ef6577e5b4aa0b16bd2d08600d405a3b27d82ee002b1f203fe35a5
S3_USE_SSL=true
S3_REGION=auto
S3_BUCKET_PHOTOS=lomi-photos
S3_BUCKET_VIDEOS=lomi-videos
S3_BUCKET_GIFTS=lomi-gifts
S3_BUCKET_VERIFICATIONS=lomi-verifications

# JWT Authentication
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_ACCESS_EXPIRY=24h
JWT_REFRESH_EXPIRY=168h

# Telegram Bot
TELEGRAM_BOT_TOKEN=8453633918:AAE6UxkHrplAxyKXXBLt56bQufhZpH-rVEM
```

### Deployment Platforms

**Recommended:**
- **Railway** - Easy PostgreSQL + Redis + Go deployment
- **Render** - Free tier available, easy setup
- **Fly.io** - Great for Docker deployments
- **DigitalOcean App Platform** - Simple and reliable
- **AWS/GCP/Azure** - For enterprise scale

### Build & Deploy

```bash
# Build the Go binary
cd backend
go build -o bin/api ./cmd/api

# Or use Docker
docker build -t lomi-backend ./backend
docker push your-registry/lomi-backend
```

## Frontend Deployment

### Update API URL

Update `frontend/src/api/client.ts`:
```typescript
const PROD_API_URL = 'https://api.lomi.social/api/v1'; // Your backend URL
```

Or use environment variable:
```bash
EXPO_PUBLIC_API_URL=https://api.lomi.social/api/v1
```

### Deploy Options

**Option 1: Expo EAS (Recommended)**
```bash
cd frontend
npm install -g eas-cli
eas build --platform web
eas update
```

**Option 2: Static Hosting**
- Vercel
- Netlify
- Cloudflare Pages
- GitHub Pages

**Option 3: Self-hosted**
- Nginx
- Apache
- Any static file server

## Telegram Mini App Setup

1. Go to [@BotFather](https://t.me/BotFather)
2. Send `/newapp`
3. Select: `lomi_social_bot`
4. Provide:
   - **Title**: Lomi Social
   - **Short name**: lomi
   - **Description**: Find your Lomi in Ethiopia
   - **Photo**: Upload app icon
   - **Web App URL**: `https://your-frontend-domain.com`
   - **Short name**: lomi

## Post-Deployment Checklist

- [ ] Backend deployed with HTTPS
- [ ] Frontend deployed with HTTPS
- [ ] Database migrations run
- [ ] R2 buckets created
- [ ] CORS configured correctly
- [ ] Environment variables set
- [ ] Telegram Mini App configured
- [ ] Test authentication flow
- [ ] Test media upload
- [ ] Monitor logs for errors

## Testing Production

1. Open bot in Telegram
2. Launch Mini App
3. Test login flow
4. Verify API calls work
5. Check Cloudflare R2 uploads

## Security Checklist

- [ ] JWT_SECRET is strong and unique
- [ ] Database credentials are secure
- [ ] R2 credentials are secure
- [ ] CORS is properly configured
- [ ] HTTPS is enabled
- [ ] Rate limiting configured
- [ ] Error messages don't leak sensitive info

