package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const schema = `
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS chat_rooms (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		creator_id INTEGER NOT NULL,
		FOREIGN KEY (creator_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS room_users (
		room_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		FOREIGN KEY (room_id) REFERENCES chat_rooms(id),
		FOREIGN KEY (user_id) REFERENCES users(id),
		PRIMARY KEY (room_id, user_id)
);

CREATE TABLE IF NOT EXISTS messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sender_id INTEGER,
    recipient_id INTEGER,
    room_id INTEGER,
    content TEXT NOT NULL,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sender_id) REFERENCES users(id),
    FOREIGN KEY (recipient_id) REFERENCES users(id),
    FOREIGN KEY (room_id) REFERENCES chat_rooms(id)
);

`

func InitDatabase(dataSourceName string) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		fmt.Println("Error opening database:", err)
		os.Exit(1)
	}
	defer db.Close()

	_, err = db.Exec(schema)
	if err != nil {
		fmt.Println("Error initializing database:", err)
		os.Exit(1)
	}
}

func ClearDatabase(dataSourceName string) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		fmt.Println("Error opening database:", err)
		os.Exit(1)
	}
	defer db.Close()

	_, err = db.Exec("DROP TABLE IF EXISTS users")
	if err != nil {
		fmt.Println("Error clearing database:", err)
		os.Exit(1)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS chat_rooms")
	if err != nil {
		fmt.Println("Error clearing database:", err)
		os.Exit(1)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS room_users")
	if err != nil {
		fmt.Println("Error clearing database:", err)
		os.Exit(1)
	}
}
