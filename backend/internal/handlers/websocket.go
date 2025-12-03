package handlers

import (
	"encoding/json"
	"log"
	"lomi-backend/internal/database"
	"lomi-backend/internal/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// WebSocket message types
type WSMessage struct {
	Type           string      `json:"type"` // "message", "typing", "read_receipt", "online_status", "delivery_status"
	MatchID        string      `json:"match_id,omitempty"`
	MessageID      string      `json:"message_id,omitempty"`
	Content        interface{} `json:"content,omitempty"`
	MessageType    string      `json:"message_type,omitempty"` // "text", "photo", "video", "voice", "gift"
	MediaURL       string      `json:"media_url,omitempty"`
	GiftID         string      `json:"gift_id,omitempty"`
	SenderID       string      `json:"sender_id,omitempty"`
	ReceiverID     string      `json:"receiver_id,omitempty"`
	IsTyping       bool        `json:"is_typing,omitempty"`
	DeliveryStatus string      `json:"delivery_status,omitempty"` // "sent", "delivered", "read"
	Timestamp      string      `json:"timestamp"`
}

// Client represents a WebSocket connection
type Client struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Conn   *websocket.Conn
	Send   chan []byte
	Hub    *Hub
}

// Hub manages WebSocket connections
type Hub struct {
	clients    map[uuid.UUID]*Client // user_id -> client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uuid.UUID]*Client),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.UserID] = client
			// Update user online status
			database.DB.Model(&models.User{}).Where("id = ?", client.UserID).Updates(map[string]interface{}{
				"is_online":    true,
				"last_seen_at": time.Now(),
			})

		case client := <-h.unregister:
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
				// Update user offline status
				database.DB.Model(&models.User{}).Where("id = ?", client.UserID).Updates(map[string]interface{}{
					"is_online":    false,
					"last_seen_at": time.Now(),
				})
			}

		case message := <-h.broadcast:
			// Broadcast to all clients (or specific clients based on message)
			var wsMsg WSMessage
			if err := json.Unmarshal(message, &wsMsg); err == nil {
				// Handle different message types
				if wsMsg.Type == "message" || wsMsg.Type == "delivery_status" || wsMsg.Type == "read_receipt" {
					// Send to specific match participants
					var match models.Match
					if err := database.DB.First(&match, "id = ?", wsMsg.MatchID).Error; err == nil {
						// Send to user1
						if client, ok := h.clients[match.User1ID]; ok {
							select {
							case client.Send <- message:
							default:
								close(client.Send)
								delete(h.clients, client.UserID)
							}
						}
						// Send to user2
						if client, ok := h.clients[match.User2ID]; ok {
							select {
							case client.Send <- message:
							default:
								close(client.Send)
								delete(h.clients, client.UserID)
							}
						}
					}
				} else if wsMsg.Type == "typing" {
					// Send typing indicator to the other user in the match
					var match models.Match
					if err := database.DB.First(&match, "id = ?", wsMsg.MatchID).Error; err == nil {
						var targetUserID uuid.UUID
						senderID, _ := uuid.Parse(wsMsg.SenderID)
						if match.User1ID == senderID {
							targetUserID = match.User2ID
						} else {
							targetUserID = match.User1ID
						}
						if client, ok := h.clients[targetUserID]; ok {
							select {
							case client.Send <- message:
							default:
								close(client.Send)
								delete(h.clients, client.UserID)
							}
						}
					}
				}
			}
		}
	}
}

var hub *Hub

func init() {
	hub = NewHub()
	go hub.Run()
}

// HandleWebSocket handles WebSocket connections
func HandleWebSocket(c *websocket.Conn) {
	// Extract user from query params (token should be in query string)
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

	client := &Client{
		ID:     uuid.New(),
		UserID: userID,
		Conn:   c,
		Send:   make(chan []byte, 256),
		Hub:    hub,
	}

	client.Hub.register <- client

	// Start goroutines
	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
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

		var wsMsg WSMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			continue
		}

		// Handle incoming messages
		switch wsMsg.Type {
		case "message":
			// Save message to database and broadcast
			matchID, _ := uuid.Parse(wsMsg.MatchID)
			msg := models.Message{
				MatchID:     &matchID,
				SenderID:    c.UserID,
				MessageType: models.MessageType(wsMsg.MessageType),
				Content:     "",
				MediaURL:    wsMsg.MediaURL,
				IsRead:      false,
				IsLive:      false,
			}

			// Handle content based on message type
			if wsMsg.Content != nil {
				if contentStr, ok := wsMsg.Content.(string); ok {
					msg.Content = contentStr
				}
			}

			// Handle gift messages
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
				if err := database.DB.Create(&msg).Error; err == nil {
					// Update message ID in WS message
					wsMsg.MessageID = msg.ID.String()
					wsMsg.SenderID = c.UserID.String()
					wsMsg.ReceiverID = msg.ReceiverID.String()
					wsMsg.DeliveryStatus = "sent"
					wsMsg.Timestamp = time.Now().Format(time.RFC3339)

					// Broadcast to hub
					broadcastMsg, _ := json.Marshal(wsMsg)
					c.Hub.broadcast <- broadcastMsg

					// Send delivery status to sender
					deliveryMsg := WSMessage{
						Type:           "delivery_status",
						MatchID:        wsMsg.MatchID,
						MessageID:      msg.ID.String(),
						DeliveryStatus: "delivered",
						Timestamp:      time.Now().Format(time.RFC3339),
					}
					deliveryBytes, _ := json.Marshal(deliveryMsg)
					select {
					case c.Send <- deliveryBytes:
					default:
					}
				}
			}

		case "typing":
			// Add sender ID to typing message
			wsMsg.SenderID = c.UserID.String()
			typingMsg, _ := json.Marshal(wsMsg)
			c.Hub.broadcast <- typingMsg

		case "read_receipt":
			// Mark messages as read
			matchID, _ := uuid.Parse(wsMsg.MatchID)
			now := time.Now()
			var readMessages []models.Message
			database.DB.Model(&models.Message{}).
				Where("match_id = ? AND receiver_id = ? AND is_read = ?", matchID, c.UserID, false).
				Find(&readMessages)

			if len(readMessages) > 0 {
				messageIDs := make([]uuid.UUID, len(readMessages))
				for i, m := range readMessages {
					messageIDs[i] = m.ID
				}

				database.DB.Model(&models.Message{}).
					Where("id IN ?", messageIDs).
					Updates(map[string]interface{}{
						"is_read": true,
						"read_at": now,
					})

				// Send read receipt for each message
				for _, m := range readMessages {
					readReceipt := WSMessage{
						Type:           "read_receipt",
						MatchID:        wsMsg.MatchID,
						MessageID:      m.ID.String(),
						DeliveryStatus: "read",
						Timestamp:      time.Now().Format(time.RFC3339),
					}
					readBytes, _ := json.Marshal(readReceipt)
					c.Hub.broadcast <- readBytes
				}
			}
		}
	}
}

func (c *Client) writePump() {
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
