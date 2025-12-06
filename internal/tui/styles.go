package tui

import "github.com/charmbracelet/lipgloss"

var (
    boxStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("#874BFD")).
        Padding(1, 0).
        BorderTop(true).
        BorderLeft(true).
        BorderRight(true).
        BorderBottom(true)

    highestBidStyle = lipgloss.NewStyle().
        Align(lipgloss.Center).
        Width(60)

    subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
)
