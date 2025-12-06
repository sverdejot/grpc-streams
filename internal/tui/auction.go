package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	bidpb "github.com/sverdejot/grpc-streams/internal/api/grpc/bid/v1"
)


type errType struct{}

type auction struct {
    currentHighestBid highestBid
    sortedHighestBids bidList
    f func() (any, error)
}

func CreateAuction(f func() (any, error)) *auction {
    return &auction{f: f}
}

func (a *auction) Init() tea.Cmd {
    return tea.Batch(a.checkBids(), a.currentHighestBid.Init(), a.sortedHighestBids.Init())
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
    cmds := make([]tea.Cmd, 0, 2)
    cmds = append(cmds, a.checkBids())
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            return a, tea.Quit
        }
    case errType:
        return a, tea.Quit
    case *bidpb.Bid:
        model, cmd := a.sortedHighestBids.Update(bid(a.currentHighestBid))
        cmds = append(cmds, cmd)
        a.sortedHighestBids = model.(bidList)

        model, cmd = a.currentHighestBid.Update(msg)
        cmds = append(cmds, cmd)
        a.currentHighestBid = model.(highestBid)

        return a, tea.Batch(cmds...)
    }
    // default cause fallthrough is not allowed in type switches
    return a, nil
}

func (a *auction) View() string {
    return lipgloss.Place(80, 20,
        lipgloss.Center, lipgloss.Center,
        lipgloss.JoinVertical(lipgloss.Center, a.currentHighestBid.View(), a.sortedHighestBids.View()),
        lipgloss.WithWhitespaceChars("猫咪"),
        lipgloss.WithWhitespaceForeground(subtle),
	)
}

