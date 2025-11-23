# 🔄 Deployment Flow Explained

## How It Works

### Current Flow (Manual)

```
1. Make changes locally
   ↓
2. Push to GitHub
   git push origin main
   ↓
3. SSH to server
   ssh user@152.53.87.200
   ↓
4. Pull latest code
   cd /opt/lomi_mini
   git pull origin main
   ↓
5. Run deploy script
   ./deploy.sh
   ↓
6. Done! ✅
```

**Yes, `git pull` will override local changes on server** - that's what we want! It updates the code to match GitHub.

---

## 🚀 Two Options for Deployment

### Option 1: Manual (Current - Simple)

**On your local computer:**
```bash
# Make your code changes...

# Push to GitHub
git add .
git commit -m "Your changes"
git push origin main
```

**On your server (SSH):**
```bash
cd /opt/lomi_mini
git pull origin main    # This pulls latest from GitHub
./deploy.sh             # This rebuilds and restarts
```

**That's it!** The `git pull` will:
- ✅ Download latest code from GitHub
- ✅ Override any local changes on server (we want this!)
- ✅ Update all files to match GitHub

---

### Option 2: Automated (With Webhook - One Command)

**Setup once:**
1. Run `setup-webhook.sh` on server
2. Add webhook in GitHub settings

**Then from local computer:**
```bash
# Make your code changes...

# One command does everything:
./deploy "Your changes"
```

This script will:
1. ✅ Commit your changes
2. ✅ Push to GitHub
3. ✅ GitHub webhook automatically triggers `deploy.sh` on server
4. ✅ Server pulls latest code and deploys

**You never need to SSH to server!**

---

## 📋 What `deploy.sh` Does on Server

```bash
./deploy.sh
```

This script:
1. ✅ Pulls latest code from GitHub (`git pull`)
2. ✅ Stops old containers
3. ✅ Builds new Docker image
4. ✅ Starts containers
5. ✅ Checks health
6. ✅ Runs database migrations (if any)
7. ✅ Reloads Caddy

**Everything is automated!**

---

## 🔄 Complete Flow Diagram

### Manual Flow:
```
Local Machine          GitHub          Server
     │                   │                │
     │─── Make changes ──│                │
     │                   │                │
     │─── git push ──────>                │
     │                   │                │
     │                   │                │
     │                   │<── git pull ───│
     │                   │                │
     │                   │                │─── ./deploy.sh
     │                   │                │
     │                   │                │─── ✅ Live!
```

### Automated Flow (With Webhook):
```
Local Machine          GitHub          Server
     │                   │                │
     │─── Make changes ──│                │
     │                   │                │
     │─── ./deploy ──────>                │
     │   (commits &      │                │
     │    pushes)        │                │
     │                   │                │
     │                   │─── Webhook ────>│
     │                   │   (triggers)   │
     │                   │                │
     │                   │                │─── Auto git pull
     │                   │                │─── Auto ./deploy.sh
     │                   │                │
     │                   │                │─── ✅ Live!
```

---

## 💡 Recommended: Use the `./deploy` Script

I created a simple script that does everything:

**From your local machine:**
```bash
# Make your changes...

# One command:
./deploy "Fixed bug in likes feature"
```

This will:
1. Commit all changes
2. Push to GitHub
3. If webhook is set up → Server auto-deploys
4. If no webhook → You SSH and run `./deploy.sh`

---

## 🎯 Quick Answer

**Yes, the flow is:**
1. Push to GitHub (from local)
2. Pull on server (overrides local server files - that's good!)
3. Run `./deploy.sh` (rebuilds and restarts)

**Or use automated webhook:**
1. Run `./deploy` (from local)
2. Everything happens automatically!

---

## 📝 Example Workflow

### Making a Bug Fix:

```bash
# 1. Local: Fix the bug
nano frontend/src/screens/likesyou/LikesYouScreen.tsx
# ... make changes ...

# 2. Local: Deploy
./deploy "Fixed likes screen bug"

# 3. Done! (If webhook is set up)
# OR SSH to server and run ./deploy.sh (if no webhook)
```

---

## 🔧 Setup Webhook (One Time)

**On server:**
```bash
cd /opt/lomi_mini
sudo ./setup-webhook.sh
# Copy the webhook secret shown
```

**On GitHub:**
1. Go to: https://github.com/yohannesjx/lomi_mini/settings/hooks
2. Add webhook:
   - URL: `http://152.53.87.200:9000/webhook`
   - Secret: (from server)
   - Events: Just the push event

**Now `./deploy` from local = automatic deployment!**

---

## Summary

- **Manual**: Push → SSH → Pull → Deploy
- **Automated**: Push → Webhook → Auto-deploy
- **Both work!** Choose what's easier for you.

The `git pull` on server **does override** local files - that's exactly what we want! It ensures server always matches GitHub.

