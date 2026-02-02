package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	vaultDir    string
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	docStyle    = lipgloss.NewStyle().Margin(1, 2)
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting home directory: %v", err)
	}

	vaultDir = filepath.Join(homeDir, ".note-takingApp")
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

// model
type model struct {
	newFileInput           textinput.Model
	noteTextArea           textarea.Model
	createFileInputVisible bool
	style                  *StyleBundle
	currentFile            *os.File
	list                   list.Model
	showinglist            bool
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

	// list
	notelist := listFiles()

	finalList := list.New(notelist, list.NewDefaultDelegate(), 0, 0)
	finalList.Title = "All notes 📋"
	finalList.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("16")).
		Background(lipgloss.Color("254")).
		Padding(0, 1)

	finalList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("ctrl+d"),
				key.WithHelp("ctrl+d", "delete file"),
			),
		}
	}

	return model{
		style:                  DefaultStyles(),
		newFileInput:           ti,
		noteTextArea:           ta,
		list:                   finalList,
		createFileInputVisible: false,
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-5)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+q":
			return m, tea.Quit
		case "ctrl+d":
			if m.showinglist {
				item, ok := m.list.SelectedItem().(item)
				if ok {
					filePath := filepath.Join(vaultDir, item.title)

					// remove the file
					if err := os.Remove(filePath); err != nil {
						log.Printf("Error deleting file: %v", err)
					}

					m.list.SetItems(listFiles())
				}
			}
			return m, nil
		case "esc":
			// list and filter state  to main UI
			if m.showinglist {
				if m.list.FilterState() == list.Filtering {
					break
				}

				m.showinglist = false
			}

			// create a new file to main UI
			if m.createFileInputVisible {
				m.createFileInputVisible = false
			}

			// textarea to main UI
			if m.currentFile != nil {
				m.noteTextArea.SetValue("")
				m.currentFile = nil
			}

			return m, nil
		case "ctrl+n":
			m.createFileInputVisible = true
			return m, nil
		case "enter":
			if m.currentFile != nil {
				break
			}

			if m.showinglist {
				item, ok := m.list.SelectedItem().(item)
				if ok {
					filePath := filepath.Join(vaultDir, item.title)

					content, err := os.ReadFile(filePath)
					if err != nil {
						log.Printf("Error reading the file: %v", err)
						return m, nil
					}

					m.noteTextArea.SetValue(string(content))

					// open the file
					f, err := os.OpenFile(filePath, os.O_RDWR, 0644)
					if err != nil {
						log.Printf("Error reading the file: %v", err)
						return m, nil
					}

					m.currentFile = f
					m.showinglist = false

					return m, nil
				}
			}

			fname := m.newFileInput.Value()
			// todo: create a file
			if fname != "" {
				fileName := fmt.Sprintf("%s.md", fname)
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
			notelist := listFiles()
			m.list.SetItems(notelist)
			m.showinglist = true
		}
	}

	if m.createFileInputVisible {
		m.newFileInput, cmd = m.newFileInput.Update(msg)
	}

	if m.currentFile != nil {
		m.noteTextArea, cmd = m.noteTextArea.Update(msg)
	}

	if m.showinglist {
		m.list, cmd = m.list.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {

	header := m.style.headerStyle.Render("Welcome to note-takingApp 📓")
	footer := m.style.footerStyle.Render(fmt.Sprintf(
		"Ctrl+n: New file . Ctrl+l: list . Ctrl+s: save . Esc: back . Ctrl+q: Quit"))

	view := ""
	if m.createFileInputVisible {
		view = m.newFileInput.View()
	}

	if m.currentFile != nil {
		view = m.noteTextArea.View()
	}

	if m.showinglist {
		view = m.list.View()
	}

	return fmt.Sprintf("\n%s\n\n%s\n\n%s\n", header, view, footer)
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func listFiles() []list.Item {
	// 1. Load the specific timezone
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		// Fallback to local if Kolkata isn't found
		loc = time.Local
	}

	items := make([]list.Item, 0)

	entries, err := os.ReadDir(vaultDir)
	if err != nil {
		log.Fatal("Error reading notes")
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		modeTime := info.ModTime().In(loc).Format("2006-01-02 03:04 PM")

		items = append(items, item{
			title: info.Name(),
			desc:  fmt.Sprintf("Modified: %s", modeTime),
		})

	}

	return items
}
