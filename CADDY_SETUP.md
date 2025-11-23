# Caddy Setup Guide

## Quick Setup

### 1. Copy Caddyfile to server
```bash
sudo cp Caddyfile /etc/caddy/Caddyfile
```

### 2. Update domain names (if needed)
```bash
sudo nano /etc/caddy/Caddyfile
```

Replace `lomi.social` and `api.lomi.social` with your actual domains.

### 3. Update frontend path (if hosting static files)
If you're hosting the frontend on the server, update the path:
```caddy
root * /var/www/lomi-frontend
```

Or if using a different location:
```caddy
root * /opt/lomi_mini/frontend/build
```

### 4. Test configuration
```bash
sudo caddy validate --config /etc/caddy/Caddyfile
```

### 5. Reload Caddy
```bash
sudo systemctl reload caddy
```

## Domain Configuration

### DNS Records Required

Point these domains to your server IP:

```
A     lomi.social        → YOUR_SERVER_IP
A     api.lomi.social    → YOUR_SERVER_IP
```

### SSL Certificates

Caddy automatically:
- ✅ Obtains SSL certificates from Let's Encrypt
- ✅ Renews certificates automatically
- ✅ Redirects HTTP → HTTPS

No manual SSL configuration needed!

## Frontend Deployment Options

### Option 1: Serve from Server (Current Setup)
```bash
# Build frontend
cd frontend
npm run build

# Copy to server
scp -r build/* user@server:/var/www/lomi-frontend/
```

### Option 2: Use CDN/External Hosting
If frontend is hosted elsewhere (Vercel, Netlify, etc.), update Caddyfile:

```caddy
lomi.social {
    reverse_proxy https://your-frontend-url.com {
        header_up Host {host}
    }
}
```

### Option 3: Redirect to External URL
```caddy
lomi.social {
    redir https://your-frontend-url.com{uri} permanent
}
```

## Backend API

The API is automatically proxied to `localhost:8080` (Docker container).

No changes needed unless:
- Backend runs on different port → Update `localhost:8080`
- Need different routing → Modify reverse_proxy block

## Testing

### Test Landing Page
```bash
curl -I https://lomi.social
```

### Test API
```bash
curl https://api.lomi.social/api/v1/health
```

Should return:
```json
{"status":"ok","message":"Lomi Backend is running 🍋"}
```

## Logs

### View Caddy logs
```bash
# All logs
sudo journalctl -u caddy -f

# Frontend logs
sudo tail -f /var/log/caddy/frontend.log

# API logs
sudo tail -f /var/log/caddy/api.log
```

## Troubleshooting

### Caddy not starting
```bash
# Check configuration
sudo caddy validate --config /etc/caddy/Caddyfile

# Check status
sudo systemctl status caddy

# View errors
sudo journalctl -u caddy -n 50
```

### SSL certificate issues
```bash
# Force certificate renewal
sudo caddy reload --config /etc/caddy/Caddyfile

# Check certificate status
sudo caddy trust
```

### Backend not responding
```bash
# Check if backend is running
curl http://localhost:8080/api/v1/health

# Check Docker container
docker-compose -f docker-compose.prod.yml ps
```

## Security Notes

- Caddy automatically handles HTTPS
- Security headers are configured
- CORS is handled (also configure in backend)
- Rate limiting can be added if needed

## Performance

- Compression enabled (gzip, zstd)
- Health checks for backend
- Connection pooling
- Automatic HTTP/2

