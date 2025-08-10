package main

import (
	"context"
	"log"
	"net"
	pb "routeguide/routeguide"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

// routeGuideServer implements the RouteGuideServer interface
type routeGuideServer struct {
	pb.UnimplementedRouteGuideServer
	savedFeatures []*pb.Feature // in-memory storage for geographical features
}

// GetFeature retrieves a feature at the given geographical point
// Returns the named feature if found, otherwise returns a feature with empty name
func (s *routeGuideServer) GetFeature(_ context.Context, point *pb.Point) (*pb.Feature, error) {
	for _, feature := range s.savedFeatures {
		if proto.Equal(feature.Location, point) {
			return feature, nil
		}
	}
	return &pb.Feature{Location: point}, nil
}

// newServer creates and initializes a new RouteGuide server instance
// with hardcoded Liberty Bell coordinates as sample data
func newServer() *routeGuideServer {
	return &routeGuideServer{
		savedFeatures: []*pb.Feature{
			{
				Name:     "Liberty Bell",
				Location: &pb.Point{Latitude: 395906000, Longitude: -753506000}, // Philadelphia coordinates in E7 format
			},
		},
	}
}

func main() {
	// Create TCP listener on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create gRPC server and register our RouteGuide service
	s := grpc.NewServer()
	pb.RegisterRouteGuideServer(s, newServer())

	log.Println("Server listening on :50051")

	// Start serving requests
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
