# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

LibreRiichi is an open-source Mahjong server written in Go with a Vue.js frontend. The project implements a real-time multiplayer Riichi Mahjong game using WebSocket communication.

## Architecture

### Backend (Go)
- **Entry Point**: `cmd/server/main.go` - Starts WebSocket server on port 3000
- **Core Package**: Contains the main game logic
  - `arena_manager.go` - Manages multiple game rooms/arenas with thread-safe operations
  - `arena.go` - Individual game room implementation
  - `game.go` - Core game state and Mahjong rule logic
  - `client.go` - WebSocket client connection handling
  - `game_data/` - Mahjong-specific data structures (tiles, hands, yaku, etc.)
  - `messages/` - Message protocol definitions and auto-generated serialization code
  - `web/` - HTTP/WebSocket server implementation
  - `util/` - Helper utilities

### Frontend (Vue.js/TypeScript)
- **Entry Point**: `app/index.html` served by Parcel
- **Structure**: 
  - `app/components/` - Vue components
  - `app/views/` - Page views
  - `app/messaging/` - WebSocket communication layer
  - `app/game/` - Game-specific UI logic
  - Uses Pinia for state management and Vue Router for navigation

### Message Generation
The project uses code generation for message serialization:
- `core/generate_message.go` - Generator that creates `*_generated.go` files
- Generated files handle message marshaling/unmarshaling between Go and frontend

## Development Commands

### Backend
```bash
# Run development server
yarn backend
# or
go run cmd/server/main.go

# Build production server
yarn build-backend
# or  
go build -o dist-server/server cmd/server/main.go

# Run tests
go test ./...

# Run specific test file
go test ./core/game_data/

# Test with evaluation client
yarn test-client
# or
go run cmd/eval_client/eval_client.go
```

### Frontend
```bash
# Run development server with hot reload
yarn frontend

# Build production frontend
yarn build-frontend

# Run both frontend and backend concurrently
yarn watch

# Run full test setup (backend + test client)
yarn test
```

### Message Generation
After modifying message structures in `core/messages/`, regenerate serialization code:
```bash
cd core
go generate ./messages/
```

## Testing
- Go tests: Use standard `go test` command
- Test files follow `*_test.go` naming convention
- Evaluation client (`cmd/eval_client/`) provides automated testing of game scenarios

## Key Dependencies
- **Go**: WebSocket (`gorilla/websocket`), UUID generation (`google/uuid`)
- **Frontend**: Vue 3, TypeScript, Tailwind CSS, Pinia, Parcel bundler

## IMPORTANT: Sound Notification

After finishing responding to my request or running a command, run this command to notify me by sound:

```bash
afplay /System/Library/Sounds/Funk.aiff
```
