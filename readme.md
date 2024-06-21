# Chat Application Documentation

## Features

- **User Registration and Login**
- **JWT-based Authentication**
- **WebSocket-based Real-time Communication**
- **Chat Room Management** (Create, Join, Leave)
- **In-room Direct Messaging**
- **Error Handling and Logging**
- **Deployment using Docker and Docker Compose**

## Technologies Used

- **Go** for the server-side application
  - Entry point: `main.go`
- **Go** for CLI application
  - Entry point: `cmd/client/client.go`
- **SQLite** as the database for business logic relevant data
  - Database file location: `chat-app.db`
  - `users` table: to store user information
    - Columns: `id`, `username`, `password_hash`
  - `chat_rooms` table: to store chat room information
		- Columns: `id`, `name`, `creator_id`
  - `room_users` table: to store user-room mapping
		- Columns: `room_id`, `user_id`
  - `messages` table: to store chat messages(both group and direct messages)
    - Columns: `id`, `sender_id`, `recipient_id`, `room_id`, `content`, `timestamp`
- **Gorilla WebSocket** for WebSocket implementation
  - Relevant code: `internal/handlers/websocket.go`
  - Whenever a new WebSocket connection is established, a new `Client` object is created to handle the connection
  - The `Client` object listens for incoming messages and broadcasts them to all users in the same room, or sends direct messages to specific users
  - A list of all connected clients is maintained in the `clients` map
  - To avoid race conditions, a `mutex` is used to synchronize access to the `clients` map
- **JWT** for user authentication
  - Relevant code: `internal/auth/*`
- **Logrus** for logging
  - Relevant code: `pkg/utils/logger.go`
  - Log file location: `log/chat-app.log`
- **Docker**, **Docker Compose** for deployment

## Usage Instructions

### Start the Server

To start the server using Docker Compose, run:

```sh
docker-compose up
```

To start the by running the server directly, run:

```sh
go run main.go
```

### Start the CLI Client

To start the CLI client, run:

```sh
go run cmd/client/client.go
```

### CLI Commands

#### User Registration and Authentication

##### Register a new user

To register a new user:

```sh
register <username> <password>
```

##### Login user

To log in as an existing user:

```sh
login <username> <password>
```

##### Logout User

To log out the current user:

```sh
logout
```

#### Chat Room Management

##### List Rooms

To list all available chat rooms:

```
list-rooms
```

#### Create Room

To create a new chat room:

```sh
create-room <room_name>
```

#### Join Room

To join an existing chat room:

```sh
join-room <room_id>
```

#### Leave Room

To leave a chat room:

```sh
leave-room <room_id>
```

#### List Users

To list all users in a chat room:

```sh
list-users <room_id>
```

### Messaging

#### Enter a Room

To enter a chat room and start participating in group conversations or direct messaging:

```sh
enter-room <room_in>
```

#### Broadcast Message in a Room

Once you have entered a room, type your message and press Enter. The message will be broadcasted to everyone in the room:

```sh
<message>
```

#### Send Direct Message in a Room

To send a direct message to a specific user in the room, use the following format:

```sh
!dm-<user_id> <messages>
```

#### Exit a Room

To exit from the room and return to the main CLI interface:

```sh
!leave
```
