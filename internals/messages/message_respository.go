package messages

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Chrisentech/aluta-market-api/internals/user"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/websocket"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// var jwtSecret = []byte(os.Getenv("SECRET_KEY"))

func extractUserIDFromRequest(r *http.Request) (uint32, error) {
	// Get the Authorization header and query parameter
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		// If Authorization header is missing, check the query parameter
		authHeader = "Bearer " + r.URL.Query().Get("token")
	}
	if authHeader == "" {
		return 0, fmt.Errorf("authorization header and token query parameter missing")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate that the token is signed using HMAC-SHA256 (HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SECRET_KEY")), nil
		}
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	})
	if err != nil {
		return 0, fmt.Errorf("error parsing token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extract user ID from claims (assuming it's a string in the "id" claim)
		userIDStr, ok := claims["id"].(string)
		if !ok {
			return 0, fmt.Errorf("user id not found in token claims or not a string")
		}

		// Convert the string userID to uint32
		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			return 0, fmt.Errorf("error converting user id to uint32: %v", err)
		}

		return uint32(userID), nil
	}

	return 0, fmt.Errorf("invalid token")
}

var (
	clients         = make(map[*websocket.Conn]bool)
	userConnections = make(map[uint32]*websocket.Conn)
	upgrader        = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			allowedOrigins := []string{
				os.Getenv("BASE_URL"),
				"http://localhost:5173",
			}

			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					return true
				}
			}
			return false
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	mutex sync.Mutex
)

// Update the updateUserStatus function to be more efficient
func (r *repository) updateUserStatus(userID uint32, online bool) error {
	// Add index hint and optimize the update
	result := r.db.Model(&user.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"online":     online,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		log.Printf("Error updating user status: %v", result.Error)
		return result.Error
	}
	return nil
}

// Modify WebSocketHandler to handle disconnections more gracefully
func (r *repository) WebSocketHandler(w http.ResponseWriter, req *http.Request) {
	userID, err := extractUserIDFromRequest(req)
	if err != nil {
		log.Printf("Failed to extract user ID: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Configure WebSocket
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true // You might want to implement proper origin checking
	}

	ws, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Printf("Failed to upgrade to websocket: %v", err)
		return
	}

	// Set read deadline to detect stale connections
	ws.SetReadDeadline(time.Now().Add(60 * time.Second))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Start ping-pong routine
	go func() {
		ticker := time.NewTicker(54 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
					return
				}
			}
		}
	}()

	r.HandleWebSocketConnection(ws, userID)
}

// Keep this original declaration and update it with the improved error handling
func (r *repository) HandleWebSocketConnection(ws *websocket.Conn, userID uint32) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Recovered from panic in HandleWebSocketConnection: %v", err)
		}
		r.updateUserStatus(userID, false)
		mutex.Lock()
		delete(clients, ws)
		delete(userConnections, userID)
		mutex.Unlock()
		ws.Close()
	}()

	r.updateUserStatus(userID, true)
	mutex.Lock()
	userConnections[userID] = ws
	clients[ws] = true
	mutex.Unlock()

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("Unexpected websocket error: %v", err)
			}
			break
		}

		if msg.ChatID > 0 {
			r.BroadcastMessages(msg.ChatID, msg)
		}
	}
}

// BroadcastMessages listens for new messages and sends them to all clients
func (r *repository) BroadcastMessages(chatID uint32, msg Message) {
	// Get the participants of the chat
	chat, err := r.FindChatByID(chatID)
	if err != nil {
		log.Printf("Error finding chat: %v", err)

	}

	// Fetch users from the chat
	var chatUsers []*user.User
	var chatUserIDs []uint32
	if err := r.db.Model(&chat).Association("Users").Find(&chatUsers); err != nil {
		log.Printf("Error fetching users from chat: %v", err)

	}

	// Extract user IDs from the chatUsers
	for _, user := range chatUsers {
		chatUserIDs = append(chatUserIDs, user.ID) // Assuming `user.ID` is uint32
	}

	mutex.Lock()
	for _, userID := range chatUserIDs {
		// Check if the user is connected
		if conn, ok := userConnections[userID]; ok {
			err := conn.WriteJSON(msg)
			if err != nil {
				log.Printf("Error sending message to user %d: %v", userID, err)
				conn.Close()
				delete(userConnections, userID)
			}
		}
	}
	mutex.Unlock()
}

func NewRepository() Repository {
	dbURI := os.Getenv("DB_URI")
	if dbURI == "" {
		log.Fatal("DB_URI environment variable is not set")
	}

	db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Verify connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	// Test the connection
	err = sqlDB.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	return &repository{
		db: db,
	}
}

func (r *repository) GetChatUsers(ctx context.Context, chatID uint32) ([]*user.User, error) {
	var chat Chat
	if err := r.db.First(&chat, chatID).Error; err != nil {
		return nil, fmt.Errorf("chat not found: %v", err)
	}

	var users []*user.User
	if err := r.db.Model(&chat).Association("Users").Find(&users); err != nil {
		return nil, fmt.Errorf("failed to get chat users: %v", err)
	}

	return users, nil
}

func (r *repository) FindChatByID(chatID uint32) (*Chat, error) {
	var chat Chat
	if err := r.db.First(&chat, chatID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("chat not found")
		}
		return nil, fmt.Errorf("error finding chat: %v", err)
	}
	return &chat, nil
}

