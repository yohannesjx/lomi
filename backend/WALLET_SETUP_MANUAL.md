# 🏦 WALLET SYSTEM - MANUAL SETUP GUIDE

## Your Server Path Structure
```
~/lomi_mini/
├── backend/          (This is where the wallet files are)
│   ├── go.mod
│   ├── internal/
│   │   ├── database/migrations/004_wallet_system.sql
│   │   ├── handlers/wallet_handler.go
│   │   ├── models/wallet.go
│   │   ├── repositories/wallet_repository.go
│   │   └── services/wallet_service.go
│   └── ...
```

## Manual Setup Steps

### 1. Install Go Dependency (sqlx)

Since `go` command is not in PATH, you need to either:

**Option A: Add Go to PATH**
```bash
export PATH=$PATH:/usr/local/go/bin
# Or wherever your Go is installed
```

**Option B: Use full path to Go**
```bash
/usr/local/go/bin/go get github.com/jmoiron/sqlx
cd ~/lomi_mini/backend
/usr/local/go/bin/go mod tidy
```

### 2. Run Database Migration

```bash
cd ~/lomi_mini/backend
psql -U postgres -d lomi_db -f internal/database/migrations/004_wallet_system.sql
```

**If you get "password required":**
```bash
psql -U postgres -h localhost -d lomi_db -f internal/database/migrations/004_wallet_system.sql
```

**Or if postgres user doesn't have password:**
```bash
sudo -u postgres psql -d lomi_db -f internal/database/migrations/004_wallet_system.sql
```

### 3. Verify Migration

```bash
psql -U postgres -d lomi_db -c "\dt"
```

You should see these new tables:
- `wallets`
- `wallet_transactions`
- `withdrawal_requests`
- `payout_methods`
- `coin_packages`
- `coin_purchases`

### 4. Update Your Backend Code

The wallet routes are already added to `internal/routes/routes.go`, but you need to wire up the dependencies.

**Find your `main.go` or wherever you initialize routes**, and add:

```go
import (
    "lomi-backend/internal/repositories"
    "lomi-backend/internal/services"
    "lomi-backend/internal/handlers"
)

// After you have your database connection (db *sqlx.DB):
walletRepo := repositories.NewWalletRepository(db)
walletService := services.NewWalletService(walletRepo)
walletHandler := handlers.NewWalletHandler(walletService)

// Then in routes.go, replace the placeholder:
// walletHandler := handlers.NewWalletHandler(/* inject wallet service */)
// With the actual walletHandler passed from main
```

### 5. Restart Your Backend

```bash
# Stop current backend
pkill -f "your_backend_binary"

# Rebuild and run
cd ~/lomi_mini/backend
go build -o lomi_backend cmd/main.go
./lomi_backend
```

---

## Quick Test

Once backend is running, test the wallet:

```bash
# Get coin packages (public endpoint)
curl http://localhost:8080/api/v1/wallet/coin-packages

# Get wallet balance (requires auth token)
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     http://localhost:8080/api/v1/wallet/v2/balance
```

---

## Troubleshooting

### "go: command not found"
Find where Go is installed:
```bash
which go
# or
find /usr -name "go" 2>/dev/null
```

### "psql: command not found"
Install PostgreSQL client:
```bash
apt-get update
apt-get install postgresql-client
```

### Database connection error
Check if PostgreSQL is running:
```bash
systemctl status postgresql
```

---

## What's Next?

After setup is complete:
1. ✅ Test wallet endpoints
2. ✅ Update Android app to use new wallet APIs
3. ✅ Continue with Profile & Settings features

See `WALLET_SYSTEM_SUMMARY.md` for full API documentation.
