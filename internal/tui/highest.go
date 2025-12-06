package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	bidpb "github.com/sverdejot/grpc-streams/internal/api/grpc/bid/v1"
)

type highestBid bid

func (h highestBid) Init() tea.Cmd {
    return nil
}
 
func (h highestBid) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case *bidpb.Bid:
        h.UserId = msg.GetUserId()
        h.QtyInCents = int(msg.GetQtyInCents())
    }
    return h, nil
}

func (h highestBid) View() string {
    return boxStyle.Render(highestBidStyle.Render(fmt.Sprintf("Highest [%s] bid is %s" , h.UserId, bid(h).FormatQty())))
}
