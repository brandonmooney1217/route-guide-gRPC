# Route Guide

A comprehensive gRPC route guide service implementation in Go that demonstrates all four gRPC streaming patterns. This project serves as a complete example of unary, server streaming, client streaming, and bidirectional streaming RPCs.

## Overview

This project demonstrates a complete gRPC client-server application using Protocol Buffers. The server maintains an in-memory collection of geographical features and implements a location-based chat system. The client demonstrates all four RPC patterns with clear examples of each streaming type.

## Features

- **Unary RPC**: GetFeature - Retrieve feature information by coordinates
- **Server Streaming**: ListFeatures - Stream all features within a geographical rectangle  
- **Client Streaming**: RecordRoute - Send route points and receive summary statistics
- **Bidirectional Streaming**: RouteChat - Location-based messaging system

## Project Structure

```
route-guide/
├── client/
│   └── client.go          # Client demonstrating all four RPC patterns
├── server/
│   └── main.go           # Complete gRPC server implementation
├── routeguide/
│   ├── routeguide.proto  # Service definition with all four RPC types
│   ├── routeguide.pb.go  # Generated protobuf Go code
│   └── routeguide_grpc.pb.go # Generated gRPC Go code
├── go.mod
├── go.sum
├── README.md
└── CLAUDE.md             # Development guidance
```

## Prerequisites

- Go 1.19 or later
- Protocol Buffer compiler (`protoc`) - only needed if modifying .proto files
- `protoc-gen-go` and `protoc-gen-go-grpc` plugins

## Getting Started

### 1. Install Dependencies

```bash
go mod tidy
```

### 2. Run the Server

```bash
go run server/main.go
```

The server will start listening on port 50051 and log:
```
Server listening on :50051
```

### 3. Run the Client

In a separate terminal:

```bash
go run client/client.go
```

The client will demonstrate all four RPC patterns:
1. GetFeature - Query for Liberty Bell coordinates
2. ListFeatures - List all features in an East Coast rectangle
3. RecordRoute - Send a route and receive statistics
4. RouteChat - Exchange location-based messages

## Service Definition

The RouteGuide service provides four RPC methods demonstrating all gRPC streaming patterns:

### RPC Methods

- `GetFeature(Point) returns (Feature)` - **Unary**: Retrieves feature information for given coordinates
- `ListFeatures(Rectangle) returns (stream Feature)` - **Server Streaming**: Streams all features within a geographical rectangle
- `RecordRoute(stream Point) returns (RouteSummary)` - **Client Streaming**: Accepts route points and returns summary statistics
- `RouteChat(stream RouteNote) returns (stream RouteNote)` - **Bidirectional Streaming**: Location-based chat system

### Message Types

- `Point` - Geographical coordinates (latitude, longitude in E7 format)
- `Feature` - Feature name and location
- `Rectangle` - Geographical boundary with corners
- `RouteSummary` - Statistics about a route (point count, feature count)
- `RouteNote` - Chat message with location and text

## RouteChat Feature

The RouteChat RPC implements a unique location-based messaging system:

- Send messages attached to specific geographical coordinates
- Receive all previous messages sent to the same location
- Multiple clients can chat by using the same coordinates
- Messages are stored in server memory using serialized coordinates as keys

## Development

### Quick Test

```bash
# Test the complete system
go run server/main.go &     # Start server in background
go run client/client.go     # Run client (demonstrates all RPCs)
kill %1                     # Stop the background server
```

### Modifying the Protocol Buffer Definition

If you need to modify the service definition:

1. Edit `routeguide/routeguide.proto`
2. Regenerate the Go code:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       routeguide/routeguide.proto
```

### Building

```bash
# Build all packages
go build ./...

# Build specific components
go build ./server
go build ./client
```

## Example Output

When running the client, you should see output similar to:

```
=== GetFeature ===
Feature name: Liberty Bell, Latitude: 395906000, Longitude: -753506000

=== ListFeatures ===
Features in rectangle:
- Liberty Bell at (395906000, -753506000)
- Statue of Liberty at (405847500, -741301800)
- Empire State Building at (407486500, -739885900)
- Lincoln Memorial at (389030600, -770494800)

=== RecordRoute ===
Route summary: 4 points, 3 features

=== RouteChat ===
Sending message: First message at location: (0, 1)
Sending message: Second message at location: (0, 2)
Received message: First message at location: (0, 1)
Sending message: Fourth message at location: (0, 1)
Received message: First message at location: (0, 1)
...
```

## Architecture Highlights

- **Thread Safety**: Server uses mutex for concurrent access to shared route notes storage
- **Streaming Patterns**: Complete implementation of all four gRPC streaming types
- **Client Organization**: Each RPC method implemented in separate functions for clarity
- **Error Handling**: Proper EOF handling for streaming operations
- **Synchronization**: Channel-based coordination for bidirectional streaming

## Current Implementation

- Server uses in-memory data storage (7 hardcoded US landmarks)
- Location-based chat system with persistent message storage
- No authentication or authorization
- Single server instance (no clustering)

## Dependencies

- `google.golang.org/grpc` - gRPC Go implementation
- `google.golang.org/protobuf` - Protocol Buffers Go implementation

## Learning Objectives

This project demonstrates:
- Protocol Buffer service definitions
- All four gRPC streaming patterns
- Go concurrency patterns with goroutines and channels
- Thread-safe server implementation
- Client-server communication patterns
- Error handling in streaming operations