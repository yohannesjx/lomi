# Test Photo Upload - Curl Commands

## Step 1: Get JWT Token

First, authenticate via Telegram to get a JWT token:

```bash
# Replace <initData> with real Telegram initData
curl -X POST http://localhost/api/v1/auth/telegram \
  -H "Authorization: tma <initData>"
```

Response will contain `access_token` - save it for next steps.

## Step 2: Get Upload URL

```bash
# Replace YOUR_JWT_TOKEN with the access_token from step 1
TOKEN="YOUR_JWT_TOKEN"

curl -X GET "http://localhost/api/v1/users/media/upload-url?media_type=photo" \
  -H "Authorization: Bearer $TOKEN" \
  -v
```

**Expected Response (200):**
```json
{
  "upload_url": "https://...",
  "file_key": "users/uuid/photo/uuid.jpg",
  "file_name": "uuid.jpg",
  "bucket": "lomi-photos",
  "expires_in": 3600,
  "method": "PUT",
  "headers": {
    "Content-Type": "image/jpeg"
  }
}
```

**If this fails, check:**
- S3Client is initialized (check backend logs)
- S3 credentials are correct
- Bucket exists in R2

## Step 3: Upload Photo to R2

```bash
# Get upload_url from step 2
UPLOAD_URL="https://..."

# Create a test image file (or use existing)
echo "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==" | base64 -d > test.jpg

# Upload to R2 using PUT
curl -X PUT "$UPLOAD_URL" \
  -H "Content-Type: image/jpeg" \
  --data-binary @test.jpg \
  -v
```

**Expected Response (200):**
- Empty body with 200 status code means upload succeeded

## Step 4: Create Media Record

```bash
# Use the file_key from step 2
FILE_KEY="users/uuid/photo/uuid.jpg"

curl -X POST "http://localhost/api/v1/users/media" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "media_type": "photo",
    "file_key": "'"$FILE_KEY"'",
    "display_order": 0
  }' \
  -v
```

**Expected Response (201):**
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "media_type": "photo",
  "url": "users/uuid/photo/uuid.jpg",
  "display_order": 0,
  "is_approved": false,
  "created_at": "..."
}
```

## Complete Test Script

```bash
#!/bin/bash

# 1. Authenticate (replace with real initData)
echo "Step 1: Authenticating..."
AUTH_RESPONSE=$(curl -s -X POST http://localhost/api/v1/auth/telegram \
  -H "Authorization: tma YOUR_INIT_DATA_HERE")

TOKEN=$(echo $AUTH_RESPONSE | jq -r '.access_token')
if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
  echo "❌ Authentication failed"
  echo $AUTH_RESPONSE
  exit 1
fi
echo "✅ Authenticated, token: ${TOKEN:0:20}..."

# 2. Get upload URL
echo "Step 2: Getting upload URL..."
UPLOAD_RESPONSE=$(curl -s -X GET "http://localhost/api/v1/users/media/upload-url?media_type=photo" \
  -H "Authorization: Bearer $TOKEN")

UPLOAD_URL=$(echo $UPLOAD_RESPONSE | jq -r '.upload_url')
FILE_KEY=$(echo $UPLOAD_RESPONSE | jq -r '.file_key')

if [ "$UPLOAD_URL" = "null" ] || [ -z "$UPLOAD_URL" ]; then
  echo "❌ Failed to get upload URL"
  echo $UPLOAD_RESPONSE
  exit 1
fi
echo "✅ Got upload URL: ${UPLOAD_URL:0:50}..."
echo "   File key: $FILE_KEY"

# 3. Create test image
echo "Step 3: Creating test image..."
echo "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==" | base64 -d > test.jpg

# 4. Upload to R2
echo "Step 4: Uploading to R2..."
UPLOAD_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X PUT "$UPLOAD_URL" \
  -H "Content-Type: image/jpeg" \
  --data-binary @test.jpg)

if [ "$UPLOAD_STATUS" != "200" ]; then
  echo "❌ Upload failed with status: $UPLOAD_STATUS"
  exit 1
fi
echo "✅ Uploaded successfully"

# 5. Create media record
echo "Step 5: Creating media record..."
MEDIA_RESPONSE=$(curl -s -X POST "http://localhost/api/v1/users/media" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"media_type\": \"photo\",
    \"file_key\": \"$FILE_KEY\",
    \"display_order\": 0
  }")

MEDIA_ID=$(echo $MEDIA_RESPONSE | jq -r '.id')
if [ "$MEDIA_ID" = "null" ] || [ -z "$MEDIA_ID" ]; then
  echo "❌ Failed to create media record"
  echo $MEDIA_RESPONSE
  exit 1
fi
echo "✅ Media record created: $MEDIA_ID"

# Cleanup
rm -f test.jpg

echo ""
echo "✅ All tests passed!"
```

## Check Backend Logs

Watch backend logs in real-time:

```bash
docker-compose -f docker-compose.prod.yml logs -f backend | grep -E "📤|📸|✅|❌|S3|R2|upload|photo"
```

## Common Issues

### 1. S3Client is nil
**Error:** "S3 storage not configured"
**Fix:** Check backend logs for S3 connection errors. Verify environment variables.

### 2. Failed to generate presigned URL
**Error:** "Failed to generate upload URL"
**Possible causes:**
- Invalid S3 credentials
- Bucket doesn't exist
- Network issue reaching R2 endpoint
- Wrong endpoint URL

### 3. Upload fails (403/404)
**Error:** HTTP 403 or 404 when uploading to R2
**Possible causes:**
- Presigned URL expired
- Wrong bucket
- CORS issue
- R2 endpoint unreachable

### 4. Media record creation fails
**Error:** "Failed to create media record"
**Possible causes:**
- Database connection issue
- Invalid file_key format
- User doesn't exist

## Test S3 Connection

Run the test script:

```bash
cd /opt/lomi_mini
chmod +x test-s3-connection.sh
./test-s3-connection.sh
```

This will:
- Check backend health
- Test S3 endpoint connectivity
- List buckets (if AWS CLI installed)
- Show S3 configuration

