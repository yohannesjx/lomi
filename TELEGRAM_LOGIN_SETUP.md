# Telegram Login Setup - Production

## ✅ Current Status

Your Telegram authentication is set up and ready! Here's what's configured:

### Backend ✅
- Telegram auth endpoint: `/api/v1/auth/telegram`
- Validates Telegram `initData` using HMAC-SHA256
- Creates/finds user and returns JWT tokens
- Bot token configured: `8453633918:AAE6UxkHrplAxyKXXBLt56bQufhZpH-rVEM`

### Frontend ✅
- Telegram WebApp SDK integration
- WelcomeScreen with "Continue with Telegram" button
- Auth store with login function
- API client configured for production: `https://api.lomi.social/api/v1`

## 🚀 How It Works

1. **User opens app in Telegram**
   - Opens Mini App from bot menu
   - Telegram injects `initData` into the page

2. **Frontend gets initData**
   - `getTelegramInitData()` reads from `window.Telegram.WebApp.initData`
   - User clicks "Continue with Telegram"

3. **Backend validates**
   - Frontend sends `initData` to `/api/v1/auth/telegram`
   - Backend validates using bot token
   - Creates/finds user in database

4. **User logged in**
   - Backend returns JWT tokens
   - Frontend stores tokens
   - Navigates to Main or ProfileSetup

## 🔧 Configuration Check

### 1. Verify Bot Token in Backend

On your server, check `.env.production`:

```bash
# On server
cat .env.production | grep TELEGRAM_BOT_TOKEN
# Should show: TELEGRAM_BOT_TOKEN=8453633918:AAE6UxkHrplAxyKXXBLt56bQufhZpH-rVEM
```

### 2. Verify API URL in Frontend

The frontend is configured to use:
- Production: `https://api.lomi.social/api/v1`
- Local dev: `http://localhost:8080/api/v1` (only when on localhost)

### 3. Verify Mini App URL in BotFather

1. Go to [@BotFather](https://t.me/BotFather)
2. Send `/myapps`
3. Select your bot
4. Verify Web App URL is: `https://lomi.social` (or `http://152.53.87.200` if DNS not ready)

## 🧪 Testing

### Test in Telegram (Production)

1. **Open your bot in Telegram**
   - Search for your bot username
   - Open the bot

2. **Open Mini App**
   - Click the menu button (☰)
   - Click your Mini App
   - App should open with Telegram context

3. **Test Login**
   - Click "Continue with Telegram"
   - Should authenticate and navigate

### Debug Steps

If login doesn't work:

1. **Check browser console** (in Telegram, tap menu → "Open in Browser")
   - Look for errors
   - Check if `window.Telegram.WebApp` exists
   - Check if `initData` is available

2. **Check backend logs**
   ```bash
   # On server
   docker-compose -f docker-compose.prod.yml logs backend | grep -i telegram
   ```

3. **Test API directly**
   ```bash
   # Get initData from browser console:
   # window.Telegram.WebApp.initData
   
   # Test endpoint
   curl -X POST https://api.lomi.social/api/v1/auth/telegram \
     -H "Content-Type: application/json" \
     -d '{"init_data":"YOUR_INIT_DATA_HERE"}'
   ```

## 🐛 Common Issues

### Issue: "Telegram WebApp not available"

**Cause:** App not opened from Telegram or script didn't load

**Fix:**
- Make sure you open the app from Telegram bot menu
- Check browser console for script loading errors
- Verify HTTPS is working (required for Telegram)

### Issue: "Invalid Telegram data"

**Cause:** Bot token mismatch or initData expired

**Fix:**
- Verify bot token in backend `.env.production`
- Make sure bot token matches the one used in BotFather
- Restart backend after changing token

### Issue: "401 Unauthorized"

**Cause:** API URL incorrect or CORS issue

**Fix:**
- Verify API URL is `https://api.lomi.social/api/v1`
- Check Caddy is proxying correctly
- Test API health: `curl https://api.lomi.social/api/v1/health`

## 📝 Next Steps

1. **Test login flow** in Telegram
2. **Check database** - verify users are being created
3. **Test profile setup** - complete onboarding flow
4. **Monitor logs** - watch for any errors

## 🔐 Security Notes

- ✅ Telegram `initData` is validated using HMAC-SHA256
- ✅ JWT tokens are used for authenticated requests
- ✅ Bot token is stored securely in environment variables
- ⚠️ Make sure `.env.production` is not committed to Git

## 📞 Support

If you encounter issues:
1. Check browser console for errors
2. Check backend logs: `docker-compose logs backend`
3. Test API endpoint directly with curl
4. Verify bot token matches BotFather

