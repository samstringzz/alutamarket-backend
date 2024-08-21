package messages

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

var (
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan Message)
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// Allow requests from certain origins
			return r.Host == os.Getenv("BASE_URL") || r.Host == "http://localhost:5173"
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	mutex sync.Mutex
)

// HandleWebSocketConnection manages WebSocket connections and messages
func HandleWebSocketConnection(ws *websocket.Conn) {
	// Register the new client
	mutex.Lock()
	clients[ws] = true
	mutex.Unlock()

	defer func() {
		// Unregister the client
		mutex.Lock()
		delete(clients, ws)
		mutex.Unlock()
		ws.Close()
	}()

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
		broadcast <- msg
	}
}

// BroadcastMessages listens for new messages and sends them to all clients
func BroadcastMessages() {
	for {
		msg := <-broadcast
		mutex.Lock()
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Error writing message to client: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
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

func (r *repository) SendMessage(ctx context.Context, req *Message) error {
	// var clients = make(map[*websocket.Conn]bool)
	// var broadcast = make(chan Message)

	newMessage := Message{
		ChatID:  req.ChatID,
		Content: req.Content,
		Media:   req.Media,
		User:    req.User,
	}
	if err := r.db.Create(newMessage).Error; err != nil {

		log.Printf("Error creating message: %v", err)
		return err
	}
	broadcast <- newMessage

	return nil
}

func (r *repository) GetChatLists(ctx context.Context, userID int64) ([]*Chat, error) {
	var chats []*Chat
	// Use the @> operator to check if userID is present in the users array.
	if err := r.db.Where("users @> ?", pq.Array([]int64{userID})).Find(&chats).Error; err != nil {
		return nil, err
	}

	return chats, nil
}

func (r *repository) CreateChat(ctx context.Context, chatID uint32, usersID pq.Int64Array) error {
	var chat Chat
	var existingChat Chat

	// Check if the chat already exists
	if err := r.db.Where("id = ?", chatID).First(&existingChat).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Chat does not exist, create a new chat
			chat = Chat{
				ID:    chatID,
				Users: usersID,
			}
			// fmt.Println("Hello1")

			if err := r.db.Create(&chat).Error; err != nil {
				// fmt.Println("Hello2")

				return err
			}
		} else {
			// fmt.Println("Hello3")

			return err
		}
	} else {
		// Chat exists, use the existingChat variable
		chat = existingChat
	}

	for _, usrID := range usersID {
		if !containsUser(chat.Users, usrID) {
			chat.Users = append(chat.Users, usrID)
		}
	}
	// fmt.Println("Hello")
	// Save the updated chat with users
	if err := r.db.Save(&chat).Error; err != nil {
		return err
	}

	return nil
}

// Helper function to check if a user is in the list
func containsUser(users pq.Int64Array, userID int64) bool {
	for _, id := range users {
		if id == userID {
			return true
		}
	}
	return false
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
