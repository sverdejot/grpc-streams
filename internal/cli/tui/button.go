package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var bidButtonPercentageMap map[int]float32 = map[int]float32{
	0: 0.10,
	1: 0.25,
	2: 0.50,
}

type refocusButtonMsg int

type bidButton struct {
	focused   bool
	pos       int
	bidToMake Bid
	fn        creator
    user    string
}

func (b bidButton) Init() tea.Cmd {
	return nil
}

func (b bidButton) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type != tea.KeyEnter {
			return b, nil
		}
		if err := b.fn(b.bidToMake.UserId, b.bidToMake.QtyInCents); err != nil {
			return b, tea.Quit
		}
		return b, nil
	case refocusButtonMsg:
		b.focused = int(msg) == b.pos
	case Bid:
		b.bidToMake = msg
		b.bidToMake.QtyInCents = int(float32(msg.QtyInCents) * (1.0 + bidButtonPercentageMap[b.pos]))
	}
	return b, nil
}

func (b bidButton) View() string {
	style := button
	if b.focused {
		style = focusedButton
	}
	return style.Render(fmt.Sprintf("%d%% (%s)", int(bidButtonPercentageMap[b.pos]*100), b.bidToMake.FormatQty()))
}

type bidButtons struct {
	bl      [3]bidButton
	focused int
}

func NewButtons(f creator, user string) bidButtons {
	bs := bidButtons{}
	for i := range bs.bl {
		bs.bl[i].pos = i
		bs.bl[i].fn = f
        bs.bl[i].user = user
	}
	bs.bl[0].focused = true
	return bs
}

func (bs bidButtons) Init() tea.Cmd {
	return nil
}

func (bs bidButtons) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			bs.focused = (bs.focused + 1) % len(bs.bl)
			var model tea.Model
			for i, b := range bs.bl {
				model, _ = b.Update(refocusButtonMsg(bs.focused))
				bs.bl[i] = model.(bidButton)
			}
		case tea.KeyEnter:
			var model tea.Model
			for i, b := range bs.bl {
				model, _ = b.Update(msg)
				bs.bl[i] = model.(bidButton)
			}
		}

	case Bid:
        fmt.Println("got a bid ina  button")
		var model tea.Model
		for i, b := range bs.bl {
			model, _ = b.Update(msg)
			bs.bl[i] = model.(bidButton)
		}
	}

	return bs, nil
}

func (bs bidButtons) View() string {
	bvw := make([]string, len(bs.bl))
	for i, b := range bs.bl {
		bvw[i] = b.View()
	}
	return lipgloss.JoinHorizontal(lipgloss.Center, bvw...)
}
