# Voting-Go

A real-time voting platform where hosts can create polls and participants see results update instantly as votes come in. This project is a full-stack application leveraging Go's concurrency model and WebSockets for real-time communication.

---

## Project Overview

This repository is split into two main components:
- **Backend**: A Go-based REST & WebSocket API with PostgreSQL persistence and Redis caching.
- **Frontend**: A modern React 19 SPA built with Vite and TypeScript.

## Features

- **Real-time Synchronization**: Result updates are broadcasted instantly via WebSockets.
- **Secure Authentication**: User accounts and JWT-protected endpoints.
- **Concurrent-safe**: High-performance voting logic using Go's synchronization primitives.
- **Responsive UI**: A clean, modern interface for creating and participating in polls.

## Technical Stack

### Backend
- **Language**: Go (v1.26+)
- **Database**: PostgreSQL
- **Cache**: Redis
- **Auth**: JWT
- **Real-time**: Gorilla WebSocket

### Frontend
- **Framework**: React 19 (Vite + TypeScript)
- **Icons**: Lucide React
- **Styling**: Vanilla CSS (CSS Modules)

## Getting Started

Detailed instructions for setting up each component can be found in their respective directories:

- [Backend Documentation](./backend/README.md)
- [Frontend Documentation](./frontend/README.md)

### Quick Run (Docker)
*Ensure you have Docker and Docker Compose installed.*

```bash
docker-compose up --build
```

## Contact

If you have questions or want to reach out:
- Email: [rafael.cr.carneiro@gmail.com]
