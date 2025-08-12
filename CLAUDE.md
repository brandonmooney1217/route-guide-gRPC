# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a gRPC route guide service implementation in Go. The project consists of:

- **Server** (`server/main.go`): Implements the RouteGuideServer with in-memory feature storage
- **Client** (`client/client.go`): Demonstrates calling both GetFeature and ListFeatures RPC methods
- **Protocol Buffers** (`routeguide/`): Service definition and generated Go code

## Development Commands

### Running the Application
```bash
# Start the server (runs on port 50051)
go run server/main.go

# Run the client (in separate terminal)
go run client/client.go
```

### Code Generation
```bash
# Regenerate protobuf Go code after modifying .proto files
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
The service is defined in `routeguide/routeguide.proto` with two RPC methods:
- `GetFeature(Point) returns (Feature)` - Retrieves feature information for geographical coordinates
- `ListFeatures(Rectangle) returns (stream Feature)` - Server streaming RPC that returns all features within a rectangle

### Server Implementation
- Uses pointer receivers for gRPC methods: `func (s *routeGuideServer) GetFeature(...)`
- Stores features in-memory via `savedFeatures []*pb.Feature` field
- Implements both `GetFeature` (unary RPC) and `ListFeatures` (server streaming RPC)
- Contains 7 hardcoded US landmarks for testing

### Key Go/gRPC Patterns
- All protobuf message types use pointers (`*pb.Point`, `*pb.Feature`, `*pb.Rectangle`)
- Server struct embeds `pb.UnimplementedRouteGuideServer` for default implementations
- Uses `&` operator to create pointers for protobuf message literals
- Unary gRPC methods return `(message, error)` tuples
- Streaming gRPC methods use `stream.Send()` and return only `error`
- Client streaming uses `stream.Recv()` in loop until EOF

## Important Notes

- The server has 7 hardcoded US landmark coordinates in `newServer()`
- Generated protobuf files (`*.pb.go`) should not be manually edited
- Both `GetFeature` and `ListFeatures` are implemented
- Coordinates are stored in E7 format (latitude/longitude × 10^7)
- `isFeatureInRectangle()` helper function validates if a point falls within rectangle boundaries