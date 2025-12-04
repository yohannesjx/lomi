#!/bin/bash
# Quick deployment commands for your server

# Run these commands on your VPS:

# 1. Pull latest code
cd /opt/lomi-backend
git pull origin main

# 2. Build backend
cd backend
go build -o lomi-backend cmd/api/main.go

# 3. Restart service
sudo systemctl restart lomi-backend

# 4. Check status
sudo systemctl status lomi-backend

# 5. View logs
sudo journalctl -u lomi-backend -f
