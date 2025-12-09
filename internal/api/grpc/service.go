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

func (s *auctionService) CreateAuction(ctx context.Context, req *bidpb.CreateAuctionRequest) (*bidpb.CreateAuctionResponse, error) {
    slog.InfoContext(ctx, "creating auction", "item_name", req.Item)
    a, err := internal.NewAuction(
        internal.Withitem(req.Item),
        internal.WithEndsAt(defaultAuctionDuration),
    )
    if err != nil {
        slog.ErrorContext(ctx, "error creating auction", "item_name", req.Item)
        return nil, status.Errorf(codes.Internal, "cannot create auction: %s", err)
    }

    if err := s.s.Upsert(a); err != nil {
        slog.ErrorContext(ctx, "error saving auction", "auction_id", a.ID, "item_name", req.Item)
        return nil, status.Errorf(codes.Internal, "cannot save auction: %s", err)
    }

    slog.InfoContext(ctx, "auction created", "auction_id", a.ID,"item_name", req.Item)
    return &bidpb.CreateAuctionResponse{
        Auction: a.ToProto(),
    }, nil
}

func (s *auctionService) CreateBid(ctx context.Context, req *bidpb.CreateBidRequest) (*bidpb.CreateBidResponse, error) {
    slog.InfoContext(ctx, "creating bid", "auction_id", req.AuctionId, "user_it", req.UserId, "qty", req.QuantityInCents)
    a, ok := s.s.Get(req.AuctionId)
    if !ok {
        slog.ErrorContext(ctx, "cannot find auction to bid", "auction_id", req.AuctionId, "user_id", req.UserId, "qty", req.QuantityInCents)
        return nil, status.Errorf(codes.NotFound, "cannot find auction")
    }

    b := internal.NewBid(
        internal.WithAuctionID(req.AuctionId),
        internal.WithUserID(req.UserId),
        internal.WithQuantity(req.QuantityInCents),
    )

    ok = a.InsertBid(b)
    if !ok {
        slog.ErrorContext(ctx, "error inserting bid", "auction_id", a.ID, "user_id", req.UserId, "qty", req.QuantityInCents)
        return nil, status.Errorf(codes.Internal, "cannot insert bid")
    }
    
    err := s.s.Upsert(a) 
    if err != nil {
        slog.InfoContext(ctx, "error saving auction after bid", "auction_id", req.AuctionId, "user_it", req.UserId, "qty", req.QuantityInCents)
        return nil, status.Errorf(codes.Internal, "cannot save auction: %s", err)
    }

    slog.InfoContext(ctx, "bid created", "auction_id", req.AuctionId, "user_it", req.UserId, "qty", req.QuantityInCents)
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
            slog.WarnContext(ctx, fmt.Sprintf("stream closed: %s", stream.Context().Err()), "auction_id", req.AuctionId)
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
