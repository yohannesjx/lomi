# Lomi Admin Dashboard

Beautiful, simple admin dashboard for Lomi Social.

## Tech Stack

- Next.js 14 App Router
- TypeScript
- Tailwind CSS + DaisyUI
- Dark luxury theme (black + neon lime #A7FF83)

## Setup

```bash
cd admin
npm install
npm run dev
```

Access at: http://localhost:3001/admin

## Features

1. **Dashboard** - Stats overview
2. **Users** - User management with search and ban
3. **Photo Moderation** - Approve/reject photos
4. **Reports** - Handle user reports
5. **Gifts & Coins** - Transaction log and refunds
6. **Cashouts** - Approve/reject cashout requests
7. **Broadcast** - Send push notifications
8. **Analytics** - Charts and insights
9. **Settings** - Feature toggles

## Authentication

- Login at `/login`
- JWT token stored in localStorage
- Admin role check on all routes
- Auto-redirect if not admin

## Deployment

Build for production:
```bash
npm run build
npm start
```

Deploy on same VPS as backend (port 3001).

## API Integration

All API calls use existing Go Fiber backend:
- Base URL: `https://lomi.social/api/v1`
- JWT authentication via Bearer token
- Admin endpoints: `/admin/*`

