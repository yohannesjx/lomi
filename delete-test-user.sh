#!/bin/bash

# Delete Test User "J" and All Associated Data
# This script safely removes a test user and all related records

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🗑️  Delete Test User Script${NC}"
echo "=================================="
echo ""

# Load environment variables
if [ -f ".env.production" ]; then
    set -a
    source .env.production
    set +a
fi

# Find user "J" or "J's Fashion"
echo -e "${YELLOW}Step 1: Finding test user...${NC}"
USER_INFO=$(docker-compose -f docker-compose.prod.yml --env-file .env.production exec -T postgres psql -U "${DB_USER:-lomi}" -d "${DB_NAME:-lomi_db}" -t -c "
SELECT id, name, telegram_id, created_at 
FROM users 
WHERE name ILIKE '%J%' OR telegram_username ILIKE '%J%'
ORDER BY created_at DESC
LIMIT 5;
" 2>/dev/null | tr -d ' \r\n' || echo "")

if [ -z "$USER_INFO" ]; then
    echo -e "${RED}❌ No user found with name containing 'J'${NC}"
    exit 1
fi

echo "Found users:"
docker-compose -f docker-compose.prod.yml --env-file .env.production exec -T postgres psql -U "${DB_USER:-lomi}" -d "${DB_NAME:-lomi_db}" -c "
SELECT id, name, telegram_id, telegram_username, created_at 
FROM users 
WHERE name ILIKE '%J%' OR telegram_username ILIKE '%J%'
ORDER BY created_at DESC;
" 2>/dev/null

echo ""
read -p "Enter the user ID (UUID) to delete (or 'all' to delete all J users): " USER_ID

if [ -z "$USER_ID" ]; then
    echo -e "${RED}❌ User ID required${NC}"
    exit 1
fi

# Confirm deletion
echo ""
echo -e "${RED}⚠️  WARNING: This will delete the user and ALL associated data:${NC}"
echo "   - Media (photos/videos)"
echo "   - Swipes"
echo "   - Matches"
echo "   - Messages"
echo "   - Reports"
echo "   - Blocks"
echo "   - Verifications"
echo "   - Transactions"
echo "   - Payouts"
echo "   - All other related records"
echo ""
read -p "Are you sure you want to continue? (type 'yes' to confirm): " CONFIRM

if [ "$CONFIRM" != "yes" ]; then
    echo -e "${YELLOW}Cancelled.${NC}"
    exit 0
fi

# Delete user(s)
if [ "$USER_ID" = "all" ]; then
    echo ""
    echo -e "${YELLOW}Deleting all users with 'J' in name...${NC}"
    
    docker-compose -f docker-compose.prod.yml --env-file .env.production exec -T postgres psql -U "${DB_USER:-lomi}" -d "${DB_NAME:-lomi_db}" <<EOF
-- Delete all users with 'J' in name
-- CASCADE will automatically delete related records
DELETE FROM users 
WHERE name ILIKE '%J%' OR telegram_username ILIKE '%J%';

-- Show what was deleted
SELECT 'Deleted users with J in name' as result;
EOF

else
    echo ""
    echo -e "${YELLOW}Deleting user: $USER_ID${NC}"
    
    # First show what will be deleted
    echo ""
    echo -e "${BLUE}Records that will be deleted:${NC}"
    docker-compose -f docker-compose.prod.yml --env-file .env.production exec -T postgres psql -U "${DB_USER:-lomi}" -d "${DB_NAME:-lomi_db}" <<EOF
SELECT 'Media' as table_name, COUNT(*) as count FROM media WHERE user_id = '$USER_ID'
UNION ALL
SELECT 'Swipes', COUNT(*) FROM swipes WHERE user_id = '$USER_ID' OR swiped_user_id = '$USER_ID'
UNION ALL
SELECT 'Matches', COUNT(*) FROM matches WHERE user1_id = '$USER_ID' OR user2_id = '$USER_ID'
UNION ALL
SELECT 'Messages', COUNT(*) FROM messages WHERE sender_id = '$USER_ID' OR recipient_id = '$USER_ID'
UNION ALL
SELECT 'Reports', COUNT(*) FROM reports WHERE reporter_id = '$USER_ID' OR reported_user_id = '$USER_ID'
UNION ALL
SELECT 'Blocks', COUNT(*) FROM blocks WHERE blocker_id = '$USER_ID' OR blocked_user_id = '$USER_ID'
UNION ALL
SELECT 'Verifications', COUNT(*) FROM verifications WHERE user_id = '$USER_ID'
UNION ALL
SELECT 'Transactions', COUNT(*) FROM coin_transactions WHERE user_id = '$USER_ID'
UNION ALL
SELECT 'Payouts', COUNT(*) FROM payouts WHERE user_id = '$USER_ID';
EOF

    echo ""
    read -p "Continue with deletion? (type 'yes'): " FINAL_CONFIRM
    
    if [ "$FINAL_CONFIRM" != "yes" ]; then
        echo -e "${YELLOW}Cancelled.${NC}"
        exit 0
    fi
    
    # Delete the user (CASCADE will handle related records)
    docker-compose -f docker-compose.prod.yml --env-file .env.production exec -T postgres psql -U "${DB_USER:-lomi}" -d "${DB_NAME:-lomi_db}" <<EOF
-- Delete user (CASCADE will delete all related records)
DELETE FROM users WHERE id = '$USER_ID';

-- Verify deletion
SELECT CASE 
    WHEN EXISTS (SELECT 1 FROM users WHERE id = '$USER_ID') 
    THEN 'User still exists - deletion failed'
    ELSE 'User deleted successfully'
END as result;
EOF

fi

echo ""
echo -e "${GREEN}✅ Deletion complete!${NC}"
echo ""
echo "Remaining users:"
docker-compose -f docker-compose.prod.yml --env-file .env.production exec -T postgres psql -U "${DB_USER:-lomi}" -d "${DB_NAME:-lomi_db}" -c "
SELECT id, name, telegram_id, created_at 
FROM users 
ORDER BY created_at DESC
LIMIT 10;
" 2>/dev/null

echo ""
echo -e "${GREEN}You can now test the photo moderation system from the beginning!${NC}"

