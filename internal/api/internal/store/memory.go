package store

import (
	"maps"
	"slices"
	"sync"

	"github.com/sverdejot/grpc-streams/internal/api/internal"
)

type memAuctionStore struct {
    auctions map[string]internal.Auction
    mu sync.Mutex
}

func NewMemoryStore() *memAuctionStore {
    return &memAuctionStore{
        auctions: make(map[string]internal.Auction),
    }
}

func (s *memAuctionStore) Upsert(auction internal.Auction) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.auctions[auction.ID] = auction
    return nil
}

func (s *memAuctionStore) Get(ID string) (internal.Auction, bool) {
    s.mu.Lock()
    defer s.mu.Unlock()
    v, ok := s.auctions[ID]
    return v, ok
}

func (s *memAuctionStore) GetAll() []internal.Auction {
    return slices.Collect(maps.Values(s.auctions))
}

