package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	pocketbase "github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	_ "github.com/pocketbase/pocketbase/migrations"
	"golang.org/x/term"
)

type model struct {
	app             *pocketbase.PocketBase
	currentView     string
	list            list.Model
	collectionsList list.Model
	recordsList     list.Model
	textInput       textinput.Model
	spinner         spinner.Model
	collections     []*core.Collection
	records         []*core.Record
	logs            []*core.Log
	settings        *core.Settings
	err             error
}

type collectionsLoadedMsg struct {
	collections []*core.Collection
}

type recordsLoadedMsg struct {
	records []*core.Record
}

type settingsLoadedMsg struct {
	settings *core.Settings
}

type logsLoadedMsg struct {
	logs []*core.Log
}

type errorMsg struct {
	err error
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func initialModel() model {
	// Initialize PocketBase
	dataDir := "./pb_data"
	if envDir := os.Getenv("POCKETBASE_DATA_DIR"); envDir != "" {
		dataDir = envDir
	}

	app := pocketbase.NewWithConfig(pocketbase.Config{
		DefaultDataDir: dataDir,
		DefaultDev:     false,
	})

	// Initialize Bubble Tea components
	items := []list.Item{
		item{title: "Collections", desc: "Manage database collections"},
		item{title: "Records", desc: "View and edit records"},
		item{title: "Settings", desc: "Application settings"},
		item{title: "Backups", desc: "Backup and restore"},
		item{title: "Logs", desc: "View application logs"},
		item{title: "Exit", desc: "Exit the application"},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = ""

	cl := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	cl.Title = ""

	rl := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	rl.Title = ""

	ti := textinput.New()
	ti.Placeholder = "Enter command..."
	ti.CharLimit = 100
	ti.Width = 50

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		app:             app,
		currentView:     "menu",
		list:            l,
		collectionsList: cl,
		recordsList:     rl,
		textInput:       ti,
		spinner:         s,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		tea.EnterAltScreen,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.currentView == "menu" {
				selectedItem := m.list.SelectedItem()
				if selectedItem != nil {
					switch selectedItem.(item).title {
					case "Collections":
						return m, m.loadCollections()
					case "Records":
						return m, m.loadCollections()
					case "Settings":
						return m, m.loadSettings()
					case "Backups":
						return m, m.createBackup()
					case "Logs":
						return m, m.loadLogs()
					case "Exit":
						return m, tea.Quit
					}
				}
			} else if m.currentView == "collections" {
				selectedItem := m.collectionsList.SelectedItem()
				if selectedItem != nil {
					collectionName := selectedItem.(item).title
					return m, m.loadRecords(collectionName)
				}
			} else if m.currentView == "select_collection" {
				collectionName := m.textInput.Value()
				if collectionName != "" {
					return m, m.loadRecords(collectionName)
				}
			}
			return m, nil
		case "esc":
			m.currentView = "menu"
			return m, nil
		default:
			return m, nil
		}
	case tea.WindowSizeMsg:
		h, v := lipgloss.NewStyle().GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		m.collectionsList.SetSize(msg.Width-h, msg.Height-v)
		m.recordsList.SetSize(msg.Width-h, msg.Height-v)
		return m, nil
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func (m model) View() string {
	var s string

	if m.currentView == "menu" {
		s = m.list.View()
	} else if m.currentView == "collections" {
		s = m.collectionsList.View()
	} else if m.currentView == "records" {
		s = m.recordsList.View()
	} else if m.currentView == "settings" {
		s = m.viewSettings()
	} else if m.currentView == "backups_done" {
		s = "Backup created successfully at backup.zip\n\nPress Esc to go back."
	} else if m.currentView == "logs" {
		s = m.viewLogs()
	}

	if m.err != nil {
		s += "\n\nError: " + m.err.Error()
	}

	return s
}

func (m model) loadCollections() tea.Cmd                  { return nil }
func (m model) loadSettings() tea.Cmd                     { return nil }
func (m model) createBackup() tea.Cmd                     { return nil }
func (m model) loadLogs() tea.Cmd                         { return nil }
func (m model) loadRecords(collectionName string) tea.Cmd { return nil }
func (m model) viewSettings() string                      { return "Settings not implemented" }
func (m model) viewLogs() string                          { return "Logs not implemented" }

func main() {
	p := tea.NewProgram(initialModel())
	_, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}

func isTTY() bool {
	// Check both stdout and stdin for TTY
	return term.IsTerminal(int(os.Stdout.Fd())) || term.IsTerminal(int(os.Stdin.Fd()))
}
