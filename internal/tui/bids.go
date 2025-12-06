package tui

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

var (
    t = table.New(
        table.WithColumns(colums),
        table.WithHeight(5),
    )

    colums = []table.Column{
        {Title: "User", Width: 36},
        {Title: "Quantity", Width: 20},
    }

    defaultRow = table.Row{"...", "..."}
)

type bidList [4]bid

func (l bidList) Init() tea.Cmd {
    t.SetRows(l.toRows())
    return nil
}

func (l bidList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    if msg, ok := msg.(bid); ok {
        l.insertHighestBid(msg)
    }
    return l, nil
} 

func (l bidList) View() string {
    return boxStyle.Render(t.View())
}

func (l bidList) toRows() []table.Row {
    rows := make([]table.Row, 0, len(l))
    for _, v := range(l) {
        if v.UserId == "" {
            rows = append(rows, defaultRow)
            continue
        }
        rows = append(rows, v.toRow())
    }
    return rows
}

func (l *bidList) insertHighestBid(b bid) {
    // shift highest bids one place
    copy(l[1:], l[:len(l)-1])
    l[0] = b
    t.SetRows(l.toRows())
}
