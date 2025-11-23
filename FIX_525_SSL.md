# Fix Error 525 - SSL Handshake Failed

## Problem
Error 525 occurs when Caddy tries to use HTTPS/SSL for IP addresses. Let's Encrypt **cannot issue SSL certificates for IP addresses**, only for domain names.

## Solution

The updated Caddyfile now has **two separate configurations**:

### 1. HTTP-only for IP addresses (`:80`)
- ✅ Works immediately
- ✅ No SSL needed
- ✅ Access via: `http://152.53.87.200`

### 2. HTTPS for domain names
- ✅ Automatic SSL certificates
- ✅ Works when DNS is configured
- ✅ Access via: `https://lomi.social`

## How to Use

### Right Now (No DNS):
```bash
# Use HTTP (not HTTPS) with IP address
http://152.53.87.200              # Frontend
http://152.53.87.200/api/v1/health  # API
```

### After DNS is Configured:
```bash
# Use HTTPS with domain names
https://lomi.social                # Frontend
https://api.lomi.social/api/v1/health  # API
```

## On Your Server

```bash
# Pull latest code
git pull origin main

# Update Caddyfile
sudo cp Caddyfile /etc/caddy/Caddyfile

# Restart Caddy
sudo systemctl restart caddy

# Test
curl http://152.53.87.200/api/v1/health
```

## Important Notes

1. **Never use HTTPS with IP addresses** - it will always fail
2. **Use HTTP (`http://`) when accessing via IP**
3. **Use HTTPS (`https://`) when accessing via domain** (after DNS is set up)

## Testing

```bash
# Test backend directly (bypasses Caddy)
curl http://localhost:8080/api/v1/health

# Test via Caddy with IP (HTTP only)
curl http://152.53.87.200/api/v1/health

# Test via domain (after DNS - HTTPS)
curl https://api.lomi.social/api/v1/health
```

