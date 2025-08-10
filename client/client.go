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
}
