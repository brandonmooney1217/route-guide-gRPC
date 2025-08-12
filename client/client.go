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

}
