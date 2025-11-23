#!/usr/bin/env node

/**
 * Simple Webhook Server for GitHub Deployment
 * Listens for GitHub webhook events and triggers deployment
 * 
 * Run: node webhook-server.js
 * Or: pm2 start webhook-server.js --name lomi-webhook
 */

const http = require('http');
const crypto = require('crypto');
const { exec } = require('child_process');
const fs = require('fs');
const path = require('path');

const PORT = process.env.WEBHOOK_PORT || 9000;
const SECRET = process.env.WEBHOOK_SECRET || 'your-webhook-secret-change-this';
const DEPLOY_PATH = process.env.DEPLOY_PATH || '/opt/lomi_mini';
const DEPLOY_SCRIPT = path.join(DEPLOY_PATH, 'deploy.sh');

// Simple logging
const log = (message) => {
    const timestamp = new Date().toISOString();
    console.log(`[${timestamp}] ${message}`);
    // Also log to file
    fs.appendFileSync('/var/log/lomi-webhook.log', `[${timestamp}] ${message}\n`);
};

// Verify GitHub webhook signature
function verifySignature(payload, signature) {
    if (!signature) return false;
    
    const hmac = crypto.createHmac('sha256', SECRET);
    const digest = 'sha256=' + hmac.update(payload).digest('hex');
    
    return crypto.timingSafeEqual(
        Buffer.from(signature),
        Buffer.from(digest)
    );
}

// Execute deployment
function deploy() {
    return new Promise((resolve, reject) => {
        log('🚀 Starting deployment...');
        
        exec(`cd ${DEPLOY_PATH} && ./deploy.sh`, (error, stdout, stderr) => {
            if (error) {
                log(`❌ Deployment failed: ${error.message}`);
                log(`STDERR: ${stderr}`);
                reject(error);
                return;
            }
            
            log('✅ Deployment successful!');
            log(`STDOUT: ${stdout}`);
            if (stderr) log(`STDERR: ${stderr}`);
            resolve(stdout);
        });
    });
}

// HTTP Server
const server = http.createServer((req, res) => {
    if (req.method === 'POST' && req.url === '/webhook') {
        let body = '';
        
        req.on('data', chunk => {
            body += chunk.toString();
        });
        
        req.on('end', () => {
            const signature = req.headers['x-hub-signature-256'];
            
            // Verify signature (optional but recommended)
            if (SECRET !== 'your-webhook-secret-change-this' && !verifySignature(body, signature)) {
                log('❌ Invalid webhook signature');
                res.writeHead(401, { 'Content-Type': 'application/json' });
                res.end(JSON.stringify({ error: 'Invalid signature' }));
                return;
            }
            
            try {
                const payload = JSON.parse(body);
                
                // Only deploy on push to main/master
                if (payload.ref === 'refs/heads/main' || payload.ref === 'refs/heads/master') {
                    log(`📦 Received push to ${payload.ref}`);
                    
                    // Respond immediately
                    res.writeHead(200, { 'Content-Type': 'application/json' });
                    res.end(JSON.stringify({ 
                        status: 'accepted',
                        message: 'Deployment started' 
                    }));
                    
                    // Deploy asynchronously
                    deploy().catch(err => {
                        log(`💥 Deployment error: ${err.message}`);
                    });
                } else {
                    log(`⏭️  Ignoring push to ${payload.ref}`);
                    res.writeHead(200, { 'Content-Type': 'application/json' });
                    res.end(JSON.stringify({ 
                        status: 'ignored',
                        message: 'Not main/master branch' 
                    }));
                }
            } catch (error) {
                log(`❌ Error parsing webhook: ${error.message}`);
                res.writeHead(400, { 'Content-Type': 'application/json' });
                res.end(JSON.stringify({ error: 'Invalid payload' }));
            }
        });
    } else if (req.method === 'GET' && req.url === '/health') {
        res.writeHead(200, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ status: 'ok', service: 'lomi-webhook' }));
    } else if (req.method === 'POST' && req.url === '/deploy') {
        // Manual trigger endpoint (for testing)
        log('🔧 Manual deployment triggered');
        res.writeHead(200, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ 
            status: 'accepted',
            message: 'Manual deployment started' 
        }));
        
        deploy().catch(err => {
            log(`💥 Manual deployment error: ${err.message}`);
        });
    } else {
        res.writeHead(404, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'Not found' }));
    }
});

server.listen(PORT, () => {
    log(`🎣 Webhook server listening on port ${PORT}`);
    log(`📁 Deploy path: ${DEPLOY_PATH}`);
    log(`🔐 Secret: ${SECRET.substring(0, 10)}...`);
});

// Graceful shutdown
process.on('SIGTERM', () => {
    log('🛑 Shutting down webhook server...');
    server.close(() => {
        process.exit(0);
    });
});

