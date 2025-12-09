package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	bidpb "github.com/sverdejot/grpc-streams/internal/api/grpc/bid/v1"
	"github.com/sverdejot/grpc-streams/internal/api/internal"
	"github.com/sverdejot/grpc-streams/internal/api/internal/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
    defaultAuctionDuration = time.Duration(time.Hour)
)

type auctionService struct {
	bidpb.UnimplementedAuctionServiceServer
    s internal.AuctionStore
    n *internal.AuctionUpdatesNotificator
}

func NewAuctionService() *auctionService {
	return &auctionService{
        s: store.NewMemoryStore(),
        n: internal.NewNotificator(),
	}
}

func (s *auctionService) CreateAuction(_ context.Context, req *bidpb.CreateAuctionRequest) (*bidpb.CreateAuctionResponse, error) {
    a, err := internal.NewAuction(
        internal.Withitem(req.Item),
        internal.WithEndsAt(defaultAuctionDuration),
    )
    if err != nil {
        return nil, status.Errorf(codes.Internal, "cannot create auction: %s", err)
    }

    if err := s.s.Upsert(a); err != nil {
        return nil, status.Errorf(codes.Internal, "cannot save auction: %s", err)
    }

    return &bidpb.CreateAuctionResponse{
        Auction: a.ToProto(),
    }, nil
}

func (s *auctionService) CreateBid(_ context.Context, req *bidpb.CreateBidRequest) (*bidpb.CreateBidResponse, error) {
    a, ok := s.s.Get(req.AuctionId)
    if !ok {
        return nil, status.Errorf(codes.NotFound, "cannot find auction")
    }

    b := internal.NewBid(
        internal.WithAuctionID(req.AuctionId),
        internal.WithUserID(req.UserId),
        internal.WithQuantity(req.QuantityInCents),
    )

    ok = a.InsertBid(b)
    if !ok {
        return nil, status.Errorf(codes.Internal, "cannot insert bid")
    }
    
    err := s.s.Upsert(a) 
    if err != nil {
        return nil, status.Errorf(codes.Internal, "cannot save auction: %s", err)
    }

    s.n.Notify(a)

    return nil, nil
}

func (s *auctionService) GetBids(req *bidpb.GetBidsRequest, stream bidpb.AuctionService_GetBidsServer) error {
    ctx, cancel := context.WithCancel(stream.Context())
    defer cancel()
    ch, err :=s.n.Subscribe(req.AuctionId)
    if err != nil {
        return status.Errorf(codes.Internal, "cannot subscribe to bid stream: %s", err)
    }
	for {
		select {
		case <-stream.Context().Done():
            slog.InfoContext(ctx, "stream closed")
            close(ch)
			return status.Error(codes.Canceled, "stream closed")
		case bid := <-ch:
			slog.InfoContext(ctx, "sending bid", "auction_id", bid.AuctionID, "user_id", bid.UserID, "quantity", bid.QuantityInCents)
			if err := stream.SendMsg(&bidpb.GetBidsResponse{
				Bid: bid.ToProto(),
			}); err != nil {
				return status.Error(codes.Internal, fmt.Errorf("error while sending msg: %w", err).Error())
			}
		}
	}
}
