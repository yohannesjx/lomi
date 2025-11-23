# 💚 Lomi Social - Find Your Lomi in Ethiopia

**Tagline:** "Find your Lomi (your lemon, your love) in Ethiopia"

**Domain:** lomi.social  
**Telegram Bot:** @lomi_social_bot

---

## 🎯 Project Overview

Lomi Social is a premium dating + social app built specifically for Ethiopian culture. It combines the best of modern dating apps with Ethiopian traditions - coffee ceremonies, Habesha music, Amharic poetry, and cultural respect.

### Core Features
- 🔥 **Vertical Explore Feed** - TikTok-style short videos + photos
- 💚 **Smart Matching** - AI-powered daily matches based on location, religion, language
- 💬 **Rich Chat** - Text, voice, photos, videos, stickers, animated gifts
- 🎁 **Gifting Economy** - Buy & send Ethiopian-themed animated gifts
- 💰 **Cashout System** - Top users can convert gifts to real Birr
- ✅ **Lomi Verified** - Optional selfie + ID verification
- 🌍 **Bilingual** - Full Amharic + English support

---

## 🛠 Tech Stack

### Frontend
- **React Native** (Expo) + TypeScript
- **Telegram WebApp SDK** for Mini App integration
- **React Navigation** for routing
- **Zustand** for state management
- **React Query** for data fetching
- **Socket.io Client** for real-time features

### Backend
- **Golang** (Fiber framework)
- **GORM** for database ORM
- **WebSocket** + **Redis Pub/Sub** for real-time
- **PostgreSQL** for primary database
- **Redis** for caching & sessions
- **Cloudflare R2** (S3-compatible) for media storage
- **JWT** for authentication

### Infrastructure
- **Docker** + **Docker Compose** for local dev
- **Nginx** for reverse proxy (production)
- **GitHub Actions** for CI/CD

### Payment Integrations
- Telebirr
- CBE Birr
- HelloCash
- Amole

---

## 📁 Project Structure

```
lomi_mini/
├── frontend/              # React Native (Expo) app
├── backend/               # Golang API server
├── docker/                # Docker configurations
├── docs/                  # Documentation & diagrams
├── scripts/               # Utility scripts
├── docker-compose.yml     # Local development setup
└── README.md
```

---

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Node.js 18+ (for frontend)
- Go 1.21+ (for backend development)
- Expo CLI (`npm install -g expo-cli`)

### Local Development

1. **Clone the repository**
```bash
git clone <repo-url>
cd lomi_mini
```

2. **Start all services**
```bash
docker-compose up -d
```

This will start:
- PostgreSQL (port 5432)
- Redis (port 6379)
- MinIO (port 9000, console: 9001)
- Backend API (port 8080)
- Frontend Metro bundler (port 19000)

3. **Access the services**
- Backend API: http://localhost:8080
- MinIO Console: http://localhost:9001 (admin/minioadmin)
- PostgreSQL: localhost:5432 (lomi/lomi123/lomi_db)

4. **Run frontend on Expo**
```bash
cd frontend
npm install
npm start
```

---

## 🎨 Design System

### Color Palette
- **Primary (Neon Lime):** #A7FF83
- **Background:** #000000 (Pure Black)
- **Secondary (Purple):** #B794F6
- **Accent (Pink):** #FF6B9D
- **Text Primary:** #FFFFFF
- **Text Secondary:** #A0A0A0
- **Success:** #10B981
- **Error:** #EF4444
- **Warning:** #F59E0B

### Typography
- **Primary Font:** SF Pro Display (iOS) / Roboto (Android)
- **Amharic Font:** Noto Sans Ethiopic

---

## 📊 Database Schema

See `backend/database/schema.sql` for complete PostgreSQL schema.

Key tables:
- `users` - User profiles & authentication
- `photos` - User photos & videos
- `matches` - Match relationships
- `messages` - Chat messages
- `gifts` - Gift catalog
- `gift_transactions` - Gift sending history
- `coins_transactions` - Coin purchases & usage
- `payouts` - Cashout requests & history
- `verifications` - ID verification data

