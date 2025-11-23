#!/bin/bash

# Test script to verify workers can download images from R2
# Run this from inside a worker container

echo "🧪 Testing R2 download access from worker container..."

# Get a test URL from Redis queue (if any job exists)
# Or manually test with a known presigned URL

if [ -z "$1" ]; then
    echo "Usage: $0 <presigned_r2_url>"
    echo "Example: $0 'https://xxx.r2.dev/bucket/key?X-Amz-Algorithm=...'"
    exit 1
fi

TEST_URL="$1"

echo "Testing URL: ${TEST_URL:0:80}..."

# Test download
response=$(curl -s -o /tmp/test_image.jpg -w "%{http_code}" "$TEST_URL")

if [ "$response" == "200" ]; then
    file_size=$(stat -f%z /tmp/test_image.jpg 2>/dev/null || stat -c%s /tmp/test_image.jpg 2>/dev/null)
    echo "✅ SUCCESS: Downloaded image (HTTP $response, Size: $file_size bytes)"
    
    # Try to open with PIL/Pillow to verify it's a valid image
    python3 -c "
from PIL import Image
import sys
try:
    img = Image.open('/tmp/test_image.jpg')
    print(f'✅ Image is valid: {img.size[0]}x{img.size[1]}, format: {img.format}')
    sys.exit(0)
except Exception as e:
    print(f'❌ Image is invalid: {e}')
    sys.exit(1)
"
else
    echo "❌ FAILED: HTTP $response"
    echo "Response body:"
    cat /tmp/test_image.jpg
    exit 1
fi

