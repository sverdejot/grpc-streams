package grpc

import (
	"context"
	"math/rand/v2"
	"fmt"
	"time"

	bidpb "github.com/sverdejot/grpc-streams/internal/api/grpc/bid/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type auctionService struct {
    bidpb.UnimplementedAuctionServiceServer
}

func NewAuctionService() *auctionService {
    return &auctionService{}
}

func (s *auctionService) CreateBid(_ context.Context, req *bidpb.CreateBidRequest) (*bidpb.CreateBidResponse, error) {
    fmt.Printf("create bid: %v\n", req)
    return nil, nil
}

func (s *auctionService) GetBids(req *bidpb.GetBidsRequest, stream bidpb.AuctionService_GetBidsServer) (error) {
    tc := time.NewTicker(1 * time.Second)
    for {
        select {
        case <-stream.Context().Done():
            return status.Error(codes.Canceled, "stream closed")
        case <-tc.C:
            if err := stream.SendMsg(&bidpb.GetBidsResponse{
                Bid: &bidpb.Bid{
                    UserId: "some-random-user-id",
                    AuctionId: req.AuctionId,
                    QtyInCents: rand.Int32N(10000),
                },
            }); err != nil {
                return status.Error(codes.Internal, fmt.Errorf("error while sending msg: %w", err).Error())
            }
        }
    }
}
