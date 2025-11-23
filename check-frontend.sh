#!/bin/bash

# Check Frontend Deployment Status
# Run this on your server

echo "🔍 Checking frontend deployment..."

# Check if directory exists
if [ ! -d "/var/www/lomi-frontend" ]; then
    echo "❌ Directory /var/www/lomi-frontend does not exist"
    echo "Creating directory..."
    sudo mkdir -p /var/www/lomi-frontend
    sudo chown -R $USER:$USER /var/www/lomi-frontend
else
    echo "✅ Directory exists"
fi

# Check if files exist
if [ -z "$(ls -A /var/www/lomi-frontend 2>/dev/null)" ]; then
    echo "❌ Directory is empty - frontend not deployed"
    echo ""
    echo "Creating temporary placeholder page..."
    sudo tee /var/www/lomi-frontend/index.html > /dev/null <<'EOF'
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Lomi Social - Coming Soon</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            display: flex;
            align-items: center;
            justify-content: center;
            min-height: 100vh;
            text-align: center;
            padding: 20px;
        }
        .container {
            max-width: 600px;
        }
        h1 {
            font-size: 3rem;
            margin-bottom: 1rem;
            text-shadow: 2px 2px 4px rgba(0,0,0,0.3);
        }
        p {
            font-size: 1.2rem;
            margin-bottom: 2rem;
            opacity: 0.9;
        }
        .status {
            background: rgba(255,255,255,0.2);
            padding: 1rem;
            border-radius: 10px;
            margin-top: 2rem;
        }
        .api-status {
            margin-top: 1rem;
            font-size: 0.9rem;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🚀 Lomi Social</h1>
        <p>Your social platform is being set up...</p>
        <div class="status">
            <p><strong>Backend API:</strong> <span id="api-status">Checking...</span></p>
            <p class="api-status">API URL: <a href="/api/v1/health" style="color: #fff; text-decoration: underline;">/api/v1/health</a></p>
        </div>
    </div>
    <script>
        fetch('/api/v1/health')
            .then(r => r.json())
            .then(data => {
                document.getElementById('api-status').textContent = '✅ Online';
                document.getElementById('api-status').style.color = '#4ade80';
            })
            .catch(() => {
                document.getElementById('api-status').textContent = '❌ Offline';
                document.getElementById('api-status').style.color = '#f87171';
            });
    </script>
</body>
</html>
EOF
    echo "✅ Placeholder page created"
    echo ""
    echo "📝 Next steps:"
    echo "1. Build your frontend: cd frontend && npm run build"
    echo "2. Deploy to server: scp -r build/* user@server:/var/www/lomi-frontend/"
    echo "3. Or use: ./deploy-frontend.sh"
else
    echo "✅ Frontend files found"
    echo "Files in /var/www/lomi-frontend:"
    ls -la /var/www/lomi-frontend | head -10
fi

echo ""
echo "✅ Check complete!"

