package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// model
type model struct {
	height int
	width  int
	header string
	footer string
	style  *StyleBundle
}

func InitialModel() model {
	return model{
		header: "Welcome to note-takingApp 📓",
		footer: "Ctrl+n: New file . Ctrl+l: list . Ctrl+s: save . Esc: back/save . Ctrl+q: Quit",
		style:  DefaultStyles(),
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {

	h := m.style.headerStyle.Render(m.header)
	f := m.style.footerStyle.Render(m.footer)
	return fmt.Sprintf("\n%s\n\n%s\n", h, f)
}

func (m model) Init() tea.Cmd {

	return nil
}
