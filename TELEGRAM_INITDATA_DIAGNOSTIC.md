# Telegram initData Diagnostic Guide

## The Problem

You're getting "network error" because `initData` is missing. This happens when the app opens in Safari instead of Telegram's in-app browser.

## How to Verify

### Step 1: Check Browser Console

Open the app and check the browser console. You should see:

**✅ CORRECT (in Telegram):**
```
🔍 Telegram WebApp Debug Info: {
  webAppExists: true,
  hasInitData: true,  ← This should be TRUE
  initDataLength: 200+,  ← Should be > 0
  platform: "ios" or "android",  ← Should NOT be "unknown"
  ...
}
```

**❌ WRONG (in Safari):**
```
🔍 Telegram WebApp Debug Info: {
  webAppExists: true,
  hasInitData: false,  ← This is FALSE
  initDataLength: 0,  ← This is 0
  platform: "unknown",  ← This is "unknown"
  userAgent: "Safari..."  ← Shows Safari, not Telegram
  ...
}
```

### Step 2: Verify How You're Opening

**✅ CORRECT Way:**
1. Open **Telegram app** (not Safari)
2. Search for your bot
3. Open the bot
4. Click **menu button** (☰) at bottom
5. Click **Mini App** from menu

**❌ WRONG Ways:**
- Typing URL in Safari ❌
- Opening from browser bookmark ❌
- Sharing link and opening in browser ❌
- Opening from external app ❌

### Step 3: Check BotFather

```
1. Open @BotFather
2. Send /myapps
3. Select your bot
4. Check "Web App URL"
5. Should be: https://lomi.social/ (or https://lomi.social)
```

### Step 4: Test on Different Device

Try on:
- ✅ Another phone (friend's phone)
- ✅ Telegram Desktop
- ✅ Different Telegram account

## Quick Test

1. **Open Telegram app**
2. **Search for your bot**
3. **Open bot → Menu → Mini App**
4. **Check console** - Should see `✅ Found initData`

If you still see `❌ initData not found`, the app is opening in Safari, not Telegram.

## Solution

The app **MUST** be opened from Telegram's in-app browser. There's no workaround - Telegram only provides `initData` when opened from within Telegram.

If you're opening it correctly and still getting the error, check:
1. BotFather Mini App URL is correct
2. Telegram app is up to date
3. Try on a different device/account

## Debug Output

After rebuilding, the console will show detailed info:
- Where it's looking for initData
- What it finds (or doesn't find)
- Exact error messages

This will help identify the exact issue.

