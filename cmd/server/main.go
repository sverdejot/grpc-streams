package main

import (
	"log"
	"log/slog"
	"net"
	"os"

	service "github.com/sverdejot/grpc-streams/internal/api/grpc"
	bidpb "github.com/sverdejot/grpc-streams/internal/api/grpc/bid/v1"
	"google.golang.org/grpc"
)

func main() {
	server := grpc.NewServer()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	auctionService := service.NewAuctionService()
	bidpb.RegisterAuctionServiceServer(server, auctionService)

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("cannot create net listened: %v", err)
	}
    defer lis.Close() // nolint: errcheck

	log.Printf("starting server at: %s\n", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