---

## 🔌 API Endpoints

### Authentication
- `POST /api/v1/auth/telegram` - Telegram login
- `POST /api/v1/auth/refresh` - Refresh JWT token

### User Profile
- `GET /api/v1/users/me` - Get current user
- `PUT /api/v1/users/me` - Update profile
- `POST /api/v1/users/photos` - Upload photo/video
- `DELETE /api/v1/users/photos/:id` - Delete photo

### Discovery
- `GET /api/v1/discover/feed` - Explore feed
- `GET /api/v1/discover/swipe` - Get swipe cards
- `POST /api/v1/discover/swipe` - Swipe action (like/pass)

### Matching
- `GET /api/v1/matches` - Get all matches
- `GET /api/v1/matches/:id` - Get match details

### Chat
- `GET /api/v1/chats` - Get all conversations
- `GET /api/v1/chats/:id/messages` - Get messages
- `POST /api/v1/chats/:id/messages` - Send message
- `WS /api/v1/ws` - WebSocket connection

### Gifts & Coins
- `GET /api/v1/gifts` - Get gift catalog
- `POST /api/v1/gifts/send` - Send gift
- `POST /api/v1/coins/purchase` - Buy coins
- `GET /api/v1/coins/balance` - Get coin balance
- `GET /api/v1/coins/earn/channels` - List reward channels
- `POST /api/v1/coins/earn/claim` - Claim channel subscription reward

### Payouts
- `GET /api/v1/payouts/balance` - Get cashout balance
- `POST /api/v1/payouts/request` - Request payout
- `GET /api/v1/payouts/history` - Payout history

---

## 🎁 Gift Economy

### How It Works
1. Users buy **Lomi Coins** using Ethiopian payment methods OR earn them by subscribing to partner Telegram channels.
2. Users send **animated gifts** to matches (costs coins)
3. Recipients accumulate gift value in their balance
4. Top users can **cashout** to real Birr (platform takes 20-30% fee)
5. Payouts processed every Monday via Telebirr/bank transfer

### Gift Catalog
- ☕ Bunna Ceremony (100 coins)
- 🍲 Doro Wot Plate (150 coins)
- 🌹 Red Rose + Tej (200 coins)
- 👗 Habesha Dress (500 coins - avatar wears for 7 days)
- 🔑 Golden "Ye Fikir Key" (1000 coins)

---

## 🔐 Security & Privacy

- End-to-end encryption for messages (planned)
- Photo moderation using AI
- Report & block functionality
- ID verification for "Lomi Verified" badge
- Rate limiting on all endpoints
- GDPR-compliant data handling

---

## 📱 Telegram Mini App Integration

The app uses Telegram WebApp SDK to:
- Authenticate users instantly via Telegram
- Access user's Telegram profile
- Send notifications via Telegram bot
- Enable sharing to Telegram chats

---

## 🌍 Localization

Full bilingual support:
- **English** (default)
- **Amharic** (አማርኛ)

All UI strings are externalized in `frontend/src/locales/`

---

## 📈 Analytics & Monitoring

- User behavior tracking
- Match success rates
- Gift economy metrics
- Payment conversion funnels
- Real-time active users

---

## 🚢 Deployment

### Production Stack
- **Frontend:** Expo EAS Build → App Stores
- **Backend:** Docker containers on VPS/Cloud
- **Database:** Managed PostgreSQL (AWS RDS / DigitalOcean)
- **Redis:** Managed Redis (AWS ElastiCache / DigitalOcean)
- **Storage:** Cloudflare R2 (Production) / MinIO (Local)
- **CDN:** Cloudflare

---

## 📄 License

Proprietary - All rights reserved

---

## 👥 Team

Built with 💚 for Ethiopia

---

**Let's find your Lomi! 🍋💚**
