package tui

import (
	"fmt"
	"strconv"
	"strings"
)

type bid struct {
    UserId string
    QtyInCents int
}

func (b bid) toRow() []string {
    return []string{b.UserId, b.FormatQty()}
}

func (b bid) FormatQty() string {
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
