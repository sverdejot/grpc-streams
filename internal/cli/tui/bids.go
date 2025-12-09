package tui

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	colums = []table.Column{
		{Title: "User", Width: 36},
		{Title: "Quantity", Width: 20},
	}
)

type bidList struct {
	l [4]Bid
	t *table.Model // ptr so only 8 bytes are copied into the stack every time bidList is passed by value
}

func NewSortedBidList() bidList {
	bl := bidList{}
	bl.t = new(table.Model)
	t := table.New(
		table.WithColumns(colums),
		table.WithHeight(5),
		table.WithRows(bl.toRows()),
	)
	bl.t = &t // aux variable bc table.New is not addressable
	return bl
}

func (bl bidList) Init() tea.Cmd {
	return nil
}

func (l bidList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(Bid); ok {
		l.insertHighestBid(msg)
	}
	return l, nil
}

func (bl bidList) View() string {
	return boxStyle.Render(bl.t.View())
}

func (bl bidList) toRows() []table.Row {
	rows := make([]table.Row, 0, len(bl.l))
	for _, v := range bl.l {
		rows = append(rows, v.toRow())
	}
	return rows
}

func (bl *bidList) insertHighestBid(b Bid) {
	// shift highest bids one place
	copy(bl.l[1:], bl.l[:len(bl.l)-1])
	bl.l[0] = b
	bl.t.SetRows(bl.toRows())
}
