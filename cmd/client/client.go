package main

import (
	"bufio"
	"bytes"
	"chat-app/pkg/models"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

const tokenFileBaseName = "token"

var loggedInUsername = ""

func saveToken(username string, token string) error {
	loggedInUsername = username
	tokenFile := fmt.Sprintf("%s_%s.txt", tokenFileBaseName, username)
	return os.WriteFile(tokenFile, []byte(token), 0644)
}

func getToken() (string, error) {
	tokenFile := fmt.Sprintf("%s_%s.txt", tokenFileBaseName, loggedInUsername)
	token, err := os.ReadFile(tokenFile)
	if err != nil {
		return "", err
	}
	return string(token), nil
}

func main() {
	fmt.Println("Chat-app started. Type 'exit' to quit.")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		args := strings.Split(input, " ")

		if len(args) == 0 {
			continue
		}

		command := args[0]

		if command == "exit" {
			fmt.Println("Exiting chat-app.")
			break
		}

		switch command {

		case "register":
			if len(args) != 3 {
				fmt.Println("Usage: register <username> <password>")

			}
			username := args[1]
			password := args[2]

			user := map[string]string{
				"username": username,
				"password": password,
			}

			jsonUser, err := json.Marshal(user)
			if err != nil {
				fmt.Println("Error marshalling user:", err)
				continue
			}

			resp, err := http.Post("http://localhost:8080/register", "application/json", bytes.NewBuffer(jsonUser))
			if err != nil {
				fmt.Println("Error making request:", err)
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusCreated {
				fmt.Println("Error registering user:", resp.Status)
				continue
			}

			fmt.Println("User registered successfully")

		case "login":
			if len(args) != 3 {
				fmt.Println("Usage: login <username> <password>")
			}
			username := args[1]
			password := args[2]

			user := map[string]string{
				"username": username,
				"password": password,
			}

			jsonUser, err := json.Marshal(user)
			if err != nil {
				fmt.Println("Error marshalling user:", err)
				continue
			}

			resp, err := http.Post("http://localhost:8080/login", "application/json", bytes.NewBuffer(jsonUser))
			if err != nil {
				fmt.Println("Error making request:", err)
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				fmt.Println("Error logging in:", resp.Status)
				continue
			}

			var token map[string]string
			err = json.NewDecoder(resp.Body).Decode(&token)
			if err != nil {
				fmt.Println("Error decoding token:", err)
				continue
			}

			err = saveToken(username, token["token"])
			if err != nil {
				fmt.Println("Error saving token:", err)
				continue
			}

			fmt.Println("User logged in successfully")

		case "logout":
			err := os.Remove(fmt.Sprintf("%s_%s.txt", tokenFileBaseName, loggedInUsername))
			if err != nil {
				fmt.Println("Error removing token:", err)
				continue
			}
			loggedInUsername = ""
			fmt.Println("User logged out successfully")

		case "create-room":
			if len(args) != 2 {
				fmt.Println("Usage: create-room <room_name>")

			}
			token, err := getToken()
			if err != nil {
				fmt.Println("Error reading token:", err)
				continue
			}
			roomName := args[1]

			room := map[string]string{
				"name": roomName,
			}

			jsonRoom, err := json.Marshal(room)
			if err != nil {
				fmt.Println("Error marshalling room:", err)
				continue
			}

			client := &http.Client{}
			req, err := http.NewRequest("POST", "http://localhost:8080/create-room", bytes.NewBuffer(jsonRoom))
			if err != nil {
				fmt.Println("Error creating request:", err)
				continue
			}
			req.Header.Add("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error making request:", err)
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusCreated {
				fmt.Println("Error creating room:", resp.Status)

			}

			fmt.Println("Chat room created successfully")

		case "join-room":
			if len(args) != 2 {
				fmt.Println("Usage: join-room <room_id>")

			}
			token, err := getToken()
			if err != nil {
				fmt.Println("Error reading token:", err)
				continue
			}
			roomID, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("Invalid room ID:", err)
				continue
			}

			room := map[string]int{
				"room_id": roomID,
			}

			jsonRoom, err := json.Marshal(room)
			if err != nil {
				fmt.Println("Error marshalling room:", err)
				continue
			}

			client := &http.Client{}
			req, err := http.NewRequest("POST", "http://localhost:8080/join-room", bytes.NewBuffer(jsonRoom))
			if err != nil {
				fmt.Println("Error creating request:", err)
				continue
			}
			req.Header.Add("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error making request:", err)
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				fmt.Println("Error joining room:", resp.Status)
				continue
			}

			fmt.Println("Joined chat room successfully")

		case "leave-room":
			if len(args) != 2 {
				fmt.Println("Usage: leave-room <room_id>")
				continue
			}
			token, err := getToken()
			if err != nil {
				fmt.Println("Error reading token:", err)
				continue
			}
			roomID, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("Invalid room ID:", err)
				continue
			}

			room := map[string]int{
				"room_id": roomID,
			}

			jsonRoom, err := json.Marshal(room)
			if err != nil {
				fmt.Println("Error marshalling room:", err)
				continue
			}

			client := &http.Client{}
			req, err := http.NewRequest("POST", "http://localhost:8080/leave-room", bytes.NewBuffer(jsonRoom))
			if err != nil {
				fmt.Println("Error creating request:", err)
				continue
			}
			req.Header.Add("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error making request:", err)
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				fmt.Println("Error leaving room:", resp.Status)
				continue
			}

			fmt.Println("Left chat room successfully")

		case "list-users":
			if len(args) != 2 {
				fmt.Println("Usage: list-users <room_id>")

			}
			token, err := getToken()
			if err != nil {
				fmt.Println("Error reading token:", err)
				continue
			}
			roomID, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("Invalid room ID:", err)
				continue
			}

			room := map[string]int{
				"room_id": roomID,
			}

			jsonRoom, err := json.Marshal(room)
			if err != nil {
				fmt.Println("Error marshalling room:", err)
				continue
			}

			client := &http.Client{}
			req, err := http.NewRequest("POST", "http://localhost:8080/list-users", bytes.NewBuffer(jsonRoom))
			if err != nil {
				fmt.Println("Error creating request:", err)

			}
			req.Header.Add("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error making request:", err)
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				fmt.Println("Error listing users:", resp.Status)

			}

			var users []string
			err = json.NewDecoder(resp.Body).Decode(&users)
			if err != nil {
				fmt.Println("Error decoding users:", err)
				continue
			}

			fmt.Println("Users in chat room:")
			for _, user := range users {
				fmt.Printf("- %s\n", user)
			}

		case "list-rooms":
			token, err := getToken()
			if err != nil {
				fmt.Println("Error reading token:", err)
				continue
			}

			client := &http.Client{}
			req, err := http.NewRequest("GET", "http://localhost:8080/list-rooms", nil)
			if err != nil {
				fmt.Println("Error creating request:", err)
				continue
			}
			req.Header.Add("Authorization", "Bearer "+token)

			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error making request:", err)
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				fmt.Println("Error listing rooms:", resp.Status)

			}

			var rooms []models.ChatRoom
			err = json.NewDecoder(resp.Body).Decode(&rooms)
			if err != nil {
				fmt.Println("Error decoding rooms:", err)
				continue
			}

			fmt.Println("Available chat rooms:")
			for _, room := range rooms {
				fmt.Printf("- %s (ID: %d)\n", room.Name, room.ID)
			}

		case "enter-room":
			if len(args) != 2 {
				fmt.Println("Usage: enter-room <room_id>")
				continue
			}
			token, err := getToken()
			if err != nil {
				fmt.Println("Error reading token:", err)
				continue
			}

			roomID := args[1]

			c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", http.Header{"Authorization": {"Bearer " + token}})
			if err != nil {
				fmt.Println("Error connecting to WebSocket:", err)
				os.Exit(1)
			}
			defer c.Close()

			go func() {
				for {
					_, message, err := c.ReadMessage()
					if err != nil {
						fmt.Println("Websocket channel closed.", err)
						return
					}

					type Message struct {
						SenderID    int    `json:"sender_id"`
						RecipientID int    `json:"recipient_id,omitempty"`
						RoomID      int    `json:"room_id,omitempty"`
						Content     string `json:"content"`
					}
					var msg Message
					if err := json.Unmarshal(message, &msg); err != nil {
						fmt.Println("Error unmarshalling message:", err)
						continue
					}

					if msg.RoomID != 0 {
						fmt.Printf("[Room %d] (User %d): %s\n", msg.RoomID, msg.SenderID, msg.Content)
					} else if msg.RecipientID != 0 {
						fmt.Printf("[DM from User %d]: %s\n", msg.SenderID, msg.Content)
					}
				}
			}()

			for {
				fmt.Println("Enter message (or !dm-<userid> <message> for direct message, or !leave to leave the room):")
				reader := bufio.NewReader(os.Stdin)

				content, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("Error reading input:", err)
					return
				}
				content = strings.TrimSpace(content) // Remove the newline character at the end

				if content == "!leave" {
					fmt.Println("Leaving the room...")
					c.Close()
					break
				}

				if content == "" {
					continue
				}

				type Message struct {
					RecipientID int    `json:"recipient_id,omitempty"`
					RoomID      int    `json:"room_id,omitempty"`
					Content     string `json:"content"`
				}

				var msg Message
				fmt.Println("message", content)
				if strings.HasPrefix(content, "!dm-") {
					parts := strings.SplitN(content, " ", 2)
					if len(parts) < 2 {
						fmt.Println("Invalid DM format. Use !dm-<userid> <message>")
						continue
					}
					userIDStr := strings.TrimPrefix(parts[0], "!dm-")
					userID, err := strconv.Atoi(userIDStr)
					if err != nil {
						fmt.Println("Invalid user ID:", err)
						continue
					}
					msg.RecipientID = userID
					msg.Content = parts[1]
				} else {
					roomIDInt, err := strconv.Atoi(roomID)
					if err != nil {
						fmt.Println("Invalid room ID:", err)
						continue
					}
					msg.RoomID = roomIDInt
					msg.Content = content
				}

				messageBytes, err := json.Marshal(msg)
				if err != nil {
					fmt.Println("Error marshalling message:", err)
					continue
				}

				err = c.WriteMessage(websocket.TextMessage, messageBytes)
				if err != nil {
					fmt.Println("Error sending message:", err)
					return
				}
			}

		case "":
			continue

		default:
			fmt.Println("Unknown command:", command)
		}
	}
}
