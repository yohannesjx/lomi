# Troubleshooting Guide

## Issue 1: "Telegram WebApp not available" Message

### What it means:
This message appears when running the app locally (not from Telegram). This is **normal and expected** in development.

### Why it happens:
- Telegram WebApp only works when the app is opened from within Telegram
- When testing locally (Expo web, simulator, etc.), Telegram WebApp is not available
- The app detects this and allows you to continue testing the UI

### Solution:
**For Development:**
- This is expected behavior
- Click "Continue Testing" to proceed with UI testing
- Authentication will work when deployed and opened from Telegram

**For Production:**
- Deploy your app to a live URL
- Configure Telegram Mini App in BotFather
- Open the app from Telegram - authentication will work

### How to test Telegram authentication:
1. Deploy frontend to production (Vercel, Netlify, etc.)
2. Go to [@BotFather](https://t.me/BotFather)
3. Create Mini App with your production URL
4. Open the app from Telegram
5. Authentication will work automatically

---

## Issue 2: Photo Upload Error

### Common causes:

#### 1. Backend not running
**Check:**
```bash
curl http://localhost:8080/api/v1/health
```
Should return: `{"status":"ok","message":"Lomi Backend is running 🍋"}`

**Fix:**
```bash
# Start Docker backend
docker-compose up backend

# Or run locally
cd backend
go run cmd/api/main.go
```

#### 2. Authentication required
Photo upload requires authentication. Make sure you're logged in.

**Fix:**
- Complete Telegram login first
- Or use development mode to skip auth

#### 3. R2 credentials incorrect
**Check:**
- Verify R2 credentials in `backend/.env`
- Ensure buckets exist in Cloudflare R2 dashboard

**Fix:**
```bash
# Update backend/.env with correct credentials
S3_ENDPOINT=https://a53cdfc7c678dac2a028159bcd178da2.r2.cloudflarestorage.com
S3_ACCESS_KEY=c7df163be474aae5317aa530bd8448bf
S3_SECRET_KEY=bef8caaaf0ef6577e5b4aa0b16bd2d08600d405a3b27d82ee002b1f203fe35a5
```

#### 4. Network/CORS issues
**Check browser console for:**
- CORS errors
- Network errors
- 401/403 errors

**Fix:**
- Ensure backend CORS is configured correctly
- Check backend logs for errors

#### 5. File size too large
**Check:**
- Image file size
- R2 upload limits

**Fix:**
- Compress images before upload
- Check R2 bucket settings

### Debugging steps:

1. **Check browser console:**
   - Open DevTools (F12)
   - Look for error messages
   - Check Network tab for failed requests

2. **Check backend logs:**
   ```bash
   docker-compose logs backend
   ```

3. **Test pre-signed URL generation:**
   ```bash
   # Get upload URL (requires auth token)
   curl -X GET "http://localhost:8080/api/v1/users/media/upload-url?media_type=photo" \
     -H "Authorization: Bearer YOUR_TOKEN"
   ```

4. **Test R2 connection:**
   - Check backend logs for S3 connection message
   - Should see: `✅ Connected to S3-compatible storage (R2/MinIO)`

### Error messages explained:

- **"Failed to generate upload URL"** → Backend S3/R2 connection issue
- **"Upload failed: 403"** → R2 credentials incorrect or bucket permissions
- **"Upload failed: 401"** → Not authenticated
- **"Network error"** → Backend not running or CORS issue
- **"Failed to upload photo"** → Check backend logs for details

---

## Quick Fixes

### Reset everything:
```bash
# Restart backend
docker-compose restart backend

# Clear frontend cache
cd frontend
npx expo start --clear
```

### Check all services:
```bash
# Check Docker services
docker-compose ps

# Check backend health
curl http://localhost:8080/api/v1/health

# Check if R2 is accessible (from backend)
# Look for connection message in logs
```

### Common fixes:
1. ✅ Backend running on port 8080
2. ✅ R2 credentials correct in `.env`
3. ✅ R2 buckets created
4. ✅ User authenticated (or in dev mode)
5. ✅ Network connection stable

---

## Still having issues?

1. Check browser console for detailed errors
2. Check backend logs: `docker-compose logs backend`
3. Verify R2 credentials in Cloudflare dashboard
4. Test with a smaller image file
5. Try uploading from a different network