func (r *repository) SendMessage(ctx context.Context, req *Message) error {
	// Find the chat by its ID
	chat, err := r.FindChatByID(req.ChatID)
	if err != nil {
		log.Printf("Error finding chat: %v", err)
		return err
	}

	// Create a new message
	newMessage := Message{
		ChatID:    req.ChatID,
		Content:   req.Content,
		Media:     req.Media,
		Sender:    req.Sender,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save the new message to the database
	if err := r.db.Create(&newMessage).Error; err != nil {
		log.Printf("Error creating message: %v", err)
		return err
	}

	// Broadcast message only to users in the chat
	r.BroadcastMessages(newMessage.ChatID, newMessage)

	// Update the chat with the latest message ID and increment the unread count
	if err := r.db.Model(&chat).Updates(map[string]interface{}{
		"LatestMessageID": newMessage.ID,
		"UnreadCount":     chat.UnreadCount + 1,
	}).Error; err != nil {
		log.Printf("Error updating chat with latest message ID and unread count: %v", err)
		return err
	}
	return nil
}

func (r *repository) GetChatLists(ctx context.Context, userID uint32) ([]*Chat, error) {
	var chats []*Chat

	err := r.db.
		Preload("Users").
		Preload("Messages").
		Preload("Messages.User"). // This now correctly references the User relationship
		Preload("LatestMessage").
		Joins("JOIN chat_users ON chat_users.chat_id = chats.id").
		Where("chat_users.user_id = ?", userID).
		Find(&chats).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get chat lists: %v", err)
	}

	return chats, nil
}
func (r *repository) FindOrCreateChat(ctx context.Context, users []*user.User) (*Chat, error) {
	var chat Chat

	// Query to check if a chat exists with the exact users
	// You need to use a join table for the chat_users relationship and group by chat ID
	err := r.db.
		Table("chats").
		Select("chats.*").
		Joins("JOIN chat_users ON chat_users.chat_id = chats.id").
		Where("chat_users.user_id IN (?)", extractUserIDs(users)). // Extract user IDs from the passed users
		Group("chats.id").
		Having("COUNT(DISTINCT chat_users.user_id) = ?", len(users)). // Ensure the number of users matches
		First(&chat).Error

	// If a matching chat is found, return it
	if err == nil {
		fmt.Printf("Existing chat found: %+v\n", chat) // Print the chat details

		return &chat, nil
	}

	// If no matching chat is found, create a new one
	if errors.Is(err, gorm.ErrRecordNotFound) {
		newChat := Chat{
			Users:       users,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			UnreadCount: 0,
		}

		// Save the new chat to the database
		if err := r.db.Create(&newChat).Error; err != nil {
			log.Printf("Error creating new chat: %v", err)
			return nil, err
		}
		// Fetch users from the chat
		var users []*user.User
		if err := r.db.Model(&newChat).Association("Users").Find(&users); err != nil {
			log.Printf("Error fetching users from chat: %v", err)
			return nil, err
		}
		var userInfoList []*user.User
		for _, u := range users {
			userInfoList = append(userInfoList, &user.User{
				ID:       u.ID,
				Fullname: u.Fullname, // Assuming Fullname is a field in the user model
				Avatar:   u.Avatar,   // Assuming Avatar is a field in the user model
			})
		}
		// In FindOrCreateChat method, replace the message creation part
		newMessage := Message{
			ChatID:    newChat.ID,
			Content:   "Hello, I was surfing through your product and would like to make some enquiry/complains",
			Sender:    users[0].ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		// Save the new message to the database
		if err := r.db.Create(&newMessage).Error; err != nil {
			log.Printf("Error creating message: %v", err)
			return nil, err
		}
		// Update the chat with the latest message ID and increment the unread count
		if err := r.db.Model(&newChat).Updates(map[string]interface{}{
			"LatestMessageID": newMessage.ID,
			"UnreadCount":     newChat.UnreadCount + 1,
		}).Error; err != nil {
			log.Printf("Error updating chat with latest message ID and unread count: %v", err)
			return nil, err
		}
		return &newChat, nil
	}

	// Return any other error encountered
	return nil, err
}

// Helper function to extract user IDs from the list of users
func extractUserIDs(users []*user.User) []uint32 {
	var ids []uint32
	for _, user := range users {
		ids = append(ids, user.ID)
	}
	return ids
}
func (r *repository) GetMessages(ctx context.Context, chatID string) ([]*Message, error) {
	var messages []*Message

	// Convert chatID string to uint32
	chatIDUint, err := strconv.ParseUint(chatID, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid chat ID: %v", err)
	}

	// Query the database for messages with the specified chatID
	if err := r.db.
		Where("chat_id = ?", uint32(chatIDUint)).
		Order("created_at ASC").
		Find(&messages).Error; err != nil {
		log.Printf("Error retrieving messages for chatID %s: %v", chatID, err)
		return nil, err
	}

	return messages, nil
}

func (r *repository) GetMessagesByChatID(ctx context.Context, chatID uint32) ([]*Message, error) {
	var messages []*Message

	// Fetch messages for a specific chatID and order by created_at
	err := r.db.
		Where("chat_id = ?", chatID).
		Order("created_at ASC"). // Ensure messages are ordered by creation time
		Find(&messages).Error

	if err != nil {
		log.Printf("Error retrieving messages for chatID %d: %v", chatID, err)
		return nil, err
	}

	return messages, nil
}
