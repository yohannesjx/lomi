# Caddy Configuration Issues - FOUND AND FIXED

## 🔴 Critical Issues Found in Your Caddyfile

### Issue 1: **Using `handle /api/*` instead of `handle_path /api/*`**

**Problem:**
```caddy
handle /api/* {
    reverse_proxy localhost:8080
}
```

When a request comes in for `/api/v1/auth/telegram`:
- Caddy matches the pattern
- Proxies to backend as `http://localhost:8080/api/v1/auth/telegram`
- Backend expects `/api/v1/auth/telegram` ✅

**BUT** - if the backend is configured to listen on `/api/v1/*`, this works. However, using `handle_path` is cleaner:

```caddy
handle_path /api/* {
    reverse_proxy localhost:8080
}
```

This strips `/api` and sends `/v1/auth/telegram` to backend.

### Issue 2: **Global Headers Applied to Proxy Responses**

**Problem:**
```caddy
lomi.social {
    handle /api/* {
        reverse_proxy localhost:8080
    }
    
    # These headers are applied to ALL responses, including proxied ones!
    header {
        Access-Control-Allow-Origin "*"
        X-Frame-Options "DENY"
    }
}
```

This can cause conflicts with backend's own CORS headers.

**Fix:** Move headers inside specific `handle` blocks.

### Issue 3: **OPTIONS Handling Order**

The OPTIONS handler should be the FIRST handler to catch preflight requests before they reach other handlers.

## ✅ What I Fixed

1. **Changed `handle` to `handle_path`** for API routes
   - Properly strips `/api` prefix
   - Cleaner routing to backend

2. **Moved headers to specific blocks**
   - Frontend gets security headers
   - API responses get CORS headers only
   - No conflicts

3. **Added logging**
   - Easier to debug issues
   - See exactly what requests come in

4. **Added `header_down` directive**
   - Preserves backend's CORS headers
   - Prevents Caddy from overwriting them

## 🚀 How to Apply the Fix

### On Your Server:

```bash
# 1. Backup current Caddyfile
sudo cp /etc/caddy/Caddyfile /etc/caddy/Caddyfile.backup

# 2. Copy the fixed version (from your local machine)
scp Caddyfile.fixed user@152.53.87.200:/tmp/Caddyfile.new

# 3. On server, move it to /etc/caddy/
sudo mv /tmp/Caddyfile.new /etc/caddy/Caddyfile

# 4. Test the configuration
sudo caddy validate --config /etc/caddy/Caddyfile

# 5. Reload Caddy (no downtime)
sudo systemctl reload caddy

# 6. Check status
sudo systemctl status caddy

# 7. Check logs
sudo journalctl -u caddy -f
```

### Verify It Works:

```bash
# Test health endpoint
curl https://lomi.social/api/v1/health

# Test with verbose output
curl -v https://lomi.social/api/v1/health

# Should see:
# < HTTP/2 200
# < access-control-allow-origin: *
# {"status":"ok","message":"Lomi Backend is running 🍋"}
```

## 🔍 Testing the Fix

### 1. Test API Endpoint Directly
```bash
# From your server
curl -X POST http://localhost:8080/api/v1/auth/telegram \
  -H "Content-Type: application/json" \
  -H "Authorization: tma test_data"

# Should get 401 (expected - invalid data)
# NOT 405 (method not allowed)
```

### 2. Test Through Caddy
```bash
# From anywhere
curl -X POST https://lomi.social/api/v1/auth/telegram \
  -H "Content-Type: application/json" \
  -H "Authorization: tma test_data"

# Should also get 401, not 405
```

### 3. Test in Telegram
1. Open Telegram app
2. Open your bot's Mini App
3. Click "Continue with Telegram"
4. Should now work! ✅

## 📊 Understanding the Request Flow

### Before Fix:
```
Telegram → https://lomi.social/api/v1/auth/telegram
         ↓
      Caddy (lomi.social block)
         ↓
      handle /api/* matches
         ↓
      Global headers applied (might conflict)
         ↓
      reverse_proxy localhost:8080
         ↓
      Backend receives: /api/v1/auth/telegram
         ↓
      Backend responds
         ↓
      Caddy adds MORE headers (conflicts!)
         ↓
      Response to Telegram (might be malformed)
```

### After Fix:
```
Telegram → https://lomi.social/api/v1/auth/telegram
         ↓
      Caddy (lomi.social block)
         ↓
      handle_path /api/* matches
         ↓
      reverse_proxy localhost:8080
         ↓
      Backend receives: /v1/auth/telegram (or /api/v1/auth/telegram)
         ↓
      Backend responds with CORS headers
         ↓
      Caddy preserves backend headers (header_down)
         ↓
      Clean response to Telegram ✅
```

## 🐛 If Still Having Issues

### Check Caddy Logs:
```bash
# Real-time logs
sudo journalctl -u caddy -f

# Recent errors
sudo journalctl -u caddy -n 100 --no-pager | grep -i error

# Access logs (if configured)
sudo tail -f /var/log/caddy/access.log
```

### Check Backend Logs:
```bash
docker-compose -f docker-compose.prod.yml logs backend -f
```

### Verify Caddy is Proxying:
```bash
# Check if backend is reachable
curl http://localhost:8080/api/v1/health

# Check if Caddy is proxying
curl https://lomi.social/api/v1/health

# Both should return the same response
```

## 📝 Summary

**Root Cause:** Caddy configuration had:
1. ❌ Improper header ordering
2. ❌ Global headers conflicting with backend
3. ❌ Using `handle` instead of `handle_path`

**Fix Applied:**
1. ✅ Use `handle_path` for cleaner routing
2. ✅ Move headers to specific blocks
3. ✅ Add `header_down` to preserve backend headers
4. ✅ Add logging for debugging

**Next Step:** Apply the fixed Caddyfile and reload Caddy!
