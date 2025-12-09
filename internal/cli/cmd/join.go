package cmd

import (
    "github.com/spf13/cobra"
	"context"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	bidpb "github.com/sverdejot/grpc-streams/internal/api/grpc/bid/v1"
	"github.com/sverdejot/grpc-streams/internal/cli/tui"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var joinCmd = &cobra.Command{
    Use: "join [auction id] [full name]",
    Short: "Join an auction",
    Args: cobra.ExactArgs(2),
    Run: func(cmd *cobra.Command, args []string) {
        auctionID := args[0]
        userID := args[1]

        join(auctionID, userID)
    },
}

func join(auctionID, userID string) {
	conn, err := grpc.NewClient(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
    defer conn.Close() // nolint: errcheck
	f := bidpb.NewAuctionServiceClient(conn)
	fetcher, err := f.GetBids(context.Background(), &bidpb.GetBidsRequest{AuctionId: auctionID})
	if err != nil {
		log.Fatal(err)
	}

	creator := func(userID string, qty int) error {
		req := &bidpb.CreateBidRequest{
			UserId:          userID,
			QuantityInCents: int32(qty),
			AuctionId:       "019b0090-3f7b-7288-913b-64b54f2f4133",
		}

		_, err := f.CreateBid(context.Background(), req)
		return err
	}
	ad := func() (any, error) {
		v, err := fetcher.Recv()
		return v.GetBid(), err
	}
	p := tea.NewProgram(tui.CreateAuction(ad, creator, userID))

	if _, err := p.Run(); err != nil {
		log.Fatalf("error while running TUI: %v", err)
	}
}
