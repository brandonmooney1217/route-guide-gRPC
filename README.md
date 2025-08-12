# Route Guide

A simple gRPC route guide service implementation in Go that allows clients to retrieve geographical feature information.

## Overview

This project demonstrates a basic gRPC client-server application using Protocol Buffers. The server maintains an in-memory collection of geographical features and provides an RPC endpoint to query features by location coordinates.

## Project Structure

```
route-guide/
├── client/
│   └── client.go          # gRPC client implementation
├── server/
│   └── main.go           # gRPC server implementation
├── routeguide/
│   ├── routeguide.proto  # Protocol Buffer service definition
│   ├── routeguide.pb.go  # Generated protobuf Go code
│   └── routeguide_grpc.pb.go # Generated gRPC Go code
├── go.mod
└── go.sum
```

## Prerequisites

- Go 1.24.6 or later
- Protocol Buffer compiler (`protoc`) - only needed if modifying .proto files

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

The client will connect to the server, request feature information for the Liberty Bell coordinates, then list all features within an East Coast rectangle.

## Service Definition

The RouteGuide service provides two RPC methods:

- `GetFeature(Point) returns (Feature)` - Retrieves feature information for given coordinates
- `ListFeatures(Rectangle) returns (stream Feature)` - Streams all features within a geographical rectangle

### Message Types

- `Point` - Represents geographical coordinates (latitude, longitude)
- `Feature` - Contains a feature name and its location
- `Rectangle` - Defines a geographical boundary with bottom-left and top-right corners

## Development

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
Feature name: Liberty Bell, Latitude: 395906000, Longitude: -753506000
Features in rectangle:
- Liberty Bell at (395906000, -753506000)
- Statue of Liberty at (405847500, -741301800)
- Empire State Building at (407486500, -739885900)
- Lincoln Memorial at (389030600, -770494800)
```

## Current Limitations

- Server uses hardcoded in-memory data (7 US landmarks)
- No persistent storage
- No authentication or authorization

## Dependencies

- `google.golang.org/grpc` - gRPC Go implementation
- `google.golang.org/protobuf` - Protocol Buffers Go implementation