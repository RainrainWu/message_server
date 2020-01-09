# message_server
This is the final project of **Network Management** course, develop with go language.

## Getting Started

### Prerequisites
- go 1.13+ (GOMODULE111=true)

## Architecture

### Server
Implemented in server.go, maintain a connected client list and deal with dispatching message between them.

### Client
Implemented in client.go, the terminal to send message and chat with others.

## Features
- Websocket communication

### Server side
- Concurrent processing
	- Websocket read, write goroutine for each user
	- Broadcast goroutine
- Multiple users supported
- Message broadcasting
- Bad user banned

### Client side
- Concurrent processing
	- Websocket read, write goroutine
	- Command line input, output goroutine
- User debug mode
- Sensitive string masked

## How to use
### Server
```bash
	go run server.go
```

### Client
```bash
	go run client.go -addr=SERVER_ADDR -debug=false
```
- **-addr** : the message server address.
- **-debug** : whether enable debug mode.
