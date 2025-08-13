package main

import (
	"context"
	"log"
	pb "routeguide/routeguide"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Establish connection to gRPC server
	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create RouteGuide client
	client := pb.NewRouteGuideClient(conn)

	point := &pb.Point{
		Latitude:  395906000,
		Longitude: -753506000,
	}

	// Set up context with timeout for the RPC call
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call GetFeature RPC method
	feature, err := client.GetFeature(ctx, point)
	if err != nil {
		log.Fatalf("GetFeature failed: %v", err)
	}

	// Extract and display feature information
	log.Printf("Feature name: %s, Latitude: %d, Longitude: %d",
		feature.Name, feature.Location.Latitude, feature.Location.Longitude)

	rect := &pb.Rectangle{
		BottomLeftCorner: &pb.Point{
			Latitude:  385000000,
			Longitude: -780000000,
		},
		TopRightCorner: &pb.Point{
			Latitude:  410000000,
			Longitude: -735000000,
		},
	}

	stream, err := client.ListFeatures(ctx, rect)

	if err != nil {
		log.Fatalf("ListFeature failed: %v", err)
	}

	log.Println("Features in rectangle:")
	for {
		feature, err = stream.Recv()
		if err != nil {
			break // End of stream (EOF is normal)
		}
		log.Printf("- %s at (%d, %d)", feature.Name, feature.Location.Latitude, feature.Location.Longitude)
	}

	points := []*pb.Point{
		{Latitude: 395906000, Longitude: -753506000},
		{Latitude: 405847500, Longitude: -741301800},
		{Latitude: 407486500, Longitude: -739885900},
		{Latitude: 407486500, Longitude: -3},
	}

	stream2, err := client.RecordRoute(ctx)
	if err != nil {
		log.Fatalf("Failed to record route: %v", err)
	}

	for _, point := range points {
		if err := stream2.Send(point); err != nil {
			log.Fatalf("Failed to send point: %v", err)
		}
	}

	summary, err := stream2.CloseAndRecv()
	if err != nil {
		log.Fatalf("Failed to receive route summary: %v", err)
	}
	log.Printf("Route summary: %d points, %d features",
		summary.PointCount, summary.FeatureCount)
}
