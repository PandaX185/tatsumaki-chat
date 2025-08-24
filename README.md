# Tatsumaki Chat

A real-time chat application built with Go, featuring SSE communications, Redis pub/sub for realtime messaging and PostgreSQL for data persistence.

## 🚀 Features

- **User authentication** with JWT tokens
- **Real-time messaging** using Server-Sent Events (SSE)
- **Multi-user chat rooms** with member management
- **Unread message tracking** with real-time notifications
- **Message persistence** with PostgreSQL
- **Redis pub/sub** to track realtime events (new unread message, mark as read)

## 🏗️ Architecture

### Tech Stack

- **Language**: Go
- **Database**: PostgreSQL
- **Pub-Sub**: Redis
- **Real-time Communication**: Server-Sent Events (SSE)
- **HTTP Router**: Go's native `net/http` with `http.ServeMux`

## 🚀 Getting Started

### Prerequisites

- Go 1.23 or higher
- Docker and Docker Compose
- PostgreSQL (via Docker)
- Redis (via Docker)

### Environment Setup

1. Clone the repository:
```bash
git clone https://github.com/PandaX185/tatsumaki-chat.git
cd tatsumaki-chat
```

2. Create `.env` file:
```env
DB_USER=
DB_PASS=
DB_HOST=
DB_NAME=
DB_SSLMODE=
PORT=
REDIS_URL=
```

3. You can simply run:
```bash
make run-all
```

The server will start on `http://localhost:8080`


## 📝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License.
