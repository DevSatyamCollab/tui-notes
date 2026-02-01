package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	vaultDir    string
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting home directory: %v", err)
	}

	vaultDir = filepath.Join(homeDir, ".note-takingApp")
}

// model
type model struct {
	height                 int
	width                  int
	newFileInput           textinput.Model
	noteTextArea           textarea.Model
	createFileInputVisible bool
	style                  *StyleBundle
	currentFile            *os.File
}

func InitialModel() model {
	err := os.MkdirAll(vaultDir, 0750)
	if err != nil {
		log.Fatalf("can't create a dir: %v", err)
	}

	// textinput
	ti := textinput.New()
	ti.Placeholder = "What would you like to call it?"
	ti.Focus()
	ti.Width = 50
	ti.CharLimit = 155
	ti.Cursor.Style = cursorStyle
	ti.PromptStyle = cursorStyle
	ti.TextStyle = cursorStyle

	// textarea
	ta := textarea.New()
	ta.Placeholder = "Write your note here..."
	ta.Focus()

	return model{
		style:                  DefaultStyles(),
		newFileInput:           ti,
		noteTextArea:           ta,
		createFileInputVisible: false,
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+q":
			return m, tea.Quit
		case "ctrl+n":
			m.createFileInputVisible = true
			return m, nil
		case "enter":

			// todo: create a file
			fileName := fmt.Sprintf("%s.md", m.newFileInput.Value())
			if fileName != "" {
				filePath := filepath.Join(vaultDir, fileName)

				// if not exists, create a file
				if _, err := os.Stat(filePath); err == nil {
					return m, nil
				}

				f, err := os.Create(filePath)
				if err != nil {
					log.Fatalf("Can't create a file: %v", err)
				}

				m.currentFile = f
				m.createFileInputVisible = false
				m.newFileInput.SetValue("")
			}

		case "ctrl+s":
			// textarea value -> write it in the field descriptor and close it
			if m.currentFile == nil {
				break
			}

			err := m.currentFile.Truncate(0)
			if err != nil {
				fmt.Println("can not save the file :(")
				return m, nil
			}

			_, err = m.currentFile.Seek(0, 0)
			if err != nil {
				fmt.Println("can not save the file :(")
			}

			// writing a file
			_, err = m.currentFile.WriteString(m.noteTextArea.Value())
			if err != nil {
				fmt.Println("can not save the file :(")
			}

			err = m.currentFile.Close()
			if err != nil {
				fmt.Println("can not close the file :(")
			}

			// update the state
			m.currentFile = nil
			m.noteTextArea.SetValue("")

			return m, nil
		case "ctrl+l":

		case "Esc":
		}
	}

	if m.createFileInputVisible {
		m.newFileInput, cmd = m.newFileInput.Update(msg)
	}

	if m.currentFile != nil {
		m.noteTextArea, cmd = m.noteTextArea.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {

	h := m.style.headerStyle.Render("Welcome to note-takingApp 📓")
	f := m.style.footerStyle.Render(fmt.Sprintf(
		"Ctrl+n: New file . Ctrl+l: list . Ctrl+s: save . Esc: back/save . Ctrl+q: Quit"))

	v := ""
	if m.createFileInputVisible {
		v = m.newFileInput.View()
	}

	if m.currentFile != nil {
		v = m.noteTextArea.View()
	}

	return fmt.Sprintf("\n%s\n\n%s\n\n%s\n", h, v, f)
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}
