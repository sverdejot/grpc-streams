package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	bidpb "github.com/sverdejot/grpc-streams/internal/api/grpc/bid/v1"
)

type bid struct {
    UserId string
    QtyInCents int32
}

type errType struct{}

type auction struct {
    bids []bid
    f func() (any, error)
}

func CreateAuction(f func() (any, error)) *auction {
    return &auction{f: f}
}

func (a *auction) Init() tea.Cmd {
    return a.checkBids()
}

func (a *auction) checkBids() tea.Cmd {
    return func() tea.Msg {
        v, err := a.f()
        if err != nil {
            return errType{}
        }
        return v
    }
}

func (a *auction) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {

    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            return a, tea.Quit
        }
    case errType:
        tea.Quit()
    case bidpb.Bid:
        b := bid{
            UserId: msg.UserId,
            QtyInCents: msg.QtyInCents,
        }
        a.bids = append(a.bids, b)
        return a, a.checkBids()
    }

    return a, nil
}

func (a *auction) View() string {
    return strings.Join(toString(a.bids), "\n")
}

func toString(bs []bid) []string {
    s := make([]string, 0, len(bs))
    for _, v := range bs {
        s = append(s, fmt.Sprintf("Bid: [%s, %d]", v.UserId, v.QtyInCents))
    }
    return s
}

