package main

import "github.com/charmbracelet/lipgloss"

type StyleBundle struct {
	headerStyle lipgloss.Style
	footerStyle lipgloss.Style
}

func DefaultStyles() *StyleBundle {
	return &StyleBundle{
		headerStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("16")).
			Background(lipgloss.Color("205")).
			Padding(1, 2),

		footerStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("69")),
	}
}
