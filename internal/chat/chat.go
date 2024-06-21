package chat

import (
	"chat-app/pkg/models"
	"database/sql"
)

// CreateChatRoom creates a new chat room with the given name
func CreateChatRoom(db *sql.DB, name string, creatorID int) error {
	_, err := db.Exec("INSERT INTO chat_rooms (name, creator_id) VALUES (?, ?)", name, creatorID)
	if err != nil {
		return err
	}
	return nil
}

// JoinChatRoom adds a user to a chat room
func JoinChatRoom(db *sql.DB, roomID, userID int) error {
	_, err := db.Exec("INSERT INTO room_users (room_id, user_id) VALUES (?, ?)", roomID, userID)
	if err != nil {
		return err
	}
	return nil
}

// LeaveChatRoom removes a user from a chat room
func LeaveChatRoom(db *sql.DB, roomID, userID int) error {
	_, err := db.Exec("DELETE FROM room_users WHERE room_id = ? AND user_id = ?", roomID, userID)
	if err != nil {
		return err
	}
	return nil
}

// ListChatRooms lists all available chat rooms
func ListChatRooms(db *sql.DB) ([]models.ChatRoom, error) {
	rows, err := db.Query("SELECT id, name FROM chat_rooms")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chatRooms := []models.ChatRoom{}
	for rows.Next() {
		var chatRoom models.ChatRoom
		err := rows.Scan(&chatRoom.ID, &chatRoom.Name)
		if err != nil {
			return nil, err
		}
		chatRooms = append(chatRooms, chatRoom)
	}
	return chatRooms, nil
}

// ListUsersInChatRoom lists all users in a chat room
func ListUsersInChatRoom(db *sql.DB, roomID int) ([]models.User, error) {
	rows, err := db.Query("SELECT users.id, users.username FROM users JOIN room_users ON users.id = room_users.user_id WHERE room_users.room_id = ?", roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
