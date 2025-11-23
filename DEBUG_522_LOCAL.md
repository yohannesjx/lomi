# Debug 522 Error for local.lomi

## Issue: Connection timed out (522)

This means Caddy is running but can't reach the backend.

## Quick Tests (Run on Server)

```bash
# 1. Check if backend container is running
docker-compose -f docker-compose.prod.yml ps

# 2. Check backend logs (see if it crashed)
docker-compose -f docker-compose.prod.yml logs backend --tail 50

# 3. Test backend directly (bypass Caddy)
curl http://localhost:8080/api/v1/health

# 4. Test via Caddy using IP
curl http://152.53.87.200/api/v1/health

# 5. Check if backend is listening on port 8080
sudo netstat -tulpn | grep 8080
# OR
sudo ss -tulpn | grep 8080
```

## Common Issues:

### Backend Not Running
```bash
# Restart backend
docker-compose -f docker-compose.prod.yml restart backend

# Check logs
docker-compose -f docker-compose.prod.yml logs backend
```

### Backend Crashed
```bash
# Check what went wrong
docker-compose -f docker-compose.prod.yml logs backend --tail 100

# Restart
docker-compose -f docker-compose.prod.yml up -d backend
```

### Caddyfile Not Routing Correctly
The simple Caddyfile only handles `:80` (IP-based). If you're accessing `local.lomi`, you need to either:
1. Use IP address: `http://152.53.87.200`
2. Add `local.lomi` to Caddyfile
3. Use the actual domain: `lomi.social`

