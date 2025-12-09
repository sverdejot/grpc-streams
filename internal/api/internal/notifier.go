package internal

import (
	"fmt"
	"sync"
)

type AuctionUpdatesNotificator struct {
    mus map[string]*sync.Mutex
    listeners map[string][]chan Bid
}

func NewNotificator() *AuctionUpdatesNotificator {
    return &AuctionUpdatesNotificator{
        mus: make(map[string]*sync.Mutex),
        listeners: make(map[string][]chan Bid),
    }
}

func (n *AuctionUpdatesNotificator) Notify(auction Auction) {
    mu, ok := n.mus[auction.ID]
    if !ok {
        mu = new(sync.Mutex)
        n.mus[auction.ID] = mu
    }
    mu.Lock()
    defer mu.Unlock()

    chs, ok := n.listeners[auction.ID]
    if !ok {
        return
    }

    for _, ch := range chs {
        ch <- auction.CurrentHighestBid
    }
}

func (n *AuctionUpdatesNotificator) Add(auction Auction) {
    n.mus[auction.ID] = new(sync.Mutex)
}

func (n *AuctionUpdatesNotificator) Subscribe(auctionID string) (chan Bid, error) {
    mu, ok := n.mus[auctionID]
    if !ok {
        return nil, fmt.Errorf("auction %s not notifying", auctionID)
    }
    mu.Lock()
    defer mu.Unlock()

    ch := make(chan Bid)

    _, ok = n.listeners[auctionID]
    if !ok {
        n.listeners[auctionID] = make([]chan Bid, 0)
    }
    n.listeners[auctionID] = append(n.listeners[auctionID], ch)

    return ch, nil
}
