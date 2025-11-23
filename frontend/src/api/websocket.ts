import { useAuthStore } from '../store/authStore';

export interface WSMessage {
    type: 'message' | 'typing' | 'read_receipt' | 'online_status' | 'delivery_status';
    match_id?: string;
    message_id?: string;
    content?: string;
    message_type?: 'text' | 'photo' | 'video' | 'voice' | 'gift';
    media_url?: string;
    gift_id?: string;
    sender_id?: string;
    receiver_id?: string;
    is_typing?: boolean;
    delivery_status?: 'sent' | 'delivered' | 'read';
    timestamp?: string;
}

class WebSocketService {
    private ws: WebSocket | null = null;
    private reconnectAttempts = 0;
    private maxReconnectAttempts = 5;
    private reconnectDelay = 1000;
    private listeners: Map<string, Set<(data: WSMessage) => void>> = new Map();
    private isConnecting = false;

    connect(token: string) {
        if (this.ws?.readyState === WebSocket.OPEN || this.isConnecting) {
            return;
        }

        this.isConnecting = true;
        // Get API URL - handle both http and https
        let apiUrl = process.env.EXPO_PUBLIC_API_URL || 'http://localhost:8080';
        
        // Remove /api/v1 if present
        apiUrl = apiUrl.replace('/api/v1', '');
        
        // Convert to WebSocket URL
        const wsProtocol = apiUrl.startsWith('https') ? 'wss' : 'ws';
        const wsBase = apiUrl.replace(/^https?:\/\//, '');
        const wsUrl = `${wsProtocol}://${wsBase}/api/v1/ws?token=${token}`;

        try {
            this.ws = new WebSocket(wsUrl);

            this.ws.onopen = () => {
                console.log('WebSocket connected');
                this.isConnecting = false;
                this.reconnectAttempts = 0;
            };

            this.ws.onmessage = (event) => {
                try {
                    const message: WSMessage = JSON.parse(event.data);
                    this.handleMessage(message);
                } catch (error) {
                    console.error('Error parsing WebSocket message:', error);
                }
            };

            this.ws.onerror = (error) => {
                console.error('WebSocket error:', error);
                this.isConnecting = false;
            };

            this.ws.onclose = () => {
                console.log('WebSocket disconnected');
                this.isConnecting = false;
                this.reconnect();
            };
        } catch (error) {
            console.error('Failed to connect WebSocket:', error);
            this.isConnecting = false;
        }
    }

    private reconnect() {
        if (this.reconnectAttempts >= this.maxReconnectAttempts) {
            console.error('Max reconnection attempts reached');
            return;
        }

        this.reconnectAttempts++;
        const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1);

        setTimeout(() => {
            const { accessToken } = useAuthStore.getState();
            if (accessToken) {
                this.connect(accessToken);
            }
        }, delay);
    }

    private handleMessage(message: WSMessage) {
        // Emit to all listeners for this message type
        const listeners = this.listeners.get(message.type);
        if (listeners) {
            listeners.forEach(listener => listener(message));
        }

        // Also emit to match-specific listeners
        if (message.match_id) {
            const matchListeners = this.listeners.get(`match:${message.match_id}`);
            if (matchListeners) {
                matchListeners.forEach(listener => listener(message));
            }
        }
    }

    send(message: WSMessage) {
        if (this.ws?.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify(message));
        } else {
            console.warn('WebSocket not connected, message not sent');
        }
    }

    sendMessage(matchId: string, content: string, messageType: 'text' | 'photo' | 'video' | 'voice' = 'text', mediaUrl?: string, giftId?: string) {
        this.send({
            type: 'message',
            match_id: matchId,
            content,
            message_type: messageType,
            media_url: mediaUrl,
            gift_id: giftId,
        });
    }

    sendTyping(matchId: string, isTyping: boolean) {
        this.send({
            type: 'typing',
            match_id: matchId,
            is_typing: isTyping,
        });
    }

    sendReadReceipt(matchId: string) {
        this.send({
            type: 'read_receipt',
            match_id: matchId,
        });
    }

    on(event: string, callback: (data: WSMessage) => void) {
        if (!this.listeners.has(event)) {
            this.listeners.set(event, new Set());
        }
        this.listeners.get(event)!.add(callback);

        // Return unsubscribe function
        return () => {
            const listeners = this.listeners.get(event);
            if (listeners) {
                listeners.delete(callback);
            }
        };
    }

    off(event: string, callback: (data: WSMessage) => void) {
        const listeners = this.listeners.get(event);
        if (listeners) {
            listeners.delete(callback);
        }
    }

    disconnect() {
        if (this.ws) {
            this.ws.close();
            this.ws = null;
        }
        this.listeners.clear();
    }
}

export const wsService = new WebSocketService();

