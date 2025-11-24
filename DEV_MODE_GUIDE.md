# Development Mode Guide

## Understanding Development vs Production

### Telegram Authentication

**Why it doesn't work locally:**
- Telegram WebApp requires the app to be opened **from within Telegram**
- Telegram needs a **publicly accessible HTTPS URL**
- Local development (localhost) is not accessible to Telegram's servers

**This is normal and expected!** ✅

### Development Mode Features

The app now includes **development mode bypasses** to allow UI testing:

#### 1. **Welcome Screen**
- Automatically detects dev mode
- Skips Telegram authentication
- Navigates directly to Profile Setup
- No error messages - just continues

#### 2. **Photo Upload**
- In dev mode without auth: Stores photos locally
- Shows photos in UI for testing
- Skips actual R2 upload
- Allows testing the full flow

#### 3. **Profile Flow**
- All screens work without authentication
- UI/UX can be tested fully
- Navigation flows work correctly

## How to Test Locally

### Current Setup (Development Mode)
1. ✅ Open app in browser/simulator
2. ✅ Click "Continue with Telegram"
3. ✅ Automatically proceeds to Profile Setup
4. ✅ Upload photos (stored locally in dev mode)
5. ✅ Complete onboarding flow
6. ✅ Test all screens and features

### What Works in Dev Mode
- ✅ All UI screens
- ✅ Navigation flows
- ✅ Photo selection (stored locally)
- ✅ Form inputs
- ✅ Swipe gestures
- ✅ Chat interface
- ✅ Profile screens

### What Requires Production
- ❌ Actual Telegram authentication
- ❌ Real photo uploads to R2
- ❌ Backend API calls (some)
- ❌ Real user data

## Testing Full Authentication

To test Telegram authentication, you need:

### Option 1: Deploy to Production
1. Deploy frontend to Vercel/Netlify/Expo
2. Get HTTPS URL
3. Configure in BotFather
4. Open from Telegram
5. Full authentication works

### Option 2: Use ngrok (Temporary)
1. Start backend: `docker-compose up backend`
2. Start ngrok: `ngrok http 8080`
3. Get HTTPS URL from ngrok
4. Configure in BotFather
5. Open from Telegram
6. Test authentication

### Option 3: Test Backend Separately
- Use Postman/curl to test API endpoints
- Test authentication flow manually
- Verify R2 uploads work

## Development Workflow

### Recommended Approach:
1. **UI Development** → Test locally in dev mode ✅
2. **Backend Testing** → Test API with Postman/curl
3. **Integration Testing** → Deploy to staging
4. **Full Testing** → Deploy to production + Telegram

### Current Status:
- ✅ UI development works perfectly in dev mode
- ✅ All screens are functional
- ✅ Navigation flows work
- ✅ Can test complete user journey
- ⚠️ Authentication requires production deployment

## Console Messages

You'll see these messages in dev mode (this is normal):

```
⚠️ Telegram WebApp not available in local development.
ℹ️  This is expected. Telegram auth only works when:
   1. App is deployed to a public HTTPS URL
   2. App is opened from within Telegram
   3. Mini App is configured in BotFather

💡 For now, you can test the UI flow without authentication.
```

**This is informational, not an error!** ✅

## Summary

- ✅ **Dev mode works great for UI testing**
- ✅ **All features are testable locally**
- ✅ **Authentication requires production**
- ✅ **This is the expected behavior**

Continue developing and testing the UI locally. When ready, deploy to test full authentication!


