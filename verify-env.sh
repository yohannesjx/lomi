#!/bin/bash

# Quick script to verify .env.production is correctly formatted
# Run this on your server

echo "🔍 Verifying .env.production file..."
echo ""

if [ ! -f ".env.production" ]; then
    echo "❌ Error: .env.production not found!"
    exit 1
fi

echo "✅ File exists"
echo ""

# Check for the common issue (unquoted APP_NAME)
if grep -q '^APP_NAME=Lomi Social API$' .env.production; then
    echo "⚠️  WARNING: APP_NAME is not quoted!"
    echo "   Fix it with: sed -i 's/^APP_NAME=Lomi Social API$/APP_NAME=\"Lomi Social API\"/' .env.production"
    echo ""
fi

# Check for required variables
echo "Checking required variables..."
REQUIRED_VARS=("DB_PASSWORD" "REDIS_PASSWORD" "JWT_SECRET" "TELEGRAM_BOT_TOKEN" "S3_ENDPOINT" "S3_ACCESS_KEY" "S3_SECRET_KEY")
MISSING=0

for var in "${REQUIRED_VARS[@]}"; do
    if grep -q "^${var}=" .env.production; then
        echo "  ✅ $var"
    else
        echo "  ❌ $var (MISSING)"
        MISSING=1
    fi
done

echo ""

# Try to source it (this will catch syntax errors)
echo "Testing if file can be sourced..."
if set -a && source .env.production 2>&1 && set +a; then
    echo "✅ File can be sourced without errors"
else
    echo "❌ Error sourcing file!"
    echo "Check for unquoted values with spaces"
    exit 1
fi

echo ""
echo "✅ .env.production looks good!"
echo ""
echo "You can now run: ./deploy-all.sh"

