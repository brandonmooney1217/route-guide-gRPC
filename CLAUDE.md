# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a gRPC route guide service implementation in Go that demonstrates all four gRPC streaming patterns. The project consists of:

- **Server** (`server/main.go`): Implements the RouteGuideServer with in-memory feature storage and location-based chat
- **Client** (`client/client.go`): Demonstrates all four RPC methods: GetFeature, ListFeatures, RecordRoute, and RouteChat
- **Protocol Buffers** (`routeguide/`): Service definition and generated Go code

## Project Structure
```
route-guide/
├── server/main.go          # gRPC server implementation
├── client/client.go        # Client demonstrating all RPC patterns
├── routeguide/             # Generated protobuf code
│   ├── routeguide.proto    # Service definition
│   ├── routeguide.pb.go    # Generated message types
│   └── routeguide_grpc.pb.go # Generated gRPC client/server code
└── CLAUDE.md              # This file
```

## Development Commands

### Running the Application
```bash
# Start the server (runs on port 50051)
go run server/main.go

# Run the client (in separate terminal)
go run client/client.go
```

### Testing and Building
```bash
# Build both server and client
go build ./...

# Test the complete flow
go run server/main.go &     # Start server in background
go run client/client.go     # Run client (will demonstrate all RPCs)
kill %1                     # Stop the background server
```

### Code Generation
```bash
# IMPORTANT: Regenerate protobuf Go code after modifying .proto files
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       routeguide/routeguide.proto
```

### Go Module Management
```bash
# Download dependencies
go mod tidy

# Build all packages
go build ./...
```

## Architecture

### gRPC Service Definition
The service is defined in `routeguide/routeguide.proto` with four RPC methods demonstrating all gRPC streaming patterns:
- `GetFeature(Point) returns (Feature)` - **Unary RPC**: Retrieves feature information for geographical coordinates
- `ListFeatures(Rectangle) returns (stream Feature)` - **Server streaming RPC**: Returns all features within a rectangle
- `RecordRoute(stream Point) returns (RouteSummary)` - **Client streaming RPC**: Client sends route points, server returns summary
- `RouteChat(stream RouteNote) returns (stream RouteNote)` - **Bidirectional streaming RPC**: Location-based chat system

### Server Implementation
- Uses pointer receivers for gRPC methods: `func (s *routeGuideServer) GetFeature(...)`
- Stores features in-memory via `savedFeatures []*pb.Feature` field
- Stores route notes in-memory via `routeNotes map[string][]*pb.RouteNote` field (location-based chat storage)
- Thread-safe access using `sync.Mutex` for concurrent RouteChat operations
- Implements all four RPC patterns:
  - **GetFeature**: Returns feature at given coordinates or empty feature
  - **ListFeatures**: Streams features within rectangle boundaries
  - **RecordRoute**: Counts points and features in client's route
  - **RouteChat**: Location-based messaging - stores and retrieves messages by coordinates
- Contains 7 hardcoded US landmarks for testing

### Key Go/gRPC Patterns
- All protobuf message types use pointers (`*pb.Point`, `*pb.Feature`, `*pb.Rectangle`, `*pb.RouteNote`)
- Server struct embeds `pb.UnimplementedRouteGuideServer` for default implementations
- Uses `&` operator to create pointers for protobuf message literals
- **Unary RPC**: Methods return `(message, error)` tuples
- **Server streaming**: Uses `stream.Send()` in loop, returns only `error`
- **Client streaming**: Uses `stream.Recv()` in loop until EOF, calls `stream.SendAndClose()`
- **Bidirectional streaming**: Uses both `stream.Send()` and `stream.Recv()`, coordinates with goroutines
- Serialization pattern: `serialize(point)` converts coordinates to string keys for map storage
- Synchronization: Channels (`waitc := make(chan struct{})`) coordinate goroutines in bidirectional streaming

## Common Development Workflows

### Adding a New RPC Method
1. Define the RPC in `routeguide/routeguide.proto`
2. Add corresponding message types if needed
3. Run protobuf code generation (see commands above)
4. Implement the method in `server/main.go`
5. Add client demonstration in `client/client.go`

### Modifying Server State
- The server uses in-memory storage only - no persistence
- `savedFeatures` contains the landmark data (initialized in `newServer()`)
- `routeNotes` stores chat messages by location (requires mutex for thread safety)
- All coordinate operations use the `serialize()` function for consistent string keys

## Important Notes

- The server has 7 hardcoded US landmark coordinates in `newServer()`
- Generated protobuf files (`*.pb.go`) should not be manually edited
- All four gRPC streaming patterns are implemented
- Coordinates are stored in E7 format (latitude/longitude × 10^7)
- `isFeatureInRectangle()` helper function validates if a point falls within rectangle boundaries
- RouteChat implements a location-based chat system where messages are stored by serialized coordinates
- Client code is organized with separate functions for each RPC method
- Thread safety is handled via mutex for concurrent access to shared route notes storage
- Server runs on port 50051 by default

## RouteChat Behavior

The RouteChat RPC implements a location-based messaging system:
1. When a client sends a message to coordinates (lat, lng), the server first sends back ALL previous messages sent to those same coordinates
2. Then the server stores the new message at that location
3. Multiple clients can "chat" by sending messages to the same geographical coordinates
4. Messages are stored indefinitely in server memory using serialized coordinates as map keys