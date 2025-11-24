# Debug: Photo Moderation Not Working

## Quick Checks

### 1. Is Frontend Updated?

**Check if frontend has the new code:**
- Frontend needs to be rebuilt with the new `uploadComplete` method
- Check browser console for errors when uploading

**If frontend is not updated:**
```bash
# Rebuild frontend
cd frontend
npm run build  # or your build command
# Then deploy/restart frontend
```

### 2. Check Backend Logs

**On your server, check if upload-complete was called:**
```bash
# Check recent backend logs
docker-compose -f docker-compose.prod.yml --env-file .env.production logs backend --tail 100 | grep -i "upload-complete\|batch_id"

# Should see: "✅ Upload complete: batch_id=..."
```

### 3. Check Database

**See if media records were created:**
```bash
# Check recent media records
docker-compose -f docker-compose.prod.yml --env-file .env.production exec postgres psql -U lomi -d lomi_db -c "
SELECT id, user_id, moderation_status, batch_id, created_at 
FROM media 
ORDER BY created_at DESC 
LIMIT 5;
"
```

### 4. Check Queue

**See if jobs were enqueued:**
```bash
# Check queue length
docker-compose -f docker-compose.prod.yml --env-file .env.production exec redis redis-cli -a "${REDIS_PASSWORD}" LLEN photo_moderation_queue

# Should be > 0 if jobs are waiting
```

### 5. Check Workers

**See if workers are processing:**
```bash
# Check worker logs
docker-compose -f docker-compose.prod.yml --env-file .env.production logs moderator-worker --tail 50 | grep -i "job\|batch\|processing"
```

---

## Common Issues

### Issue 1: Frontend Not Calling upload-complete

**Symptom**: No logs in backend, no queue activity

**Solution**: 
- Rebuild frontend with new code
- Check browser console for errors
- Verify `UserService.uploadComplete` exists

### Issue 2: Upload-complete Called But No Queue Activity

**Symptom**: Backend logs show "Upload complete" but queue is empty

**Check**:
```bash
# Check for enqueue errors
docker-compose -f docker-compose.prod.yml --env-file .env.production logs backend | grep -i "enqueue\|redis"
```

**Possible causes**:
- Redis connection issue
- Queue enqueue failed (check logs)

### Issue 3: Jobs in Queue But Workers Not Processing

**Symptom**: Queue length > 0 but workers idle

**Check**:
```bash
# Check worker logs for errors
docker-compose -f docker-compose.prod.yml --env-file .env.production logs moderator-worker | grep -i "error\|failed\|exception"
```

**Possible causes**:
- Workers crashed
- Redis connection issue
- CompreFace not accessible

### Issue 4: Photos Processed But No Notifications

**Symptom**: Photos moderated but no Telegram message

**Check**:
```bash
# Check subscriber logs
docker-compose -f docker-compose.prod.yml --env-file .env.production logs backend | grep -i "push\|notification\|telegram"
```

**Possible causes**:
- Notification service not initialized
- Telegram bot token missing
- User's telegram_id is 0 or invalid

---

## Step-by-Step Debugging

Run this script to check everything:
```bash
./debug-moderation-flow.sh
```

Or check manually:

**Step 1: Verify upload-complete was called**
```bash
docker-compose -f docker-compose.prod.yml --env-file .env.production logs backend --tail 200 | grep "Upload complete"
```

**Step 2: Check if media records exist**
```bash
docker-compose -f docker-compose.prod.yml --env-file .env.production exec postgres psql -U lomi -d lomi_db -c "SELECT COUNT(*) FROM media WHERE moderation_status = 'pending';"
```

**Step 3: Check queue**
```bash
docker-compose -f docker-compose.prod.yml --env-file .env.production exec redis redis-cli -a "${REDIS_PASSWORD}" LLEN photo_moderation_queue
```

**Step 4: Check workers**
```bash
docker-compose -f docker-compose.prod.yml --env-file .env.production ps moderator-worker
docker-compose -f docker-compose.prod.yml --env-file .env.production logs moderator-worker --tail 20
```

---

## Quick Fix Commands

**If queue has jobs but workers not processing:**
```bash
# Restart workers
docker-compose -f docker-compose.prod.yml --env-file .env.production restart moderator-worker
```

**If Redis connection issue:**
```bash
# Test Redis
docker-compose -f docker-compose.prod.yml --env-file .env.production exec redis redis-cli -a "${REDIS_PASSWORD}" PING
```

**If CompreFace issue:**
```bash
# Check CompreFace
curl http://localhost:8000/api/v1/health
docker-compose -f docker-compose.prod.yml --env-file .env.production logs compreface --tail 20
```

---

## Expected Flow

1. ✅ User uploads photos → Frontend calls `upload-complete`
2. ✅ Backend logs: "✅ Upload complete: batch_id=..."
3. ✅ Media records created with `moderation_status='pending'`
4. ✅ Job enqueued to Redis queue
5. ✅ Worker picks up job → Logs: "📥 Received job: batch_id=..."
6. ✅ Worker processes → Logs: "✅ Completed batch: batch_id=..."
7. ✅ Subscriber receives result → Logs: "📥 Received moderation result: batch_id=..."
8. ✅ Database updated → Logs: "✅ Updated media record: media_id=..."
9. ✅ Push notification sent → Logs: "✅ Sent push notification: user_id=..."

**If any step is missing, that's where the issue is!**

