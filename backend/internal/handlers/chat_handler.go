package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"lomi-backend/internal/database"
	"lomi-backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// ==================== UNIFIED CHAT HANDLER ====================
// Handles both 1-on-1 dating chat AND TikTok-style live streaming chat
// - Private chat: Direct WebSocket + PostgreSQL
// - Live chat: Redis Pub/Sub + Redis Stream + PostgreSQL (async)

var (
	redisClient *redis.Client
	ctx         = context.Background()
	rateLimiter = NewRateLimiter()
)

// InitRedis initializes Redis client for live chat
func InitRedis(addr, password string, db int) error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		PoolSize:     100,
		MinIdleConns: 10,
	})

	// Test connection
	if err := redisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis connection failed: %w", err)
	}

	log.Printf("✅ Redis connected: %s", addr)
	return nil
}

// ==================== MESSAGE TYPES ====================

type ChatMode string

const (
	ChatModePrivate ChatMode = "private"
	ChatModeLive    ChatMode = "live"
)

type WSChatMessage struct {
	Type string   `json:"type"` // "message", "typing", "read_receipt", "join", "leave", "pin", "gift", "system"
	Mode ChatMode `json:"mode"` // "private" or "live"

	// Private chat fields
	MatchID string `json:"match_id,omitempty"`

	// Live chat fields
	LiveStreamID string `json:"live_stream_id,omitempty"`
	Seq          int64  `json:"seq,omitempty"`          // Sequence number for live messages
	ViewerCount  int    `json:"viewer_count,omitempty"` // Current viewer count
	IsPinned     bool   `json:"is_pinned,omitempty"`    // Pinned message flag
	IsSystem     bool   `json:"is_system,omitempty"`    // System message flag

	// Common fields
	MessageID      string                 `json:"message_id,omitempty"`
	Content        interface{}            `json:"content,omitempty"`
	MessageType    string                 `json:"message_type,omitempty"` // "text", "photo", "video", "voice", "gift", "system"
	MediaURL       string                 `json:"media_url,omitempty"`
	GiftID         string                 `json:"gift_id,omitempty"`
	SenderID       string                 `json:"sender_id,omitempty"`
	ReceiverID     string                 `json:"receiver_id,omitempty"`
	SenderName     string                 `json:"sender_name,omitempty"`
	SenderAvatar   string                 `json:"sender_avatar,omitempty"`
	IsTyping       bool                   `json:"is_typing,omitempty"`
	DeliveryStatus string                 `json:"delivery_status,omitempty"` // "sent", "delivered", "read"
	Timestamp      string                 `json:"timestamp"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// ==================== CLIENT CONNECTION ====================

type ChatClient struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	UserName string
	Conn     *websocket.Conn
	Send     chan []byte
	Hub      *ChatHub

	// Connection context
	Mode          ChatMode
	MatchID       *uuid.UUID // For private chat
	LiveStreamID  *uuid.UUID // For live chat
	IsBroadcaster bool       // True if user owns the live stream

	// Redis subscription (for live chat)
	RedisSub *redis.PubSub

	mu sync.RWMutex
}

// ==================== HUB MANAGEMENT ====================

type ChatHub struct {
	// Private chat clients: user_id -> client
	privateClients map[uuid.UUID]*ChatClient

	// Live chat clients: live_stream_id -> map[user_id]client
	liveClients map[uuid.UUID]map[uuid.UUID]*ChatClient

	broadcast  chan []byte
	register   chan *ChatClient
	unregister chan *ChatClient

	mu sync.RWMutex
}

func NewChatHub() *ChatHub {
	return &ChatHub{
		privateClients: make(map[uuid.UUID]*ChatClient),
		liveClients:    make(map[uuid.UUID]map[uuid.UUID]*ChatClient),
		broadcast:      make(chan []byte, 1024),
		register:       make(chan *ChatClient),
		unregister:     make(chan *ChatClient),
	}
}

func (h *ChatHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			// Handle broadcast (for private chat)
			var wsMsg WSChatMessage
			if err := json.Unmarshal(message, &wsMsg); err == nil {
				h.handleBroadcast(&wsMsg, message)
			}
		}
	}
}

func (h *ChatHub) registerClient(client *ChatClient) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if client.Mode == ChatModePrivate {
		h.privateClients[client.UserID] = client

		// Update user online status
		database.DB.Model(&models.User{}).Where("id = ?", client.UserID).Updates(map[string]interface{}{
			"is_online":    true,
			"last_seen_at": time.Now(),
		})

		log.Printf("✅ Private chat client registered: user=%s", client.UserID)

	} else if client.Mode == ChatModeLive && client.LiveStreamID != nil {
		if h.liveClients[*client.LiveStreamID] == nil {
			h.liveClients[*client.LiveStreamID] = make(map[uuid.UUID]*ChatClient)
		}
		h.liveClients[*client.LiveStreamID][client.UserID] = client

		// Increment viewer count in Redis
		viewerKey := fmt.Sprintf("live:%s:viewers", client.LiveStreamID.String())
		redisClient.Incr(ctx, viewerKey)
		redisClient.Expire(ctx, viewerKey, 24*time.Hour)

		// Get current viewer count
		viewerCount, _ := redisClient.Get(ctx, viewerKey).Int()

		// Subscribe to Redis channel for live updates
		channel := fmt.Sprintf("live:%s", client.LiveStreamID.String())
		client.RedisSub = redisClient.Subscribe(ctx, channel)

		// Start Redis listener
		go client.listenRedis()

		// Broadcast join message
		joinMsg := WSChatMessage{
			Type:         "join",
			Mode:         ChatModeLive,
			LiveStreamID: client.LiveStreamID.String(),
			SenderID:     client.UserID.String(),
			SenderName:   client.UserName,
			ViewerCount:  viewerCount,
			Timestamp:    time.Now().Format(time.RFC3339),
		}
		h.publishToLive(client.LiveStreamID.String(), &joinMsg)

		log.Printf("✅ Live chat client registered: user=%s, stream=%s, viewers=%d",
			client.UserID, client.LiveStreamID, viewerCount)
	}
}

func (h *ChatHub) unregisterClient(client *ChatClient) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if client.Mode == ChatModePrivate {
		if _, ok := h.privateClients[client.UserID]; ok {
			delete(h.privateClients, client.UserID)
			close(client.Send)

			// Update user offline status
			database.DB.Model(&models.User{}).Where("id = ?", client.UserID).Updates(map[string]interface{}{
				"is_online":    false,
				"last_seen_at": time.Now(),
			})

			log.Printf("✅ Private chat client unregistered: user=%s", client.UserID)
		}

	} else if client.Mode == ChatModeLive && client.LiveStreamID != nil {
		if clients, ok := h.liveClients[*client.LiveStreamID]; ok {
			if _, exists := clients[client.UserID]; exists {
				delete(clients, client.UserID)

				// Decrement viewer count
				viewerKey := fmt.Sprintf("live:%s:viewers", client.LiveStreamID.String())
				redisClient.Decr(ctx, viewerKey)
				viewerCount, _ := redisClient.Get(ctx, viewerKey).Int()

				// Unsubscribe from Redis
				if client.RedisSub != nil {
					client.RedisSub.Close()
				}

				close(client.Send)

				// Broadcast leave message
				leaveMsg := WSChatMessage{
					Type:         "leave",
					Mode:         ChatModeLive,
					LiveStreamID: client.LiveStreamID.String(),
					SenderID:     client.UserID.String(),
					SenderName:   client.UserName,
					ViewerCount:  viewerCount,
					Timestamp:    time.Now().Format(time.RFC3339),
				}
				h.publishToLive(client.LiveStreamID.String(), &leaveMsg)

				log.Printf("✅ Live chat client unregistered: user=%s, stream=%s, viewers=%d",
					client.UserID, client.LiveStreamID, viewerCount)
			}

			// Clean up empty stream
			if len(clients) == 0 {
				delete(h.liveClients, *client.LiveStreamID)
			}
		}
	}
}

func (h *ChatHub) handleBroadcast(wsMsg *WSChatMessage, rawMsg []byte) {
	if wsMsg.Mode == ChatModePrivate {
		// Private chat: send to specific match participants
		if wsMsg.MatchID != "" {
			matchID, _ := uuid.Parse(wsMsg.MatchID)
			var match models.Match
			if err := database.DB.First(&match, "id = ?", matchID).Error; err == nil {
				h.sendToUser(match.User1ID, rawMsg)
				h.sendToUser(match.User2ID, rawMsg)
			}
		}
	}
	// Live chat messages are handled via Redis Pub/Sub, not this broadcast channel
}

func (h *ChatHub) sendToUser(userID uuid.UUID, message []byte) {
	h.mu.RLock()
	client, ok := h.privateClients[userID]
	h.mu.RUnlock()

	if ok {
		select {
		case client.Send <- message:
		default:
			// Client send buffer is full, close connection
			h.mu.Lock()
			delete(h.privateClients, userID)
			close(client.Send)
			h.mu.Unlock()
		}
	}
}

func (h *ChatHub) publishToLive(liveStreamID string, msg *WSChatMessage) {
	channel := fmt.Sprintf("live:%s", liveStreamID)
	msgBytes, _ := json.Marshal(msg)
	redisClient.Publish(ctx, channel, msgBytes)
}

var chatHub *ChatHub

func init() {
	chatHub = NewChatHub()
	go chatHub.Run()
}

// ==================== WEBSOCKET HANDLER ====================

func HandleUnifiedChat(c *websocket.Conn) {
	// Extract token from query params
	tokenString := c.Query("token")
	if tokenString == "" {
		c.WriteJSON(fiber.Map{"error": "Token required"})
		c.Close()
		return
	}

	// Parse JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-super-secret-jwt-key-change-in-production"), nil // TODO: Get from config
	})

	if err != nil || !token.Valid {
		c.WriteJSON(fiber.Map{"error": "Invalid token"})
		c.Close()
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.WriteJSON(fiber.Map{"error": "Invalid user ID"})
		c.Close()
		return
	}

	// Get user info
	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.WriteJSON(fiber.Map{"error": "User not found"})
		c.Close()
		return
	}

	// Get connection mode and context from query params
	mode := ChatMode(c.Query("mode", "private"))
	matchIDStr := c.Query("match_id")
	liveStreamIDStr := c.Query("live_stream_id")
	lastSeqStr := c.Query("last_seq", "0")

	client := &ChatClient{
		ID:       uuid.New(),
		UserID:   userID,
		UserName: user.Name,
		Conn:     c,
		Send:     make(chan []byte, 256),
		Hub:      chatHub,
		Mode:     mode,
	}

	// Setup connection context
	if mode == ChatModePrivate {
		if matchIDStr != "" {
			matchID, _ := uuid.Parse(matchIDStr)
			client.MatchID = &matchID
		}
	} else if mode == ChatModeLive {
		if liveStreamIDStr != "" {
			liveStreamID, _ := uuid.Parse(liveStreamIDStr)
			client.LiveStreamID = &liveStreamID

			// Check if user is broadcaster
			var stream models.User // TODO: Add LiveStream model
			// For now, assume broadcaster if query param is set
			client.IsBroadcaster = c.Query("is_broadcaster") == "true"
			_ = stream
		}
	}

	// Register client
	client.Hub.register <- client

	// Send missed messages on reconnection (for live chat)
	if mode == ChatModeLive && client.LiveStreamID != nil {
		lastSeq, _ := strconv.ParseInt(lastSeqStr, 10, 64)
		if lastSeq > 0 {
			go client.sendMissedMessages(lastSeq)
		}
	}

	// Start goroutines
	go client.writePump()
	go client.readPump()
}

// ==================== CLIENT READ/WRITE ====================

func (c *ChatClient) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var wsMsg WSChatMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			continue
		}

		// Route message based on mode
		if c.Mode == ChatModePrivate {
			c.handlePrivateMessage(&wsMsg)
		} else if c.Mode == ChatModeLive {
			c.handleLiveMessage(&wsMsg)
		}
	}
}

func (c *ChatClient) writePump() {
	defer c.Conn.Close()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("Write error: %v", err)
				return
			}
		}
	}
}

// ==================== PRIVATE CHAT HANDLERS ====================

func (c *ChatClient) handlePrivateMessage(wsMsg *WSChatMessage) {
	switch wsMsg.Type {
	case "message":
		c.handlePrivateChatMessage(wsMsg)
	case "typing":
		c.handleTypingIndicator(wsMsg)
	case "read_receipt":
		c.handleReadReceipt(wsMsg)
	}
}

func (c *ChatClient) handlePrivateChatMessage(wsMsg *WSChatMessage) {
	if c.MatchID == nil {
		return
	}

	matchID := *c.MatchID
	msg := models.Message{
		MatchID:     &matchID,
		SenderID:    c.UserID,
		MessageType: models.MessageType(wsMsg.MessageType),
		Content:     "",
		MediaURL:    wsMsg.MediaURL,
		IsRead:      false,
		IsLive:      false, // Private message
	}

	// Handle content
	if wsMsg.Content != nil {
		if contentStr, ok := wsMsg.Content.(string); ok {
			msg.Content = contentStr
		}
	}

	// Handle gift
	if wsMsg.GiftID != "" {
		giftID, _ := uuid.Parse(wsMsg.GiftID)
		msg.GiftID = &giftID
	}

	// Get receiver from match
	var match models.Match
	if err := database.DB.First(&match, "id = ?", matchID).Error; err == nil {
		if match.User1ID == c.UserID {
			receiverID := match.User2ID
			msg.ReceiverID = &receiverID
		} else {
			receiverID := match.User1ID
			msg.ReceiverID = &receiverID
		}

		// Save to database
		if err := database.DB.Create(&msg).Error; err == nil {
			// Update message in WS response
			wsMsg.MessageID = msg.ID.String()
			wsMsg.SenderID = c.UserID.String()
			wsMsg.ReceiverID = msg.ReceiverID.String()
			wsMsg.DeliveryStatus = "sent"
			wsMsg.Timestamp = time.Now().Format(time.RFC3339)

			// Broadcast to both users
			broadcastMsg, _ := json.Marshal(wsMsg)
			c.Hub.broadcast <- broadcastMsg
		}
	}
}

func (c *ChatClient) handleTypingIndicator(wsMsg *WSChatMessage) {
	wsMsg.SenderID = c.UserID.String()
	typingMsg, _ := json.Marshal(wsMsg)
	c.Hub.broadcast <- typingMsg
}

func (c *ChatClient) handleReadReceipt(wsMsg *WSChatMessage) {
	if c.MatchID == nil {
		return
	}

	matchID := *c.MatchID
	now := time.Now()

	database.DB.Model(&models.Message{}).
		Where("match_id = ? AND receiver_id = ? AND is_read = ?", matchID, c.UserID, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		})

	// Broadcast read receipt
	readReceipt := WSChatMessage{
		Type:           "read_receipt",
		Mode:           ChatModePrivate,
		MatchID:        wsMsg.MatchID,
		DeliveryStatus: "read",
		Timestamp:      time.Now().Format(time.RFC3339),
	}
	readBytes, _ := json.Marshal(readReceipt)
	c.Hub.broadcast <- readBytes
}

// ==================== LIVE CHAT HANDLERS ====================

func (c *ChatClient) handleLiveMessage(wsMsg *WSChatMessage) {
	if c.LiveStreamID == nil {
		return
	}

	switch wsMsg.Type {
	case "message":
		c.handleLiveChatMessage(wsMsg)
	case "gift":
		c.handleLiveGift(wsMsg)
	case "pin":
		c.handlePinMessage(wsMsg)
	case "system":
		c.handleSystemMessage(wsMsg)
	}
}

func (c *ChatClient) handleLiveChatMessage(wsMsg *WSChatMessage) {
	// Rate limiting: 5 messages per second per user
	if !rateLimiter.Allow(c.UserID.String(), 5, time.Second) {
		errorMsg := WSChatMessage{
			Type:      "error",
			Content:   "Rate limit exceeded. Please slow down.",
			Timestamp: time.Now().Format(time.RFC3339),
		}
		errorBytes, _ := json.Marshal(errorMsg)
		c.Send <- errorBytes
		return
	}

	liveStreamID := c.LiveStreamID.String()

	// Generate sequence number
	seqKey := fmt.Sprintf("live:%s:seq", liveStreamID)
	seq, err := redisClient.Incr(ctx, seqKey).Result()
	if err != nil {
		log.Printf("❌ Failed to generate sequence: %v", err)
		return
	}
	redisClient.Expire(ctx, seqKey, 24*time.Hour)

	// Prepare message
	wsMsg.Seq = seq
	wsMsg.SenderID = c.UserID.String()
	wsMsg.SenderName = c.UserName
	wsMsg.Timestamp = time.Now().Format(time.RFC3339)
	wsMsg.LiveStreamID = liveStreamID

	// Get viewer count
	viewerKey := fmt.Sprintf("live:%s:viewers", liveStreamID)
	viewerCount, _ := redisClient.Get(ctx, viewerKey).Int()
	wsMsg.ViewerCount = viewerCount

	// Publish to Redis Pub/Sub for real-time delivery
	c.Hub.publishToLive(liveStreamID, wsMsg)

	// Add to Redis Stream for persistence and replay
	streamKey := fmt.Sprintf("live:%s:history", liveStreamID)
	msgJSON, _ := json.Marshal(wsMsg)
	redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: streamKey,
		Values: map[string]interface{}{
			"seq":     seq,
			"message": string(msgJSON),
		},
	})
	redisClient.Expire(ctx, streamKey, 24*time.Hour)

	// Async save to PostgreSQL
	go c.saveLiveMessageToDB(wsMsg)
}

func (c *ChatClient) handleLiveGift(wsMsg *WSChatMessage) {
	// Similar to live message but with gift animation trigger
	wsMsg.Type = "gift"
	wsMsg.MessageType = "gift"
	c.handleLiveChatMessage(wsMsg)
}

func (c *ChatClient) handlePinMessage(wsMsg *WSChatMessage) {
	// Only broadcaster can pin messages
	if !c.IsBroadcaster {
		return
	}

	liveStreamID := c.LiveStreamID.String()

	// Store pinned message in Redis
	pinnedKey := fmt.Sprintf("live:%s:pinned", liveStreamID)
	msgJSON, _ := json.Marshal(wsMsg)
	redisClient.Set(ctx, pinnedKey, msgJSON, 24*time.Hour)

	// Broadcast pinned message
	wsMsg.Type = "pin"
	wsMsg.IsPinned = true
	wsMsg.Timestamp = time.Now().Format(time.RFC3339)
	c.Hub.publishToLive(liveStreamID, wsMsg)

	// Update in database
	if wsMsg.MessageID != "" {
		msgID, _ := uuid.Parse(wsMsg.MessageID)
		database.DB.Model(&models.Message{}).Where("id = ?", msgID).Update("pinned", true)
	}
}

func (c *ChatClient) handleSystemMessage(wsMsg *WSChatMessage) {
	// Only broadcaster can send system messages
	if !c.IsBroadcaster {
		return
	}

	wsMsg.Type = "system"
	wsMsg.MessageType = "system"
	wsMsg.IsSystem = true
	c.handleLiveChatMessage(wsMsg)
}

// ==================== REDIS LISTENER ====================

func (c *ChatClient) listenRedis() {
	if c.RedisSub == nil {
		return
	}

	ch := c.RedisSub.Channel()
	for msg := range ch {
		// Forward Redis message to WebSocket
		select {
		case c.Send <- []byte(msg.Payload):
		default:
			// Send buffer full, skip message
		}
	}
}

// ==================== MISSED MESSAGES ====================

func (c *ChatClient) sendMissedMessages(lastSeq int64) {
	if c.LiveStreamID == nil {
		return
	}

	liveStreamID := c.LiveStreamID.String()
	streamKey := fmt.Sprintf("live:%s:history", liveStreamID)

	// Read from Redis Stream
	messages, err := redisClient.XRange(ctx, streamKey, "-", "+").Result()
	if err != nil {
		log.Printf("❌ Failed to read stream: %v", err)
		return
	}

	// Send missed messages
	for _, msg := range messages {
		if seqStr, ok := msg.Values["seq"].(string); ok {
			seq, _ := strconv.ParseInt(seqStr, 10, 64)
			if seq > lastSeq {
				if msgStr, ok := msg.Values["message"].(string); ok {
					c.Send <- []byte(msgStr)
				}
			}
		}
	}

	// Also fetch from PostgreSQL for older messages
	var dbMessages []models.Message
	database.DB.Where("live_stream_id = ? AND seq > ? AND is_live = ?", c.LiveStreamID, lastSeq, true).
		Order("seq ASC").
		Limit(100).
		Find(&dbMessages)

	for _, dbMsg := range dbMessages {
		wsMsg := WSChatMessage{
			Type:         "message",
			Mode:         ChatModeLive,
			LiveStreamID: liveStreamID,
			MessageID:    dbMsg.ID.String(),
			Seq:          dbMsg.Seq,
			Content:      dbMsg.Content,
			MessageType:  string(dbMsg.MessageType),
			SenderID:     dbMsg.SenderID.String(),
			Timestamp:    dbMsg.CreatedAt.Format(time.RFC3339),
		}
		msgBytes, _ := json.Marshal(wsMsg)
		c.Send <- msgBytes
	}
}

// ==================== DATABASE PERSISTENCE ====================

func (c *ChatClient) saveLiveMessageToDB(wsMsg *WSChatMessage) {
	if c.LiveStreamID == nil {
		return
	}

	msg := models.Message{
		LiveStreamID: c.LiveStreamID,
		SenderID:     c.UserID,
		MessageType:  models.MessageType(wsMsg.MessageType),
		Content:      "",
		IsLive:       true,
		IsSystem:     wsMsg.Type == "system",
		Seq:          wsMsg.Seq,
		Pinned:       wsMsg.IsPinned,
	}

	if wsMsg.Content != nil {
		if contentStr, ok := wsMsg.Content.(string); ok {
			msg.Content = contentStr
		}
	}

	if wsMsg.GiftID != "" {
		giftID, _ := uuid.Parse(wsMsg.GiftID)
		msg.GiftID = &giftID
	}

	if wsMsg.MediaURL != "" {
		msg.MediaURL = wsMsg.MediaURL
	}

	if wsMsg.Metadata != nil {
		msg.Metadata = models.JSONMap(wsMsg.Metadata)
	}

	// Async save (non-blocking)
	if err := database.DB.Create(&msg).Error; err != nil {
		log.Printf("❌ Failed to save live message: %v", err)
	}
}

// ==================== RATE LIMITER ====================

type RateLimiter struct {
	counters map[string]*rateLimitCounter
	mu       sync.RWMutex
}

type rateLimitCounter struct {
	count     int
	resetTime time.Time
}

func NewRateLimiter() *RateLimiter {
	rl := &RateLimiter{
		counters: make(map[string]*rateLimitCounter),
	}

	// Cleanup old counters every minute
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for range ticker.C {
			rl.cleanup()
		}
	}()

	return rl
}

func (rl *RateLimiter) Allow(key string, limit int, window time.Duration) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	counter, exists := rl.counters[key]

	if !exists || now.After(counter.resetTime) {
		rl.counters[key] = &rateLimitCounter{
			count:     1,
			resetTime: now.Add(window),
		}
		return true
	}

	if counter.count >= limit {
		return false
	}

	counter.count++
	return true
}

func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for key, counter := range rl.counters {
		if now.After(counter.resetTime.Add(5 * time.Minute)) {
			delete(rl.counters, key)
		}
	}
}

// ==================== HTTP ENDPOINTS ====================

// GetLiveViewerCount returns current viewer count for a live stream
func GetLiveViewerCount(c *fiber.Ctx) error {
	liveStreamID := c.Params("id")
	viewerKey := fmt.Sprintf("live:%s:viewers", liveStreamID)

	viewerCount, err := redisClient.Get(ctx, viewerKey).Int()
	if err != nil {
		viewerCount = 0
	}

	return c.JSON(fiber.Map{
		"live_stream_id": liveStreamID,
		"viewer_count":   viewerCount,
	})
}

// GetPinnedMessage returns the pinned message for a live stream
func GetPinnedMessage(c *fiber.Ctx) error {
	liveStreamID := c.Params("id")
	pinnedKey := fmt.Sprintf("live:%s:pinned", liveStreamID)

	msgJSON, err := redisClient.Get(ctx, pinnedKey).Result()
	if err != nil {
		return c.JSON(fiber.Map{"pinned_message": nil})
	}

	var wsMsg WSChatMessage
	json.Unmarshal([]byte(msgJSON), &wsMsg)

	return c.JSON(fiber.Map{
		"pinned_message": wsMsg,
	})
}
