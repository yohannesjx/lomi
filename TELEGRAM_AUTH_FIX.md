# Telegram Authentication - Complete Fix Guide

## 🔍 Current Issues

1. **405 Method Not Allowed** - POST requests not reaching backend
2. **initData Missing** - App opening in Safari instead of Telegram
3. **Authentication Failing** - Multiple attempts, still not working

## ✅ Step-by-Step Solution

### Step 1: Verify Backend is Reachable

Test if your backend is accessible:

```bash
# From your local machine
curl https://lomi.social/api/v1/test
curl https://lomi.social/api/v1/health

# Should return JSON responses
```

### Step 2: Test POST Request

```bash
# Test POST to auth endpoint
curl -X POST https://lomi.social/api/v1/test/auth \
  -H "Content-Type: application/json" \
  -d '{"test": "data"}'

# Should return: {"status":"ok","message":"Auth endpoint is reachable! ✅"}
```

### Step 3: Check Backend Logs

On your server:

```bash
# Watch backend logs in real-time
docker-compose -f docker-compose.prod.yml logs -f backend

# When you try to login, you should see:
# 🔐 Login request received. Method: POST, Path: /api/v1/auth/telegram...
```

**If you DON'T see this log**, the request isn't reaching the backend (Caddy routing issue).

### Step 4: Verify Caddy Configuration

Check if Caddy is routing correctly:

```bash
# On server
sudo caddy validate --config /etc/caddy/Caddyfile

# Check Caddy logs
sudo journalctl -u caddy -f
```

### Step 5: Test Telegram Login Widget (Alternative)

The Login Widget is simpler and doesn't require Mini App setup:

1. **Link your domain to bot:**
   - Go to [@BotFather](https://t.me/BotFather)
   - Send `/setdomain`
   - Select your bot
   - Enter: `lomi.social`

2. **Create a test HTML page:**

Create `/var/www/lomi-frontend/test-login.html`:

```html
<!DOCTYPE html>
<html>
<head>
    <title>Telegram Login Test</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 600px;
            margin: 50px auto;
            padding: 20px;
            background: #000;
            color: #fff;
        }
        .widget-container {
            text-align: center;
            margin: 30px 0;
        }
        .status {
            margin: 20px 0;
            padding: 10px;
            border-radius: 5px;
        }
        .success { background: #0a0; }
        .error { background: #a00; }
        .info { background: #006; }
    </style>
</head>
<body>
    <h1>🍋 Lomi Social - Telegram Login Test</h1>
    
    <div class="status info">
        <p><strong>Step 1:</strong> Link your domain in BotFather: <code>/setdomain</code></p>
        <p><strong>Step 2:</strong> Replace YOUR_BOT_USERNAME below with your bot username</p>
        <p><strong>Step 3:</strong> Click the Telegram button</p>
    </div>

    <div class="widget-container">
        <script async src="https://telegram.org/js/telegram-widget.js?22" 
            data-telegram-login="YOUR_BOT_USERNAME" 
            data-size="large" 
            data-auth-url="https://lomi.social/api/v1/auth/telegram/widget"
            data-request-access="write">
        </script>
    </div>

    <div id="status" class="status" style="display:none;"></div>

    <script>
        // Check if redirected back with auth data
        const urlParams = new URLSearchParams(window.location.search);
        if (urlParams.has('id') && urlParams.has('hash')) {
            document.getElementById('status').style.display = 'block';
            document.getElementById('status').className = 'status info';
            document.getElementById('status').innerHTML = 
                '<p>✅ Received auth data! Check backend logs to see if it was processed.</p>' +
                '<p>User ID: ' + urlParams.get('id') + '</p>';
        }
    </script>
</body>
</html>
```

3. **Access test page:**
   - Go to: `https://lomi.social/test-login.html`
   - Click the Telegram button
   - Should redirect back with auth data

### Step 6: Debug Mini App Login

If using Mini App (initData method):

1. **Verify app opens in Telegram:**
   - Open Telegram app (NOT Safari)
   - Find your bot
   - Tap menu (☰) → Mini App
   - Check browser console for `initData`

2. **Check initData availability:**
   ```javascript
   // In browser console (when opened from Telegram)
   console.log(window.Telegram?.WebApp?.initData);
   console.log(window.Telegram?.WebApp?.platform);
   ```

3. **Test API call manually:**
   ```javascript
   // In browser console
   const initData = window.Telegram?.WebApp?.initData;
   if (initData) {
       fetch('/api/v1/auth/telegram', {
           method: 'POST',
           headers: {
               'Content-Type': 'application/json',
               'Authorization': 'tma ' + initData
           }
       })
       .then(r => r.json())
       .then(console.log)
       .catch(console.error);
   }
   ```

## 🐛 Common Issues & Fixes

### Issue 1: 405 Method Not Allowed

**Cause:** Request not reaching backend or Caddy blocking POST

**Fix:**
```bash
# Check Caddyfile has proper routing
sudo cat /etc/caddy/Caddyfile | grep -A 5 "handle /api"

# Should show:
# handle /api/* {
#     reverse_proxy localhost:8080
# }

# Reload Caddy
sudo systemctl reload caddy
```

### Issue 2: initData Missing

**Cause:** App opened in Safari, not Telegram

**Fix:**
- MUST open from Telegram app
- Go to bot → Menu (☰) → Mini App
- NOT from Safari/bookmark

### Issue 3: CORS Errors

**Fix:** Already configured in Caddyfile, but verify:
```bash
# Check CORS headers
curl -I -X OPTIONS https://lomi.social/api/v1/auth/telegram \
  -H "Origin: https://lomi.social"
```

## 🎯 Quick Test Commands

```bash
# 1. Test backend health
curl https://lomi.social/api/v1/health

# 2. Test POST endpoint
curl -X POST https://lomi.social/api/v1/test/auth

# 3. Test with Authorization header (fake data)
curl -X POST https://lomi.social/api/v1/auth/telegram \
  -H "Authorization: tma test_data" \
  -H "Content-Type: application/json"

# 4. Check backend logs
docker-compose -f docker-compose.prod.yml logs backend | tail -50
```

## 📋 Checklist

- [ ] Backend is running (`/api/v1/health` returns OK)
- [ ] POST requests work (`/api/v1/test/auth` returns OK)
- [ ] Caddy is routing `/api/*` to backend
- [ ] Domain linked in BotFather (`/setdomain`)
- [ ] App opened from Telegram (not Safari)
- [ ] Backend logs show incoming requests
- [ ] Bot token is correct in `.env.production`

## 🚀 Next Steps

1. Run the test commands above
2. Check backend logs when testing
3. Use the test HTML page for Login Widget
4. Share the results so we can pinpoint the exact issue

The test endpoints will help us identify if it's:
- **Routing issue** (request not reaching backend)
- **Authentication issue** (initData validation failing)
- **Configuration issue** (bot token, domain, etc.)

