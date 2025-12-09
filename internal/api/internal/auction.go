package internal

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	bidpb "github.com/sverdejot/grpc-streams/internal/api/grpc/bid/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuctionStore interface {
    Upsert(Auction) error
    Get(ID string) (Auction, bool)
    GetAll() []Auction
}

type Auction struct {
    ID string
    Item string
    CurrentHighestBid Bid
    Bids []Bid // slice for the sake of simplicity
    EndsAt time.Time
}

type auctionOpt func(*Auction)

func NewAuction(opts ...auctionOpt) (Auction, error) {
    id, err := uuid.NewV7()
    if err != nil {
        return Auction{}, fmt.Errorf("error generating id for auction: %w", err)
    }
    a := Auction{
        ID: id.String(),
    }

    for _, fn := range opts {
        fn(&a)
    }
    return a, nil
}

func WithEndsAt(t time.Duration) auctionOpt {
    return func(a *Auction) {
        a.EndsAt = time.Now().Add(t)
    }
}

func Withitem(i string) auctionOpt {
    return func(a *Auction) {
        a.Item = i
    }
}

func (a *Auction) InsertBid(bid Bid) bool {
    if !a.isHigherThanCurrent(bid) || a.isFromSameUserAsCurrent(bid) {
        return false
    }

    a.Bids = append(a.Bids, a.CurrentHighestBid)
    a.CurrentHighestBid = bid
    return true
}

func (a Auction) isHigherThanCurrent(bid Bid) bool {
    return bid.QuantityInCents > a.CurrentHighestBid.QuantityInCents
}

func (a Auction) isFromSameUserAsCurrent(bid Bid) bool {
    return a.CurrentHighestBid.UserID == bid.UserID
}

func (a Auction) ToProto() *bidpb.Auction {
    bp := make([]*bidpb.Bid, 0, len(a.Bids))
    for _, b := range a.Bids {
        bp = append(bp, b.ToProto())
    }

    return &bidpb.Auction{
        AuctionId: a.ID,
        Item: a.Item,
        Bids: bp,
        EndsAt: timestamppb.New(a.EndsAt),
    }
}
