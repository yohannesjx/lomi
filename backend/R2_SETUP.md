# Cloudflare R2 Setup Guide

## Your R2 Credentials

Your Cloudflare R2 credentials have been configured. Here's what you need to know:

### Endpoint
```
https://a53cdfc7c678dac2a028159bcd178da2.r2.cloudflarestorage.com
```

### Access Credentials
- **Access Key ID**: `c7df163be474aae5317aa530bd8448bf`
- **Secret Access Key**: `bef8caaaf0ef6577e5b4aa0b16bd2d08600d405a3b27d82ee002b1f203fe35a5`

## Environment Variables

Create a `.env` file in the `backend/` directory with these values:

```bash
# Cloudflare R2 Configuration
S3_ENDPOINT=https://a53cdfc7c678dac2a028159bcd178da2.r2.cloudflarestorage.com
S3_ACCESS_KEY=c7df163be474aae5317aa530bd8448bf
S3_SECRET_KEY=bef8caaaf0ef6577e5b4aa0b16bd2d08600d405a3b27d82ee002b1f203fe35a5
S3_USE_SSL=true
S3_REGION=auto

# Bucket Names (create these in Cloudflare R2 dashboard)
S3_BUCKET_PHOTOS=lomi-photos
S3_BUCKET_VIDEOS=lomi-videos
S3_BUCKET_GIFTS=lomi-gifts
S3_BUCKET_VERIFICATIONS=lomi-verifications
```

## Setting Up Buckets in Cloudflare R2

1. Go to [Cloudflare Dashboard](https://dash.cloudflare.com/) → **R2**
2. Create the following buckets:
   - `lomi-photos`
   - `lomi-videos`
   - `lomi-gifts`
   - `lomi-verifications`

3. For each bucket, configure CORS if needed (for direct browser uploads):
   - Go to bucket → Settings → CORS
   - Add CORS rule:
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

## Testing the Connection

You can test the R2 connection by:

1. Starting the backend server:
   ```bash
   cd backend
   go run cmd/api/main.go
   ```

2. The server should log: `✅ Connected to S3-compatible storage (R2/MinIO)`

3. Test the upload URL endpoint:
   ```bash
   curl -X GET "http://localhost:8080/api/v1/users/media/upload-url?media_type=photo" \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
   ```

## Security Notes

⚠️ **Important:**
- Never commit `.env` files to version control
- The `.gitignore` file is configured to exclude `.env` files
- Rotate your R2 API tokens periodically
- Use different credentials for development and production

## Using with Docker Compose

If using Docker Compose, update the `docker-compose.yml` environment variables:

```yaml
environment:
  S3_ENDPOINT: https://a53cdfc7c678dac2a028159bcd178da2.r2.cloudflarestorage.com
  S3_ACCESS_KEY: c7df163be474aae5317aa530bd8448bf
  S3_SECRET_KEY: bef8caaaf0ef6577e5b4aa0b16bd2d08600d405a3b27d82ee002b1f203fe35a5
  S3_USE_SSL: "true"
  S3_REGION: "auto"
```

## API Usage

### 1. Get Upload URL
```bash
GET /api/v1/users/media/upload-url?media_type=photo
```

Response:
```json
{
  "upload_url": "https://...",
  "file_key": "users/{user_id}/photos/{uuid}.jpg",
  "file_name": "{uuid}.jpg",
  "bucket": "lomi-photos",
  "expires_in": 3600,
  "method": "PUT",
  "headers": {
    "Content-Type": "image/jpeg"
  }
}
```

### 2. Upload File to R2
```bash
PUT <upload_url>
Content-Type: image/jpeg
[file binary data]
```

### 3. Register Media
```bash
POST /api/v1/users/media
{
  "media_type": "photo",
  "file_key": "users/{user_id}/photos/{uuid}.jpg",
  "display_order": 1
}
```

### 4. Get Media with Download URLs
```bash
GET /api/v1/users/{user_id}/media
```

Returns media with pre-signed download URLs (valid for 24 hours).

