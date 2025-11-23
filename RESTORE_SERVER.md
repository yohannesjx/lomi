# Restore Server Files from GitHub

If `/opt/lomi_mini` is empty or deleted, follow these steps to restore everything:

## Quick Restore (Run on Server)

```bash
# 1. Navigate to /opt
cd /opt

# 2. Remove old directory (if exists)
sudo rm -rf lomi_mini

# 3. Clone fresh from GitHub
sudo git clone https://github.com/yohannesjx/lomi_mini.git

# 4. Set ownership (if needed)
sudo chown -R $USER:$USER lomi_mini

# 5. Navigate to project
cd lomi_mini

# 6. Restore environment file (if you have backup)
# Copy your .env.production file back, or recreate it

# 7. Deploy everything
./deploy-all.sh
```

## Or Use the Restore Script

```bash
# Copy restore-server.sh to your server
scp restore-server.sh root@lomi.social:/root/

# SSH to server and run
ssh root@lomi.social
chmod +x restore-server.sh
./restore-server.sh
```

## What Gets Restored

✅ All backend code
✅ All frontend code  
✅ Database migrations
✅ Docker compose files
✅ Deployment scripts
✅ Configuration files

## What You Need to Restore Manually

⚠️ `.env.production` - Environment variables (passwords, tokens, etc.)
⚠️ Database data - If database was wiped, you'll need backups
⚠️ SSL certificates - If using Let's Encrypt, they should auto-renew
⚠️ Uploaded media - If R2/S3 bucket is separate, files should still be there

## After Restore

1. **Restore .env.production** (if you have it backed up)
   ```bash
   # Copy from backup
   cp /path/to/backup/.env.production /opt/lomi_mini/.env.production
   ```

2. **Run migration** (if database needs onboarding columns)
   ```bash
   cd /opt/lomi_mini
   ./run-migration.sh
   ```

3. **Deploy everything**
   ```bash
   ./deploy-all.sh
   ```

4. **Verify services**
   ```bash
   docker-compose -f docker-compose.prod.yml ps
   ```

## If You Don't Have .env.production Backup

You'll need to recreate it with:
- Database credentials
- Redis password
- JWT secret
- Telegram bot token
- S3/R2 credentials
- Frontend URL

Check your server's environment or previous deployment logs for these values.

