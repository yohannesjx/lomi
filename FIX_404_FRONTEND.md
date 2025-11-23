# Fix 404 Error - Frontend Not Deployed

## Problem
Getting 404 errors because frontend files aren't deployed to `/var/www/lomi-frontend`.

## Quick Fix (On Server)

Run this on your server to create a placeholder page:

```bash
# Pull latest code
git pull origin main

# Run the check script
chmod +x check-frontend.sh
./check-frontend.sh
```

This will:
1. Create `/var/www/lomi-frontend` if it doesn't exist
2. Create a temporary "Coming Soon" page
3. Show you what's in the directory

## Deploy Frontend (From Local Machine)

### Option 1: Use Deployment Script

```bash
# On your local machine
chmod +x deploy-frontend.sh

# Set server details (if needed)
export SERVER_USER=root
export SERVER_HOST=152.53.87.200

# Deploy
./deploy-frontend.sh
```

### Option 2: Manual Deployment

```bash
# 1. Build frontend locally
cd frontend
npm install
npm run build  # or npx expo export:web for Expo

# 2. Upload to server
scp -r build/* root@152.53.87.200:/var/www/lomi-frontend/

# Or use rsync
rsync -avz --delete build/ root@152.53.87.200:/var/www/lomi-frontend/
```

### Option 3: Build on Server

```bash
# On your server
cd ~/lomi_mini/frontend
npm install
npm run build

# Copy to web directory
sudo mkdir -p /var/www/lomi-frontend
sudo cp -r build/* /var/www/lomi-frontend/
sudo chown -R www-data:www-data /var/www/lomi-frontend
```

## Verify Deployment

```bash
# On server, check if files exist
ls -la /var/www/lomi-frontend/

# Should see:
# - index.html
# - static/ (or assets/)
# - Other build files
```

## Test

```bash
# Test locally on server
curl http://localhost/api/v1/health  # Should work
curl http://localhost/                # Should return HTML

# Test from browser
http://152.53.87.200
```

## Frontend Build Commands

### Expo/React Native
```bash
cd frontend
npx expo export:web
# or
npm run build
```

### React/Next.js
```bash
cd frontend
npm run build
```

### Vite
```bash
cd frontend
npm run build
```

## Troubleshooting

### Still getting 404?

1. **Check Caddyfile path:**
   ```bash
   sudo cat /etc/caddy/Caddyfile | grep root
   # Should show: root * /var/www/lomi-frontend
   ```

2. **Check permissions:**
   ```bash
   sudo ls -la /var/www/lomi-frontend
   sudo chown -R www-data:www-data /var/www/lomi-frontend
   ```

3. **Check Caddy logs:**
   ```bash
   sudo journalctl -u caddy -n 50
   ```

4. **Reload Caddy:**
   ```bash
   sudo systemctl reload caddy
   ```

## Next Steps

After deploying the frontend:
1. Update API URL in frontend config
2. Test the full flow
3. Set up DNS for domain access

