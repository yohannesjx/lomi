# DNS Setup for Lomi Social

## Server IP
**152.53.87.200**

## Required DNS Records

Add these A records in your DNS provider (wherever you manage lomi.social):

```
Type    Name              Value           TTL
A       @                152.53.87.200   3600
A       api              152.53.87.200   3600
```

Or if using subdomain:
```
Type    Name              Value           TTL
A       lomi.social       152.53.87.200   3600
A       api.lomi.social   152.53.87.200   3600
```

## DNS Providers

### Cloudflare (Recommended)
1. Go to Cloudflare Dashboard
2. Select your domain (lomi.social)
3. DNS → Records
4. Add A record: `@` → `152.53.87.200`
5. Add A record: `api` → `152.53.87.200`
6. Proxy status: DNS only (gray cloud) for now

### Namecheap
1. Domain List → Manage
2. Advanced DNS
3. Add A Record: `@` → `152.53.87.200`
4. Add A Record: `api` → `152.53.87.200`

### GoDaddy
1. My Products → DNS
2. Add A Record: `@` → `152.53.87.200`
3. Add A Record: `api` → `152.53.87.200`

## Verify DNS

After adding records, verify they're working:

```bash
# Check landing page
dig lomi.social +short
# Should return: 152.53.87.200

# Check API
dig api.lomi.social +short
# Should return: 152.53.87.200

# Or use nslookup
nslookup lomi.social
nslookup api.lomi.social
```

## Testing Before DNS

If DNS isn't set up yet, you can test with IP:

```bash
# Test backend directly
curl http://152.53.87.200:8080/api/v1/health

# Or use Caddyfile.test temporarily
sudo cp Caddyfile.test /etc/caddy/Caddyfile
sudo systemctl reload caddy
```

## After DNS is Configured

1. Wait for DNS propagation (5 minutes to 48 hours)
2. Verify DNS:
   ```bash
   dig lomi.social +short
   dig api.lomi.social +short
   ```
3. Use main Caddyfile:
   ```bash
   sudo cp Caddyfile /etc/caddy/Caddyfile
   sudo systemctl reload caddy
   ```
4. Test:
   ```bash
   curl https://api.lomi.social/api/v1/health
   ```

## SSL Certificates

Caddy will automatically:
- ✅ Obtain SSL certificates from Let's Encrypt
- ✅ Renew certificates automatically
- ✅ Redirect HTTP → HTTPS

No manual SSL configuration needed once DNS is working!

