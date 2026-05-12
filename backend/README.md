# Voting-Go

Project built with [Gini](https://gini-webserver.up.railway.app/)

This project provides a platform where a host can create a poll and voters can see the results update instantly as votes come in. It is built using Go and utilizes WebSockets for real-time data broadcasting.

---

## Features

- Real-time result updates via WebSockets.
- Concurrent-safe voting logic using sync primitives.
- RESTful API for poll creation and voting.
- Lightweight in-memory data management.

## Technical Stack

- Language: Go (v1.26+)
- WebSocket Library: Gorilla WebSocket
- Routing: Standard Library http.ServeMux

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/rafaeldepontes/voting-go.git
   cd voting-go
   ```

2. Download dependencies:
   ```bash
   go mod tidy
   ```

## Usage

### Starting the Server

```bash
go run cmd/server/main.go
```
The server defaults to port 8080.

### List all polls

```bash
curl http://localhost:8080/polls
```

### Creating a Poll

```bash
curl -X POST http://localhost:8080/polls \
     -d '{"name": "The Best Programming Language", "options": ["Golang", "Rust", "Kotlin", "C#"]}'
```

### Voting

```bash
curl -X POST http://localhost:8080/polls/{id}/vote \
     -d '{"optionId": 1}'
```

### Real-Time Monitoring

Connect to the WebSocket endpoint for a specific poll:
`ws://localhost:8080/ws/polls/{id}`

## Testing

Run unit tests for the core voting service:
```bash
go test ./internal/voting/service/...
```

## Contact

If you have questions or want to reach out, you can find me at:
- Email: [rafael.cr.carneiro@gmail.com]
