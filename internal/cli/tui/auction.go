package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	bidpb "github.com/sverdejot/grpc-streams/internal/api/grpc/bid/v1"
)

type errType struct{}

type auction struct {
	currentHighestBid highestBid
	sortedHighestBids bidList
	biddingButtons    bidButtons
    currentUserID   string
	f                 func() (any, error)
}

type creator func(string, int) error

func CreateAuction(f func() (any, error), fc creator, currentUser string) (a *auction) {
	a = &auction{
		f:                 f,
		sortedHighestBids: NewSortedBidList(),
		biddingButtons:    NewButtons(fc, currentUser),
        currentUserID: currentUser,
	}

	return a
}

func (a *auction) Init() tea.Cmd {
	return tea.Batch(a.checkBids(), a.currentHighestBid.Init(), a.sortedHighestBids.Init(), textinput.Blink)
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
	var cmds []tea.Cmd
	cmds = append(cmds, a.checkBids())
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return a, tea.Quit
		case tea.KeyTab, tea.KeyEnter:
			model, cmd := a.biddingButtons.Update(msg)
			cmds = append(cmds, cmd)
			a.biddingButtons = model.(bidButtons)
			return a, tea.Batch(cmds...)
		}
	case errType:
		return a, tea.Quit
	case *bidpb.Bid:
		model, cmd := a.sortedHighestBids.Update(Bid(a.currentHighestBid))
		cmds = append(cmds, cmd)
		a.sortedHighestBids = model.(bidList)

		model, cmd = a.currentHighestBid.Update(msg)
		cmds = append(cmds, cmd)
		a.currentHighestBid = model.(highestBid)

		model, cmd = a.biddingButtons.Update(a.currentHighestBid)
		cmds = append(cmds, cmd)
		a.biddingButtons = model.(bidButtons)

		return a, tea.Batch(cmds...)
	}
	// default cause fallthrough is not allowed in type switches
	return a, tea.Batch(cmds...)
}

func (a *auction) View() string {
	return lipgloss.Place(80, 23,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, a.currentHighestBid.View(), a.sortedHighestBids.View(), a.biddingButtons.View()),
		lipgloss.WithWhitespaceChars("競売"),
		lipgloss.WithWhitespaceForeground(subtle),
	)
}
