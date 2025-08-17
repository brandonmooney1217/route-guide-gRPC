package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	pb "routeguide/routeguide"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

// routeGuideServer implements the RouteGuideServer interface
type routeGuideServer struct {
	pb.UnimplementedRouteGuideServer
	savedFeatures []*pb.Feature              // in-memory storage for geographical features
	routeNotes    map[string][]*pb.RouteNote // in-memory storage for map of route notes at each point; use serialized point as the key
	mu            sync.Mutex                 // mutex for thread-safe access to routeNotes
}

// findFeatureAtPoint checks if a point exists in the saved features list
// Returns the feature if found, nil otherwise
func (s *routeGuideServer) findFeatureAtPoint(point *pb.Point) *pb.Feature {
	for _, feature := range s.savedFeatures {
		if proto.Equal(feature.Location, point) {
			return feature
		}
	}
	return nil
}

// GetFeature retrieves a feature at the given geographical point
// Returns the named feature if found, otherwise returns a feature with empty name
func (s *routeGuideServer) GetFeature(_ context.Context, point *pb.Point) (*pb.Feature, error) {
	if feature := s.findFeatureAtPoint(point); feature != nil {
		return feature, nil
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

func (s *routeGuideServer) RecordRoute(stream pb.RouteGuide_RecordRouteServer) error {
	var point_count, feature_count int32

	for {
		point, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.RouteSummary{
				PointCount:   point_count,
				FeatureCount: feature_count,
			})
		}
		if err != nil {
			return err
		}

		point_count = point_count + 1
		if s.findFeatureAtPoint(point) != nil {
			feature_count = feature_count + 1
		}
	}
}

// Receives a stream of Route Notes, which is Point Message pair, and returns back
// stream of all route notes at that location

func serialize(point *pb.Point) string {
	return fmt.Sprintf("%d,%d", point.Latitude, point.Longitude)
}

func (s *routeGuideServer) RouteChat(stream pb.RouteGuide_RouteChatServer) error {
	for {

		// 1. Process route note
		note, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		// 2. serialize note using the location
		key := serialize(note.Location)

		// 3. Lock
		s.mu.Lock()
		// 4. add to route notes
		s.routeNotes[key] = append(s.routeNotes[key], note)

		// 5. create copy
		rn := make([]*pb.RouteNote, len(s.routeNotes[key]))
		copy(rn, s.routeNotes[key])
		s.mu.Unlock()

		// 6. write stream
		for _, note := range rn {
			if err := stream.Send(note); err != nil {
				return err
			}
		}

	}

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
		routeNotes: make(map[string][]*pb.RouteNote),
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
