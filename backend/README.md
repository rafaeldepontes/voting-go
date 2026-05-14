# Voting-Go

Project built with [Gini](https://gini-webserver.up.railway.app/)

This project provides a platform where a host can create a poll and voters can see the results update instantly as votes come in. It is built using Go and utilizes WebSockets for real-time data broadcasting.

---

## Features

- **User Authentication**: Secure register and login using JWT.
- **Real-time Updates**: Instant result updates via WebSockets.
- **Persistent Storage**: Data persistence using PostgreSQL.
- **Caching**: Efficient session and data management with Redis.
- **Concurrent-safe**: Robust voting logic using sync primitives.
- **Automated Migrations**: Database schema management with `golang-migrate`.

## Technical Stack

- **Language**: Go (v1.26+)
- **Database**: PostgreSQL (pgx/v5)
- **Cache**: Redis (go-redis/v9)
- **Auth**: JWT (golang-jwt/v5)
- **WebSocket**: Gorilla WebSocket
- **Routing**: Standard Library http.ServeMux (Go 1.22+ patterns)

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/rafaeldepontes/voting-go.git
   cd voting-go/backend
   ```

2. Configure environment:
   ```bash
   cp .env.example .env
   # Edit .env with your local credentials
   ```

3. Download dependencies:
   ```bash
   go mod tidy
   ```

## Usage

### Starting the Server

```bash
go run cmd/server/main.go
```
The server defaults to port 8080.

### Authentication

#### Register
```bash
curl -X POST http://localhost:8080/register \
     -d '{"email": "user@example.com", "password": "yourpassword"}'
```

#### Login
```bash
curl -X POST http://localhost:8080/login \
     -d '{"email": "user@example.com", "password": "yourpassword"}'
```
*Note: Most subsequent requests require the `Authorization: Bearer <token>` header.*

### List all polls

```bash
curl -H "Authorization: Bearer <token>" http://localhost:8080/polls
```

### Creating a Poll

```bash
curl -X POST http://localhost:8080/polls \
     -H "Authorization: Bearer <token>" \
     -d '{"name": "The Best Programming Language", "options": ["Golang", "Rust", "Kotlin", "C#"]}'
```

### Voting

```bash
curl -X POST http://localhost:8080/polls/{id}/vote \
     -H "Authorization: Bearer <token>" \
     -d '{"optionId": 1}'
```

### Real-Time Monitoring

Connect to the WebSocket endpoint for a specific poll. The token must be passed as a query parameter:
`ws://localhost:8080/ws/polls/{id}?token=<token>`

## Testing

Run unit tests for the core voting service:
```bash
go test ./internal/voting/service/...
```

## Contact

If you have questions or want to reach out, you can find me at:
- Email: [rafael.cr.carneiro@gmail.com]
