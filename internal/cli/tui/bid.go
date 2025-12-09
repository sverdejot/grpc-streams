package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
)

var (
	defaultRow = table.Row{"...", "..."}
)

type Bid struct {
	UserId     string
	QtyInCents int
}

func (b Bid) toRow() []string {
	if b.UserId == "" {
		return defaultRow
	}
	return []string{b.UserId, b.FormatQty()}
}

func (b Bid) FormatQty() string {
	dllr := b.QtyInCents / 100
	cents := b.QtyInCents % 100

	if dllr < 1000 {
		return fmt.Sprintf("$%d,%02d", dllr, cents)
	}

	sdllr := strconv.Itoa(dllr)
	parts := make([]string, 0, (len(sdllr)%3)+1)
	for len(sdllr) > 3 {
		parts = append([]string{sdllr[len(sdllr)-3:]}, parts...)
		sdllr = sdllr[:len(sdllr)-3]
	}
	parts = append([]string{sdllr}, parts...)

	return "$" + strings.Join(parts, ".") + "," + strconv.Itoa(cents)
}
