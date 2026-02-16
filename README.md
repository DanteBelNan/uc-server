# Universal Clipboard Server

A real-time, cross-platform clipboard synchronization solution that enables instant sharing of clipboard content across multiple devices connected to the same secure room.

## Overview

Universal Clipboard Server (`uc-server`) is the central hub that orchestrates clipboard synchronization between connected clients. It manages WebSocket connections, room-based authentication, and message broadcasting to enable seamless clipboard sharing across desktop and mobile devices.

## Features

- **Real-time Synchronization**: Instant clipboard updates across all connected devices
- **Room-based Isolation**: Secure, password-protected rooms for private clipboard sharing
- **Low Latency Communication**: WebSocket-based bidirectional communication
- **Scalable Architecture**: Redis-backed pub/sub for horizontal scaling
- **End-to-End Encryption Ready**: Designed to support E2EE with AES-256

## Architecture

The server follows a **Hexagonal Architecture** pattern:

- **Core**: Domain entities and business logic
- **Ports**: Interface definitions (input/output)
- **Adapters**: Infrastructure implementations (WebSocket, Redis, HTTP)

## Tech Stack

- **Language**: Go 1.25.5
- **Web Framework**: Gin
- **WebSocket**: Gorilla WebSocket
- **State Management**: Redis
- **Development**: Air (hot-reload), Docker

## Getting Started

### Prerequisites

- Docker & Docker Compose
- Go 1.23.5+ (for local development)
- Make

### Quick Start

1. Clone the repository
2. Run the development environment:
```bash
make dev
```

The server will start on `http://localhost:8080`

### Available Commands
```bash
make help          # Show all available commands
make dev           # Start development server with hot-reload
make build         # Build production binary
make test          # Run tests
make docker-build  # Build production Docker image
```

## API Endpoints

### Health Check
```
GET /health
```

Returns server status.

## License

MIT