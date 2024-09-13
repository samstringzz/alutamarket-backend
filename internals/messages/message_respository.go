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

	// The header typically looks like "Bearer <token>"
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

func (r *repository) updateUserStatus(userID uint32, online bool) error {
	// Update the user status in the database
	if err := r.db.Model(&user.User{}).Where("id = ?", userID).Update("online", online).Error; err != nil {
		log.Printf("Error updating user status: %v", err)
		return err
	}
	return nil
}

// HandleWebSocketConnection manages WebSocket connections and messages
func (r *repository) HandleWebSocketConnection(ws *websocket.Conn, userID uint32) {
	// Register the new client connection for the user
	r.updateUserStatus(userID, true)
	mutex.Lock()
	userConnections[userID] = ws
	clients[ws] = true
	mutex.Unlock()

	defer func() {
		// Unregister the client and remove the connection
		r.updateUserStatus(userID, false)
		mutex.Lock()
		delete(clients, ws)
		delete(userConnections, userID)
		mutex.Unlock()
		ws.Close()
	}()

	for {
		var msg Message
		// Read the message from the WebSocket connection
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		// Ensure the message has the userID and chatID context
		// msg.SenderID = userID  // Associate the sender ID
		chatID := msg.ChatID // Extract the chatID from the message

		// Broadcast the message only to users in the same chat
		r.BroadcastMessages(chatID, msg)
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

// WebSocketHandler upgrades an HTTP request to a WebSocket connection
func (r *repository) WebSocketHandler(w http.ResponseWriter, req *http.Request) {
	// Extract userID from request or session
	userID, err := extractUserIDFromRequest(req) // You need to implement this
	if err != nil {
		log.Printf("Failed to upgrade to websocket: %v", err)
		return
	}
	ws, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Printf("Failed to upgrade to websocket: %v", err)
		return
	}

	// Register the new connection for the user
	mutex.Lock()
	userConnections[userID] = ws
	clients[ws] = true
	mutex.Unlock()

	r.HandleWebSocketConnection(ws, userID)
}

func NewRepository() Repository {
	dbURI := os.Getenv("DB_URI")
	// Initialize the database connection
	db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	return &repository{
		db: db,
	}
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

	// Fetch users from the chat
	var users []*user.User
	if err := r.db.Model(&chat).Association("Users").Find(&users); err != nil {
		log.Printf("Error fetching users from chat: %v", err)
		return err
	}
	var userInfoList []*user.User
	for _, u := range users {
		userInfoList = append(userInfoList, &user.User{
			ID:       u.ID,
			Fullname: u.Fullname, // Assuming Fullname is a field in the user model
			Avatar:   u.Avatar,   // Assuming Avatar is a field in the user model
		})
	}

	// senderId,_:= extractUserIDFromRequest(*http.Request)

	// Create a new message, assigning the fetched users to the message
	newMessage := Message{
		ChatID:    req.ChatID,
		Content:   req.Content,
		Media:     req.Media,
		Sender:    req.Sender,
		Users:     userInfoList, // Automatically assigned users from the chat
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

	// Fetch chats where the user is a participant
	err := r.db.
		Preload("Users"). // Eager load the Messages relationship
		Preload("LatestMessage").
		Joins("JOIN chat_users ON chat_users.chat_id = chats.id").
		Where("chat_users.user_id = ?", userID).
		Group("chats.id"). // Group by chat ID to ensure distinct chats
		Find(&chats).Error

	if err != nil {
		log.Printf("Error fetching chats for user: %v", err)
		return nil, err
	}

	// For each chat, retrieve messages
	for _, chat := range chats {
		messages, err := r.GetMessagesByChatID(ctx, chat.ID)
		if err != nil {
			log.Printf("Error fetching messages for chat %d: %v", chat.ID, err)
			return nil, err
		}
		chat.Messages = messages // Assuming Chat has a Messages field
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

	// Query the database for messages with the specified chatID
	if err := r.db.Where("chat_id = ?", chatID).Find(&messages).Error; err != nil {
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
