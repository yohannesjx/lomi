# Quick Start - Cloudflare R2 Integration

## ✅ Setup Complete

Your Cloudflare R2 credentials have been configured:

- **Endpoint**: `https://a53cdfc7c678dac2a028159bcd178da2.r2.cloudflarestorage.com`
- **Access Key**: `c7df163be474aae5317aa530bd8448bf`
- **Secret Key**: `bef8caaaf0ef6577e5b4aa0b16bd2d08600d405a3b27d82ee002b1f203fe35a5`

## Next Steps

### 1. Create R2 Buckets

Go to [Cloudflare Dashboard → R2](https://dash.cloudflare.com/) and create these buckets:
- `lomi-photos`
- `lomi-videos`
- `lomi-gifts`
- `lomi-verifications`

### 2. Test the Connection

Start the backend server:
```bash
cd backend
go run cmd/api/main.go
```

You should see: `✅ Connected to S3-compatible storage (R2/MinIO)`

### 3. Test Upload URL Generation

```bash
# Get an upload URL (requires authentication)
curl -X GET "http://localhost:8080/api/v1/users/media/upload-url?media_type=photo" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## File Structure

Files are stored in R2 with this structure:
```
users/{user_id}/photos/{uuid}.jpg
users/{user_id}/videos/{uuid}.mp4
users/{user_id}/photos/{uuid}_thumb.jpg  (thumbnails)
```

## Environment Variables

All configuration is in `backend/.env` (already created).

For Docker Compose, update `docker-compose.yml` or use the production override:
```bash
docker-compose -f docker-compose.yml -f docker-compose.prod.example.yml up
```

## Security

- ✅ `.env` file is in `.gitignore` (won't be committed)
- ⚠️ Never share your R2 credentials
- 🔄 Rotate API tokens periodically

## Documentation

See `R2_SETUP.md` for detailed setup instructions.

