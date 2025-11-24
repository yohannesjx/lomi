# Admin Dashboard Deployment Guide

## Quick Start

```bash
cd admin
npm install
npm run build
npm start
```

## Environment Variables

Create `.env.local`:
```
NEXT_PUBLIC_API_URL=https://lomi.social/api/v1
```

## Development

```bash
npm run dev
```

Access at: http://localhost:3001/admin

## Production Deployment

### Option 1: PM2 (Recommended)

```bash
cd admin
npm install
npm run build
pm2 start npm --name "lomi-admin" -- start
pm2 save
```

### Option 2: Systemd Service

Create `/etc/systemd/system/lomi-admin.service`:

```ini
[Unit]
Description=Lomi Admin Dashboard
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/root/lomi_mini/admin
ExecStart=/usr/bin/npm start
Restart=always
Environment=NODE_ENV=production
Environment=NEXT_PUBLIC_API_URL=https://lomi.social/api/v1

[Install]
WantedBy=multi-user.target
```

Then:
```bash
sudo systemctl enable lomi-admin
sudo systemctl start lomi-admin
```

### Option 3: Caddy Reverse Proxy

Add to Caddyfile:
```
admin.lomi.social {
    reverse_proxy localhost:3001
}
```

## Admin Login

For now, admin login requires:
1. Get JWT token from backend (use existing auth)
2. User must have `role: "admin"` in database
3. Login page accepts token directly (development mode)

**TODO**: Create `/admin/login` endpoint in backend for proper admin authentication.

## Backend Endpoints Needed

The dashboard expects these admin endpoints:

- `GET /admin/stats` - Dashboard statistics
- `GET /users` - List all users
- `POST /admin/users/:id/ban` - Ban user
- `GET /admin/moderation/pending` - Pending photos
- `PUT /admin/moderation/:id/approve` - Approve photo
- `PUT /admin/moderation/:id/reject` - Reject photo
- `GET /admin/reports/pending` - Pending reports
- `PUT /admin/reports/:id/review` - Review report
- `GET /admin/transactions` - Transaction log
- `POST /admin/transactions/:id/refund` - Refund transaction
- `GET /admin/payouts/pending` - Pending cashouts
- `PUT /admin/payouts/:id/process` - Process cashout
- `POST /admin/broadcast` - Send broadcast
- `GET /admin/analytics` - Analytics data
- `GET /admin/settings` - Get settings
- `PUT /admin/settings` - Update settings

Most of these already exist in the backend. Check `backend/internal/routes/routes.go` for existing admin routes.

