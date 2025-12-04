#!/bin/bash

# Wallet System Setup Script
# This script sets up the wallet management system

set -e

echo "🏦 Setting up Wallet Management System..."

# Get the directory where the script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# 1. Install dependencies
echo "📦 Installing dependencies..."
cd "$SCRIPT_DIR"
go get github.com/jmoiron/sqlx

# 2. Run database migration
echo "🗄️  Running database migration..."
echo "Enter your PostgreSQL password when prompted..."
psql -U postgres -d lomi_db -f internal/database/migrations/004_wallet_system.sql

# 3. Tidy dependencies
echo "🧹 Tidying Go modules..."
go mod tidy

echo "✅ Wallet system setup complete!"
echo ""
echo "📝 Next steps:"
echo "1. Update main.go to wire up wallet dependencies"
echo "2. Test endpoints with curl or Postman"
echo "3. Update Android app to use new endpoints"
echo ""
echo "📚 See WALLET_SYSTEM_SUMMARY.md for full documentation"
