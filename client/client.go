package main

import (
	"context"
	"io"
	"log"
	pb "routeguide/routeguide"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getFeature(client pb.RouteGuideClient, ctx context.Context) {
	log.Println("=== GetFeature ===")
	point := &pb.Point{
		Latitude:  395906000,
		Longitude: -753506000,
	}

	feature, err := client.GetFeature(ctx, point)
	if err != nil {
		log.Fatalf("GetFeature failed: %v", err)
	}

	log.Printf("Feature name: %s, Latitude: %d, Longitude: %d",
		feature.Name, feature.Location.Latitude, feature.Location.Longitude)
}

func listFeatures(client pb.RouteGuideClient, ctx context.Context) {
	log.Println("=== ListFeatures ===")
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
		feature, err := stream.Recv()
		if err != nil {
			break // End of stream (EOF is normal)
		}
		log.Printf("- %s at (%d, %d)", feature.Name, feature.Location.Latitude, feature.Location.Longitude)
	}
}

func recordRoute(client pb.RouteGuideClient, ctx context.Context) {
	log.Println("=== RecordRoute ===")
	points := []*pb.Point{
		{Latitude: 395906000, Longitude: -753506000},
		{Latitude: 405847500, Longitude: -741301800},
		{Latitude: 407486500, Longitude: -739885900},
		{Latitude: 407486500, Longitude: -3},
	}

	stream, err := client.RecordRoute(ctx)
	if err != nil {
		log.Fatalf("Failed to record route: %v", err)
	}

	for _, point := range points {
		if err := stream.Send(point); err != nil {
			log.Fatalf("Failed to send point: %v", err)
		}
	}

	summary, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Failed to receive route summary: %v", err)
	}
	log.Printf("Route summary: %d points, %d features",
		summary.PointCount, summary.FeatureCount)
}

func routeChat(client pb.RouteGuideClient, ctx context.Context) {
	log.Println("=== RouteChat ===")
	notes := []*pb.RouteNote{
		{Message: "First message", Location: &pb.Point{Latitude: 0, Longitude: 1}},
		{Location: &pb.Point{Latitude: 0, Longitude: 2}, Message: "Second message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 3}, Message: "Third message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 1}, Message: "Fourth message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 2}, Message: "Fifth message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 3}, Message: "Sixth message"},
	}

	stream, err := client.RouteChat(ctx)
	if err != nil {
		log.Fatalf("client.RouteChat failed: %v", err)
	}

	waitc := make(chan struct{})
	go func() {
		for {
			input, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("client.RouteChat failed: %v", err)
			}
			log.Printf("Received message: %v at location: (%d, %d)",
				input.Message, input.Location.Latitude, input.Location.Longitude)
		}
	}()

	for _, note := range notes {
		log.Printf("Sending message: %v at location: (%d, %d)",
			note.Message, note.Location.Latitude, note.Location.Longitude)
		if err := stream.Send(note); err != nil {
			log.Fatalf("Failed to send note: %v", err)
		}
	}
	stream.CloseSend()
	<-waitc
}

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

	// Set up context with timeout for the RPC calls
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Call each RPC method
	getFeature(client, ctx)
	listFeatures(client, ctx)
	recordRoute(client, ctx)
	routeChat(client, ctx)
}