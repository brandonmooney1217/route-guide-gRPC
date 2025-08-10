# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a gRPC route guide service implementation in Go. The project consists of:

- **Server** (`server/main.go`): Implements the RouteGuideServer with in-memory feature storage
- **Client** (`client/client.go`): Demonstrates calling the GetFeature RPC method
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
The service is defined in `routeguide/routeguide.proto` with a single RPC method:
- `GetFeature(Point) returns (Feature)` - Retrieves feature information for geographical coordinates

### Server Implementation
- Uses pointer receivers for gRPC methods: `func (s *routeGuideServer) GetFeature(...)`
- Stores features in-memory via `savedFeatures []*pb.Feature` field
- Currently implements only `GetFeature`; other methods return "not implemented"

### Key Go/gRPC Patterns
- All protobuf message types use pointers (`*pb.Point`, `*pb.Feature`)
- Server struct embeds `pb.UnimplementedRouteGuideServer` for default implementations
- Uses `&` operator to create pointers for protobuf message literals
- gRPC methods require pointer receivers and return `(message, error)` tuples

## Important Notes

- The server has hardcoded Liberty Bell coordinates in `newServer()`
- Generated protobuf files (`*.pb.go`) should not be manually edited
- Only `GetFeature` is currently implemented - adding other RPC methods requires updating the .proto file first