# Debug Photo Upload Issues

## Current Status

From logs, we can see:
- ✅ Upload URLs are being generated successfully (200 status)
- ❌ No POST requests to `/api/v1/users/media` (media record creation not happening)
- ⚠️ This suggests R2 upload is failing silently

## Flow Analysis

1. **Get Upload URL** ✅ Working
   - Frontend calls: `GET /api/v1/users/media/upload-url?media_type=photo`
   - Backend returns presigned URL
   - Status: 200 ✅

2. **Upload to R2** ❓ Unknown
   - Frontend uploads file to presigned URL using PUT
   - Should return 200/204 if successful
   - **This might be failing silently**

3. **Create Media Record** ❌ Not happening
   - Only called when user clicks "Next"
   - Requires `fileKey` to be set (only set if R2 upload succeeds)
   - **Not being called = R2 upload likely failing**

## Debugging Steps

### 1. Check Browser Console

Open browser DevTools (F12) → Console tab, then try uploading a photo. Look for:

```
📤 Uploading to R2...
📤 Upload response: { status: ..., ... }
✅ Upload to R2 successful  (or ❌ if failed)
✅ Upload completed successfully, setting fileKey: ...
```

### 2. Check Network Tab

1. Open DevTools (F12) → Network tab
2. Try uploading a photo
3. Look for:
   - Request to `/api/v1/users/media/upload-url` (should be 200)
   - Request to R2 endpoint (the presigned URL) - **Check this one!**
   - Status code should be 200 or 204
   - If 403/404/CORS error, that's the problem

### 3. Test R2 Upload Manually

```bash
# Get upload URL
TOKEN="your-jwt-token"
UPLOAD_RESPONSE=$(curl -s -X GET "http://localhost/api/v1/users/media/upload-url?media_type=photo" \
  -H "Authorization: Bearer $TOKEN")

UPLOAD_URL=$(echo $UPLOAD_RESPONSE | jq -r '.upload_url')
FILE_KEY=$(echo $UPLOAD_RESPONSE | jq -r '.file_key')

# Create test image
echo "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==" | base64 -d > /tmp/test.jpg

# Upload to R2
curl -v -X PUT "$UPLOAD_URL" \
  -H "Content-Type: image/jpeg" \
  --data-binary @/tmp/test.jpg
```

**Check the response:**
- 200/204 = Success ✅
- 403 = Permission/CORS issue ❌
- 404 = Bucket/URL issue ❌
- CORS error = CORS not configured ❌

### 4. Check R2 CORS Settings

R2 buckets need CORS configured for browser uploads. In Cloudflare R2 dashboard:

1. Go to your bucket (e.g., `lomi-photos`)
2. Settings → CORS
3. Add CORS rule:
```json
[
  {
    "AllowedOrigins": ["*"],
    "AllowedMethods": ["GET", "PUT", "POST", "DELETE", "HEAD"],
    "AllowedHeaders": ["*"],
    "ExposeHeaders": ["ETag"],
    "MaxAgeSeconds": 3600
  }
]
```

### 5. Check Backend Logs

```bash
# Watch for all upload-related logs
docker-compose -f docker-compose.prod.yml logs -f backend | grep -E "📤|📸|upload|media|R2|S3"
```

## Common Issues

### Issue 1: CORS Error
**Symptom:** Browser console shows CORS error
**Fix:** Configure CORS in R2 bucket settings (see above)

### Issue 2: 403 Forbidden
**Symptom:** Upload returns 403
**Possible causes:**
- Presigned URL expired
- Wrong bucket
- Credentials issue
- R2 endpoint unreachable

### Issue 3: 404 Not Found
**Symptom:** Upload returns 404
**Possible causes:**
- Bucket doesn't exist
- Wrong endpoint URL
- Path-style addressing issue

### Issue 4: Network Error
**Symptom:** Request fails with network error
**Possible causes:**
- R2 endpoint unreachable from browser
- Firewall blocking
- DNS issue

## Quick Test

Run the complete test script:

```bash
cd /opt/lomi_mini
chmod +x test-photo-upload-complete.sh
./test-photo-upload-complete.sh
```

This will test the entire flow and show exactly where it fails.

## Next Steps

1. **Check browser console** - Look for upload errors
2. **Check network tab** - See the actual R2 upload request
3. **Test manually** - Use curl to test R2 upload
4. **Check CORS** - Verify R2 bucket CORS settings
5. **Check logs** - Backend logs will show if media records are being created

