package main

import (
	"context"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	bidpb "github.com/sverdejot/grpc-streams/internal/api/grpc/bid/v1"
	"github.com/sverdejot/grpc-streams/internal/tui"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
    conn, err := grpc.NewClient(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    f := bidpb.NewAuctionServiceClient(conn)
    fetcher, err := f.GetBids(context.Background(),&bidpb.GetBidsRequest{AuctionId: "auction-1"})
    if err != nil {
        log.Fatal(err)
    }
    ad := func() (any, error) {
        v, err := fetcher.Recv()
        return v.GetBid(), err
    }
    p := tea.NewProgram(tui.CreateAuction(ad))

    if _, err := p.Run(); err != nil {
        log.Fatalf("error while running TUI: %v", err)
    }
}
