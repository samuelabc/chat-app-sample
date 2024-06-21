package server

import (
	"chat-app/internal/auth"
	"chat-app/pkg/models"
	"chat-app/pkg/utils"
	"database/sql"
	"encoding/json"
	"net/http"
)

type Server struct {
	DB *sql.DB
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *Server) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.Log.WithError(err).Error("Error decoding request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		utils.Log.WithError(err).Error("Error hashing password")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = auth.RegisterUser(s.DB, req.Username, hashedPassword)
	if err != nil {
		utils.Log.WithError(err).Error("Error registering user")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("User registered successfully")
}

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.Log.WithError(err).Error("Error decoding request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := auth.LoginUser(s.DB, req.Username, req.Password)
	if err != nil {
		utils.Log.WithError(err).Error("Error logging in user")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(token)
}

func (s *Server) CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	var room models.ChatRoom
	err := json.NewDecoder(r.Body).Decode(&room)
	if err != nil {
		utils.Log.WithError(err).Error("Error decoding request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := r.Context().Value("username").(string)
	row := s.DB.QueryRow("SELECT id FROM users WHERE username = ?", username)
	var userID int
	err = row.Scan(&userID)
	if err != nil {
		utils.Log.WithError(err).Error("Error fetching user ID")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = s.DB.Exec("INSERT INTO chat_rooms (name, creator_id) VALUES (?, ?)", room.Name, userID)
	if err != nil {
		utils.Log.WithError(err).Error("Error creating chat room")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Chat room created successfully")
}

func (s *Server) JoinRoomHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RoomID int `json:"room_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.Log.WithError(err).Error("Error decoding request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := r.Context().Value("username").(string)
	row := s.DB.QueryRow("SELECT id FROM users WHERE username = ?", username)
	var userID int
	err = row.Scan(&userID)
	if err != nil {
		utils.Log.WithError(err).Error("Error fetching user ID")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = s.DB.Exec("INSERT INTO room_users (room_id, user_id) VALUES (?, ?)", req.RoomID, userID)
	if err != nil {
		utils.Log.WithError(err).Error("Error joining chat room")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Joined chat room successfully")
}

func (s *Server) LeaveRoomHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RoomID int `json:"room_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.Log.WithError(err).Error("Error decoding request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := r.Context().Value("username").(string)
	row := s.DB.QueryRow("SELECT id FROM users WHERE username = ?", username)
	var userID int
	err = row.Scan(&userID)
	if err != nil {
		utils.Log.WithError(err).Error("Error fetching user ID")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = s.DB.Exec("DELETE FROM room_users WHERE room_id = ? AND user_id = ?", req.RoomID, userID)
	if err != nil {
		utils.Log.WithError(err).Error("Error leaving chat room")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Left chat room successfully")
}

func (s *Server) ListUsersInRoomHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RoomID int `json:"room_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.Log.WithError(err).Error("Error decoding request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rows, err := s.DB.Query("SELECT users.username FROM users JOIN room_users ON users.id = room_users.user_id WHERE room_users.room_id = ?", req.RoomID)
	if err != nil {
		utils.Log.WithError(err).Error("Error fetching users in room")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users := []string{}
	for rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			utils.Log.WithError(err).Error("Error scanning user row")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, username)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func (s *Server) ListRoomsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := s.DB.Query("SELECT id, name FROM chat_rooms")
	if err != nil {
		utils.Log.WithError(err).Error("Error fetching chat rooms")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	rooms := []models.ChatRoom{}
	for rows.Next() {
		var room models.ChatRoom
		err := rows.Scan(&room.ID, &room.Name)
		if err != nil {
			utils.Log.WithError(err).Error("Error scanning room row")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rooms = append(rooms, room)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rooms)
}
