package websocket

import (
	"chat-app/pkg/utils"
	"database/sql"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	ID     int
	Conn   *websocket.Conn
	Send   chan []byte
	DB     *sql.DB
	UserID int
}

type Message struct {
	SenderID    int    `json:"sender_id"`
	RecipientID int    `json:"recipient_id,omitempty"`
	RoomID      int    `json:"room_id,omitempty"`
	Content     string `json:"content"`
}

var (
	clients   = make(map[int]*Client)
	broadcast = make(chan Message)
	mutex     sync.Mutex
)

func HandleConnections(w http.ResponseWriter, r *http.Request, db *sql.DB, userID int) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.Log.WithError(err).Error("Error upgrading to WebSocket")
		return
	}

	mutex.Lock()

	// Check if a client with the same userID already exists
	if _, exists := clients[userID]; exists {
		utils.Log.WithField("userID", userID).Error("Client with this userID already exists")
		conn.Close() // Close the new connection as it's a duplicate
		return
	}

	client := &Client{
		Conn:   conn,
		Send:   make(chan []byte),
		DB:     db,
		UserID: userID,
	}

	clients[userID] = client

	mutex.Unlock()

	go client.readMessages()
	go client.writeMessages()
}

func (c *Client) readMessages() {
	defer func() {
		mutex.Lock()
		delete(clients, c.UserID)
		mutex.Unlock()
		c.Conn.Close()
	}()
	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			utils.Log.WithError(err).Error("Error reading message")
			return
		}
		isSenderInRoom := senderInRoom(c.DB, c.UserID, 1)
		if !isSenderInRoom {
			utils.Log.WithField("userID", c.UserID).Error("Sender not in room")
			continue
		}

		var message Message
		if err := json.Unmarshal(msg, &message); err != nil {
			utils.Log.WithError(err).Error("Error unmarshalling message")
			continue
		}
		message.SenderID = c.UserID
		broadcast <- message
	}
}

func (c *Client) writeMessages() {
	defer c.Conn.Close()
	for msg := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			utils.Log.WithError(err).Error("Error writing message")
			return
		}
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		jsonMsg, _ := json.Marshal(msg)

		if msg.RoomID != 0 {
			mutex.Lock()
			for userID, client := range clients {
				if isUserInRoom(client.DB, userID, msg.RoomID) {
					client.Send <- jsonMsg
				}
			}

			mutex.Unlock()
		} else if msg.RecipientID != 0 {
			mutex.Lock()

			for userId, client := range clients {
				if userId == msg.RecipientID || userId == msg.SenderID {
					client.Send <- jsonMsg
				}
			}
			mutex.Unlock()
		}

		saveMessageToDB(msg)
	}
}

func isUserInRoom(db *sql.DB, userID, roomID int) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM room_users WHERE user_id = ? AND room_id = ?", userID, roomID).Scan(&count)
	if err != nil || count == 0 {
		return false
	}
	return true
}

func saveMessageToDB(msg Message) {
	db := getDBConnection()
	defer db.Close()

	if msg.RoomID != 0 {
		_, err := db.Exec("INSERT INTO messages (sender_id, room_id, content) VALUES (?, ?, ?)", msg.SenderID, msg.RoomID, msg.Content)
		if err != nil {
			utils.Log.WithError(err).Error("Error saving room message to database")
		}
	} else if msg.RecipientID != 0 {
		_, err := db.Exec("INSERT INTO messages (sender_id, recipient_id, content) VALUES (?, ?, ?)", msg.SenderID, msg.RecipientID, msg.Content)
		if err != nil {
			utils.Log.WithError(err).Error("Error saving direct message to database")
		}
	}
}

func senderInRoom(db *sql.DB, senderID, roomID int) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM room_users WHERE user_id = ? AND room_id = ? limit 1", senderID, roomID).Scan(&count)
	if err != nil || count == 0 {
		return false
	}
	return true
}

func getDBConnection() *sql.DB {
	db, err := sql.Open("sqlite3", "./chat-app.db")
	if err != nil {
		utils.Log.WithError(err).Fatal("Failed to connect to database")
	}
	return db
}

func Init() {
	go handleMessages()
}
