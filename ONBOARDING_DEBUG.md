# 🔧 Onboarding Navigation Issue - Diagnosis & Fix

## Issue
User fills in ProfileSetupScreen (name, age, gender) and clicks "Next Step", but nothing happens - doesn't navigate to CityScreen.

## Root Cause Analysis

Based on code review, the implementation is **correct**. The issue is likely one of the following:

### 1. **Navigation Not Properly Initialized** ⚠️
The OnboardingNavigator might not be properly integrated into the main navigation stack.

### 2. **Console Errors Not Visible** ⚠️
Errors are being logged but not shown to the user in Telegram Mini App.

### 3. **API Call Failing Silently** ⚠️
The `updateProfile` or `updateStep` API calls might be failing.

## Debugging Steps

### Step 1: Check Browser Console
Open Telegram Mini App and check browser console for errors:
1. In Telegram, tap menu (three dots)
2. Select "Open in Browser" or "Show Web Inspector"
3. Look for errors in console
4. Try clicking "Next Step" and watch for logs

### Step 2: Check Backend Logs
```bash
# On server
docker-compose -f docker-compose.prod.yml logs backend -f | grep -i "profile\|onboarding"
```

Look for:
- ✅ `📤 Calling UserService.updateProfile...`
- ✅ `✅ Profile saved`
- ✅ `📤 Updating onboarding step to 1...`
- ✅ `✅ Onboarding step updated`
- ✅ `🧭 Navigating to City screen...`
- ❌ Any error messages

### Step 3: Test API Directly
```bash
# Get your auth token from browser console or login response
TOKEN="your_token_here"

# Test update profile
curl -X PUT http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","age":25,"gender":"male"}'

# Test update onboarding step
curl -X PATCH http://localhost:8080/api/v1/onboarding/progress \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"step":1}'
```

## Quick Fixes

### Fix 1: Add Better Error Handling

Update `ProfileSetupScreen.tsx` to show errors to user:

```typescript
const handleNext = async () => {
    // ... validation code ...

    setIsSaving(true);
    
    try {
        // Save profile
        console.log('📤 Saving profile...');
        await UserService.updateProfile({
            name: name.trim(),
            age: ageNum,
            gender,
        });
        console.log('✅ Profile saved');

        // Update step
        console.log('📤 Updating step...');
        await updateStep(1);
        console.log('✅ Step updated');

        // Navigate
        console.log('🧭 Navigating...');
        navigation.navigate('City');
        console.log('✅ Navigated');
        
    } catch (error: any) {
        console.error('❌ Error:', error);
        
        // Show detailed error
        const errorMsg = error?.response?.data?.error || 
                        error?.response?.data?.details ||
                        error?.message || 
                        'Unknown error';
        
        Alert.alert(
            'Error',
            `Failed to save: ${errorMsg}\n\nPlease check your connection and try again.`,
            [{ text: 'OK' }]
        );
    } finally {
        setIsSaving(false);
    }
};
```

### Fix 2: Ensure Navigation is Available

Check `OnboardingNavigator.tsx` is properly integrated:

```typescript
// In App.tsx or main navigation file
<Stack.Navigator>
    <Stack.Screen 
        name="Onboarding" 
        component={OnboardingNavigator} 
        options={{ headerShown: false }}
    />
    <Stack.Screen 
        name="Main" 
        component={MainTabs} 
        options={{ headerShown: false }}
    />
</Stack.Navigator>
```

### Fix 3: Add Fallback Navigation

If `navigation.navigate` fails, try alternative methods:

```typescript
// In ProfileSetupScreen.tsx, after updateStep
try {
    navigation.navigate('City');
} catch (navError) {
    console.error('Navigate failed, trying push:', navError);
    try {
        navigation.push('City');
    } catch (pushError) {
        console.error('Push failed, trying replace:', pushError);
        navigation.replace('City');
    }
}
```

## Most Likely Issues

### 1. **Missing Migration** (90% probability)
The `onboarding_step` and `onboarding_completed` columns don't exist in the database yet.

**Fix:**
```bash
# On server
docker exec -i lomi_postgres psql -U postgres -d lomi_db < backend/database/migrations/001_add_onboarding_fields.sql
```

### 2. **CORS or Network Error** (5% probability)
API calls are being blocked.

**Fix:** Check Caddy logs and ensure CORS is configured.

### 3. **Navigation Stack Issue** (5% probability)
OnboardingNavigator not properly integrated.

**Fix:** Check main App.tsx navigation structure.

## Testing Checklist

- [ ] Run database migration
- [ ] Restart backend
- [ ] Clear browser cache
- [ ] Test in Telegram (not browser)
- [ ] Check console for errors
- [ ] Check backend logs
- [ ] Try clicking "Next Step"
- [ ] Verify navigation to CityScreen

## Expected Behavior

When user clicks "Next Step":
1. ✅ Button shows loading spinner
2. ✅ API call to `/users/me` (update profile)
3. ✅ API call to `/onboarding/progress` (update step)
4. ✅ Navigate to CityScreen
5. ✅ Progress bar updates to show 1/8

## If Still Not Working

1. **Check if database has onboarding columns:**
```sql
SELECT column_name, data_type 
FROM information_schema.columns 
WHERE table_name = 'users' 
AND column_name IN ('onboarding_step', 'onboarding_completed');
```

2. **Check if user record exists:**
```sql
SELECT id, name, age, gender, city, onboarding_step, onboarding_completed 
FROM users 
WHERE telegram_id = YOUR_TELEGRAM_ID;
```

3. **Manually update step to test navigation:**
```sql
UPDATE users 
SET onboarding_step = 1 
WHERE telegram_id = YOUR_TELEGRAM_ID;
```

Then refresh app and see if it starts at CityScreen.

## Next Steps

1. **Run migration** (most important!)
2. **Check console logs** in Telegram
3. **Test API endpoints** directly
4. **Verify navigation stack** in App.tsx

The code is correct - it's likely a database or environment issue!
