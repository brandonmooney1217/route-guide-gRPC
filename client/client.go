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
	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// create client

	client := pb.NewRouteGuideClient(conn)

	var point *pb.Point = &pb.Point{Latitude: 395906000, Longitude: -753506000}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	feature, err := client.GetFeature(ctx, point)

	if err != nil {
		log.Fatalf("GetFeature failed: %v", err)
	}

	var name string = feature.Name
	var location *pb.Point = feature.Location
	var latitude int32 = location.Latitude
	longitude := location.Longitude
	log.Printf("%s", name)
	log.Printf("Latitude: %v", latitude)
	log.Printf("Longitude: %v", longitude)

}
