package internal

import (
	bidpb "github.com/sverdejot/grpc-streams/internal/api/grpc/bid/v1"
)

type Bid struct {
	AuctionID, UserID string
	QuantityInCents   int
}

func (b Bid) ToProto() *bidpb.Bid {
	return &bidpb.Bid{
		UserId:     b.UserID,
		QtyInCents: int32(b.QuantityInCents),
	}
}

type bidBuilderOpts func(*Bid)

func NewBid(opts ...bidBuilderOpts) Bid {
    var b Bid
    for _, fn := range opts {
        fn(&b)
    }
    return b
}

func WithAuctionID(auctionID string) bidBuilderOpts{
    return func(b *Bid) {
        b.AuctionID = auctionID
    }
}

func WithUserID(userID string) bidBuilderOpts {
    return func(b *Bid) {
        b.UserID = userID
    }
}

func WithQuantity(qty int32) bidBuilderOpts {
    return func(b *Bid) {
        b.QuantityInCents = int(qty)
    }
}
