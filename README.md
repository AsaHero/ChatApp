# ChatApp

ChatApp is a simple chat application written in Go, developed for educational purposes. It serves as a demonstration of clean architecture principles, database connections, WebSocket implementation, and the utilization of Go's concurrency and channels. The application is built using the Fiber framework for HTTP server.

## Features

- **Clean Architecture:** The project is structured following clean architecture principles, ensuring separation of concerns and maintainability.
- **Database Connections:** Demonstrates the integration of database connections for storing and retrieving chat-related data.
- **WebSocket Implementation:** Utilizes WebSockets to enable real-time communication between clients.
- **Go Concurrency and Channels:** Showcases the power of Go through its concurrency features and the usage of channels for communication between goroutines.
- **Fiber Framework:** Showcases the power of [Fiber](https://github.com/gofiber) runnnig HTTP API server. 

## Prerequisites

To run this application, ensure you have the following installed:

- Go (version 1.19+)
- Other dependencies as specified in the `go.mod` file.

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/AsaHero/ChatApp.git

## Run

   ```bash
   go run cmd/main.go
   
