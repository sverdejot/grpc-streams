package main

import (
	"log"
	"net"

	bidpb "github.com/sverdejot/grpc-streams/internal/api/grpc/bid/v1"
	service "github.com/sverdejot/grpc-streams/internal/api/grpc"
	"google.golang.org/grpc"
)

func main() {
    server := grpc.NewServer()

    auctionService := service.NewAuctionService()
    bidpb.RegisterAuctionServiceServer(server, auctionService)

    lis, err := net.Listen("tcp", ":8080")
    if err != nil {
        log.Fatalf("cannot create net listened: %v", err)
    }
    defer lis.Close()

    log.Printf("starting server at: %s\n", lis.Addr())
    if err := server.Serve(lis); err != nil {
        log.Fatalf("server stopped: %v", err)
    }
}
