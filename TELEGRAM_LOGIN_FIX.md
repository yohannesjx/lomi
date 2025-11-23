# Telegram Login 405 Error - Fixed

## Issues Found and Fixed

### 1. **Request Interceptor Overwriting Authorization Header** ✅ FIXED
**Problem:** The axios request interceptor was overwriting the `Authorization: tma <initData>` header with `Authorization: Bearer <token>` (which doesn't exist during login).

**Fix:** Updated `frontend/src/api/client.ts` to check if Authorization header already exists before adding Bearer token:
```typescript
if (!config.headers.Authorization) {
    const token = await storage.getItem('lomi_access_token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
}
```

### 2. **API URL Configuration** ✅ FIXED
**Problem:** Using absolute URLs was causing issues with different environments (ngrok, IP, domain).

**Fix:** Changed to use relative path `/api/v1` for all production environments. Caddy handles the routing:
```typescript
if (typeof window !== 'undefined') {
    return '/api/v1';  // Relative path - works everywhere
}
```

### 3. **Backend Dependencies** ✅ FIXED
**Problem:** Missing Go dependency for Telegram init-data validation library.

**Fix:** Ran `go mod tidy` to install `github.com/telegram-mini-apps/init-data-golang`.

## How to Test

### 1. Rebuild and Deploy Backend
```bash
# On your server
cd /path/to/lomi_mini
docker-compose -f docker-compose.prod.yml down
docker-compose -f docker-compose.prod.yml up -d --build
```

### 2. Rebuild Frontend
```bash
# On your local machine
cd frontend
npm run build  # or expo build:web
```

### 3. Deploy Frontend
```bash
# Copy build to server
scp -r dist/* user@152.53.87.200:/var/www/lomi-frontend/
```

### 4. Test in Telegram
1. Open Telegram app
2. Search for your bot
3. Open Mini App from bot menu
4. Click "Continue with Telegram"
5. Should now login successfully!

## What Changed

### Frontend (`frontend/src/api/client.ts`)
- ✅ Request interceptor now preserves Authorization header
- ✅ API URL uses relative path `/api/v1`
- ✅ Better error logging

### Backend (`backend/internal/handlers/auth.go`)
- ✅ Using official Telegram init-data library
- ✅ Enhanced logging for debugging
- ✅ Support for both Mini App and Widget login

## Verification

Test the endpoint directly:
```bash
# From your server
curl -X POST http://localhost:8080/api/v1/auth/telegram \
  -H "Content-Type: application/json" \
  -H "Authorization: tma <real_init_data_here>"
```

## Common Issues

### Still getting Network Error?
1. **Check backend is running:**
   ```bash
   docker-compose -f docker-compose.prod.yml ps
   docker-compose -f docker-compose.prod.yml logs backend
   ```

2. **Check Caddy is routing correctly:**
   ```bash
   curl http://localhost:8080/api/v1/health
   curl https://lomi.social/api/v1/health
   ```

3. **Check frontend build:**
   - Make sure you rebuilt the frontend after changes
   - Clear browser cache
   - Check browser console for errors

### Getting 401 Unauthorized?
- This is expected if initData is invalid/expired
- Make sure you're opening from Telegram (not browser)
- Check backend logs for validation errors

### Getting CORS errors?
- Check Caddy configuration has CORS headers
- Restart Caddy: `sudo systemctl restart caddy`

## Next Steps

1. ✅ Deploy changes to server
2. ✅ Test login in Telegram
3. ✅ Monitor backend logs
4. ✅ Verify user is created in database

## Support

If you still have issues:
1. Check backend logs: `docker-compose logs backend -f`
2. Check browser console in Telegram
3. Verify Caddy is running: `sudo systemctl status caddy`
4. Test API health endpoint: `curl https://lomi.social/api/v1/health`
