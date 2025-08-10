package main

import (
	"context"
	"log"
	"net"
	pb "routeguide/routeguide"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

// server struct implements the RouteGuideServer interface
type routeGuideServer struct {
	pb.UnimplementedRouteGuideServer
	savedFeatures []*pb.Feature
}

func (s *routeGuideServer) GetFeature(_ context.Context, point *pb.Point) (*pb.Feature, error) {
	for _, feature := range s.savedFeatures {
		if proto.Equal(feature.Location, point) {
			return feature, nil
		}
	}
	return &pb.Feature{Location: point}, nil
}

func newServer() *routeGuideServer {
	return &routeGuideServer{
		savedFeatures: []*pb.Feature{
			{
				Name:     "Liberty Bell",
				Location: &pb.Point{Latitude: 395906000, Longitude: -753506000},
			},
		},
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterRouteGuideServer(s, newServer())

	log.Println("Server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
