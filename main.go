package main

import (
	"chat-app/internal/auth"
	"chat-app/internal/database"
	"chat-app/internal/server"
	"chat-app/internal/websocket"
	"chat-app/pkg/utils"
	"database/sql"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbPath := os.Getenv("DATABASE_URL")
	if dbPath == "" {
		dbPath = "./chat-app.db"
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		utils.Log.WithError(err).Fatal("Failed to connect to database")
	}
	defer func() {
		if err := db.Close(); err != nil {
			utils.Log.WithError(err).Error("Error closing database")
		}
	}()

	// Initialize the database
	database.InitDatabase(dbPath)

	srv := &server.Server{DB: db}

	http.Handle("/register", http.HandlerFunc(srv.RegisterHandler))
	http.Handle("/login", http.HandlerFunc(srv.LoginHandler))
	http.Handle("/create-room", auth.JWTMiddleware(http.HandlerFunc(srv.CreateRoomHandler)))
	http.Handle("/join-room", auth.JWTMiddleware(http.HandlerFunc(srv.JoinRoomHandler)))
	http.Handle("/leave-room", auth.JWTMiddleware(http.HandlerFunc(srv.LeaveRoomHandler)))
	http.Handle("/list-users", auth.JWTMiddleware(http.HandlerFunc(srv.ListUsersInRoomHandler)))
	http.Handle("/list-rooms", auth.JWTMiddleware(http.HandlerFunc(srv.ListRoomsHandler)))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		userID, _, err := auth.ValidateJWT(r.Header.Get("Authorization")[7:])
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		websocket.HandleConnections(w, r, db, userID)
	})

	websocket.Init()

	utils.Log.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		utils.Log.WithError(err).Fatal("Server failed")
	}
}
