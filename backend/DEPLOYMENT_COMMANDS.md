# 🚀 Wallet System Deployment Commands

## Step 1: Pull Latest Code

```bash
cd ~/lomi_mini
git pull origin main
```

## Step 2: Run Database Migration

```bash
cd ~/lomi_mini/backend
docker exec -i lomi_postgres psql -U lomi -d lomi_db < internal/database/migrations/004_wallet_system.sql
```

**Expected Output:**
- You'll see CREATE TABLE, CREATE INDEX, CREATE TRIGGER messages
- If you see "already exists" errors, that's fine - it means the migration ran before

## Step 3: Verify Migration

```bash
docker exec lomi_postgres psql -U lomi -d lomi_db -c "\dt" | grep wallet
```

**Expected Output:**
```
 public | coin_packages       | table | lomi
 public | coin_purchases      | table | lomi
 public | payout_methods      | table | lomi
 public | wallet_transactions | table | lomi
 public | wallets             | table | lomi
 public | withdrawal_requests | table | lomi
```

## Step 4: Restart Backend

```bash
cd ~/lomi_mini
docker compose restart backend
```

**Or if you need to rebuild:**
```bash
docker compose up -d --build backend
```

## Step 5: Check Backend Logs

```bash
docker logs -f lomi_backend --tail 50
```

**Look for:**
- ✅ Connected to PostgreSQL Database (sqlx)
- 🚀 Server starting on port 8080

## Step 6: Test Wallet Endpoints

### Test 1: Get Coin Packages (Public)
```bash
curl http://localhost:8080/api/v1/wallet/coin-packages
```

### Test 2: Get Wallet Balance (Requires Auth)
First, get a JWT token:
```bash
# Replace USER_ID with actual user UUID from database
curl -X POST http://localhost:8080/api/v1/test/jwt \
  -H "Content-Type: application/json" \
  -d '{"user_id": "YOUR_USER_UUID"}'
```

Then test balance:
```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     http://localhost:8080/api/v1/wallet/v2/balance
```

### Test 3: Purchase Coins
```bash
curl -X POST http://localhost:8080/api/v1/wallet/purchase-coins \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "coins": 100,
    "amount": 0.99,
    "payment_method": "telebirr"
  }'
```

## 🔍 Troubleshooting

### If migration fails:
```bash
# Check if Docker is running
docker ps

# Check PostgreSQL is accessible
docker exec lomi_postgres psql -U lomi -d lomi_db -c "SELECT version();"
```

### If backend won't start:
```bash
# Check logs
docker logs lomi_backend --tail 100

# Check if all dependencies are running
docker compose ps
```

### If endpoints return 404:
```bash
# Verify routes are registered
docker logs lomi_backend | grep "wallet"
```

## 📋 Quick Health Check

Run all these commands to verify everything is working:

```bash
# 1. Check Docker services
docker compose ps

# 2. Check database tables
docker exec lomi_postgres psql -U lomi -d lomi_db -c "\dt" | grep wallet

# 3. Check backend is running
curl http://localhost:8080/api/v1/health

# 4. Check wallet endpoint
curl http://localhost:8080/api/v1/wallet/coin-packages
```

## ✅ Success Indicators

You'll know everything is working when:
1. ✅ All 6 wallet tables exist in database
2. ✅ Backend logs show "Connected to PostgreSQL Database (sqlx)"
3. ✅ Health endpoint returns `{"status":"ok"}`
4. ✅ Coin packages endpoint returns JSON array with packages
5. ✅ Balance endpoint returns wallet data (with valid JWT)

## 🎯 What's Next?

After successful deployment:
1. Update Android app to use new wallet endpoints
2. Test coin purchase flow end-to-end
3. Test withdrawal request flow
4. Monitor transaction logs in database

---

**Need help?** Check the logs:
```bash
docker logs lomi_backend -f
```
