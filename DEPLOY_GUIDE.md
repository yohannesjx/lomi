# 🚀 Quick Deployment Guide - Telegram Login Fix

## ✅ What's Been Fixed

1. **Frontend** - Request interceptor no longer overwrites Authorization header
2. **Caddyfile** - Fixed routing and header conflicts
3. **Deploy Script** - Updated to use fixed Caddyfile

## 📋 Deployment Steps

### Step 1: Commit and Push Changes

```bash
# On your local machine
cd /Users/gashawarega/Documents/Projects/lomi_mini

# Add all changes
git add .

# Commit with descriptive message
git commit -m "Fix: Telegram login 405 error - Fixed interceptor and Caddy config"

# Push to GitHub
git push origin main
```

### Step 2: Deploy on Server

```bash
# SSH to your server
ssh user@152.53.87.200

# Navigate to project directory
cd /opt/lomi_mini  # or wherever your project is

# Run the deployment script
./deploy-all.sh
```

**That's it!** The script will:
- ✅ Pull latest code from GitHub
- ✅ Rebuild and restart backend
- ✅ Build and deploy frontend
- ✅ Update Caddy with fixed configuration
- ✅ Reload all services

### Step 3: Verify Deployment

```bash
# On server, check services
docker-compose -f docker-compose.prod.yml ps

# Test backend health
curl http://localhost:8080/api/v1/health

# Test through Caddy
curl https://lomi.social/api/v1/health

# Test auth endpoint (should get 401, not 405)
curl -X POST https://lomi.social/api/v1/auth/telegram \
  -H "Authorization: tma test_data"
```

### Step 4: Test in Telegram

1. Open **Telegram app** (not browser!)
2. Search for your bot
3. Open **Mini App** from bot menu
4. Click **"Continue with Telegram"**
5. Should login successfully! ✅

## 🔍 Monitoring

### Watch Backend Logs
```bash
docker-compose -f docker-compose.prod.yml logs backend -f
```

You should see:
```
🔐 Login request received. Method: POST, Path: /api/v1/auth/telegram
📋 Authorization header present: true
✅ Hash verified successfully
```

### Watch Caddy Logs
```bash
sudo journalctl -u caddy -f
```

### Check Access Logs
```bash
sudo tail -f /var/log/caddy/lomi-access.log
```

## ⚠️ Troubleshooting

### If deploy-all.sh fails:

1. **Check script permissions:**
   ```bash
   chmod +x deploy-all.sh
   ```

2. **Check .env.production exists:**
   ```bash
   ls -la .env.production
   ```

3. **Run with verbose output:**
   ```bash
   bash -x ./deploy-all.sh
   ```

### If Caddy validation fails:

```bash
# Check Caddyfile syntax
sudo caddy validate --config /etc/caddy/Caddyfile

# View Caddy status
sudo systemctl status caddy

# Restart Caddy if needed
sudo systemctl restart caddy
```

### If backend doesn't start:

```bash
# Check backend logs
docker-compose -f docker-compose.prod.yml logs backend

# Rebuild backend
docker-compose -f docker-compose.prod.yml build backend --no-cache
docker-compose -f docker-compose.prod.yml up -d
```

### If still getting network error:

1. **Clear browser cache** in Telegram
2. **Force close** Telegram app and reopen
3. **Check backend is running:**
   ```bash
   curl http://localhost:8080/api/v1/health
   ```
4. **Check Caddy is proxying:**
   ```bash
   curl https://lomi.social/api/v1/health
   ```

## 📊 What Changed

### Files Modified:
- ✅ `frontend/src/api/client.ts` - Fixed interceptor
- ✅ `Caddyfile` - Fixed routing and headers
- ✅ `deploy-all.sh` - Updated to use fixed Caddyfile
- ✅ `backend/go.mod` - Added Telegram init-data library

### Changes Summary:
1. **Frontend interceptor** - Preserves Authorization header
2. **API URL** - Uses relative path `/api/v1`
3. **Caddyfile** - Uses `handle_path` and proper header ordering
4. **Backend** - Enhanced logging and validation

## ✅ Success Indicators

You'll know it's working when:
- ✅ `curl https://lomi.social/api/v1/health` returns 200 OK
- ✅ POST to `/auth/telegram` returns 401 (not 405)
- ✅ Telegram login button works without network error
- ✅ Backend logs show successful login attempts

## 🎉 Next Steps After Successful Deployment

1. Test full login flow in Telegram
2. Verify user is created in database
3. Test other app features
4. Monitor logs for any issues

---

**Need help?** Check the logs:
- Backend: `docker-compose logs backend -f`
- Caddy: `sudo journalctl -u caddy -f`
- Access: `sudo tail -f /var/log/caddy/lomi-access.log`
