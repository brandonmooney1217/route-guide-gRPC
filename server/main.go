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

func (s *routeGuideServer) ListFeatures(rect *pb.Rectangle, stream pb.RouteGuide_ListFeaturesServer) error {
	for _, feature := range s.savedFeatures {
		if isFeatureInRectangle(rect, feature) {
			if err := stream.Send(feature); err != nil {
				return err
			}
		}
	}
	return nil
}

func isFeatureInRectangle(rect *pb.Rectangle, feature *pb.Feature) bool {
	lat := feature.Location.Latitude
	lon := feature.Location.Longitude

	return lat > rect.BottomLeftCorner.Latitude && lat < rect.TopRightCorner.Latitude &&
		lon > rect.BottomLeftCorner.Longitude && lon < rect.TopRightCorner.Longitude
}

// newServer creates and initializes a new RouteGuide server instance
// with hardcoded Liberty Bell coordinates as sample data
func newServer() *routeGuideServer {
	return &routeGuideServer{
		savedFeatures: []*pb.Feature{
			{
				Name:     "Liberty Bell",
				Location: &pb.Point{Latitude: 395906000, Longitude: -753506000},
			},
			{
				Name:     "Statue of Liberty",
				Location: &pb.Point{Latitude: 405847500, Longitude: -741301800},
			},
			{
				Name:     "Empire State Building",
				Location: &pb.Point{Latitude: 407486500, Longitude: -739885900},
			},
			{
				Name:     "Golden Gate Bridge",
				Location: &pb.Point{Latitude: 378197400, Longitude: -1224650700},
			},
			{
				Name:     "Lincoln Memorial",
				Location: &pb.Point{Latitude: 389030600, Longitude: -770494800},
			},
			{
				Name:     "Mount Rushmore",
				Location: &pb.Point{Latitude: 438813500, Longitude: -1031032800},
			},
			{
				Name:     "Space Needle",
				Location: &pb.Point{Latitude: 476203100, Longitude: -1221315600},
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
