package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

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
		case "esc":
			m.currentView = "menu"
			return m, nil
		}
	case tea.WindowSizeMsg:
		h, v := lipgloss.NewStyle().GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		m.collectionsList.SetSize(msg.Width-h, msg.Height-v)
		m.recordsList.SetSize(msg.Width-h, msg.Height-v)
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case collectionsLoadedMsg:
		m.collections = msg.collections
		m.currentView = "collections"
		items := []list.Item{}
		for _, col := range msg.collections {
			items = append(items, item{title: col.Name, desc: col.Type})
		}
		m.collectionsList.SetItems(items)
		return m, nil
	case recordsLoadedMsg:
		m.records = msg.records
		m.currentView = "records"
		items := []list.Item{}
		for _, rec := range msg.records {
			items = append(items, item{title: rec.Id, desc: "Record"})
		}
		m.recordsList.SetItems(items)
		return m, nil
	case logsLoadedMsg:
		m.logs = msg.logs
		m.currentView = "logs"
		return m, nil
	case settingsLoadedMsg:
		m.settings = msg.settings
		m.currentView = "settings"
		return m, nil
	case backupCreatedMsg:
		m.currentView = "backups_done"
		return m, nil
	case errorMsg:
		m.err = msg.err
		return m, nil
	}

	if m.currentView == "menu" {
		m.list, cmd = m.list.Update(msg)
	} else if m.currentView == "collections" {
		var listCmd tea.Cmd
		m.collectionsList, listCmd = m.collectionsList.Update(msg)
		cmd = listCmd
	} else if m.currentView == "records" {
		var listCmd tea.Cmd
		m.recordsList, listCmd = m.recordsList.Update(msg)
		cmd = listCmd
	}



func (m model) viewSettings() string {
	if m.settings == nil {
		return "Loading settings..."
	}

	s := "Settings:\n\n"
	s += fmt.Sprintf("App Name: %s\n", m.settings.Meta.AppName)
	s += fmt.Sprintf("App URL: %s\n", m.settings.Meta.AppURL)
	s += fmt.Sprintf("Hide Controls: %t\n", m.settings.Meta.HideControls)
	s += "\nPress Esc to go back."
	return s
}

func (m model) viewLogs() string {
	if len(m.logs) == 0 {
		return "No logs found. Press Esc to go back."
	}

	s := "Logs:\n\n"
	for i, log := range m.logs {
		s += fmt.Sprintf("%d. Level %d: %s\n", i+1, log.Level, log.Message)
	}
	s += "\nPress Esc to go back."
	return s
}

// Messages
type collectionsLoadedMsg struct {
	collections []*core.Collection
}

type recordsLoadedMsg struct {
	records []*core.Record
}

type logsLoadedMsg struct {
	logs []*core.Log
}

type settingsLoadedMsg struct {
	settings *core.Settings
}

type backupCreatedMsg struct{}

type errorMsg struct {
	err error
}

// Commands
func (m model) loadCollections() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		collections := []*core.Collection{}
		err := m.app.CollectionQuery().All(&collections)
		if err != nil {
			// If table doesn't exist, suggest initializing the database
			if strings.Contains(err.Error(), "no such table") {
				return errorMsg{err: errors.New("Database not initialized. Please run PocketBase server with migrations first to set up the database.")}
			}
			return errorMsg{err: err}
		}

		return collectionsLoadedMsg{collections: collections}
	})
}

func (m model) loadRecords(collectionName string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		records := []*core.Record{}
		err := m.app.RecordQuery(collectionName).All(&records)
		if err != nil {
			return errorMsg{err: err}
		}

		return recordsLoadedMsg{records: records}
	})
}

func (m model) loadLogs() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		logs := []*core.Log{}
		err := m.app.LogQuery().All(&logs)
		if err != nil {
			return errorMsg{err: err}
		}

		return logsLoadedMsg{logs: logs}
	})
}

func (m model) loadSettings() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		settings := m.app.Settings()
		return settingsLoadedMsg{settings: settings}
	})
}

func (m model) createBackup() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		err := m.app.CreateBackup(context.Background(), "backup.zip")
		if err != nil {
			return errorMsg{err: err}
		}

		return backupCreatedMsg{}
	})
}

func main() {
	// Check if running in TTY
	if !isTTY() {
		fmt.Println("This application requires a TTY. Please run in a terminal.")
		os.Exit(1)
	}

	m := initialModel()

	// Start PocketBase server in background
	go func() {
		if err := m.app.Bootstrap(); err != nil {
			log.Printf("Bootstrap error: %v", err)
			return
		}
		if err := m.app.Start(); err != nil {
			log.Printf("PocketBase server error: %v", err)
		}
	}()

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

// isTTY checks if the program is running in a TTY environment
func isTTY() bool {
	// Check both stdout and stdin for TTY
	return term.IsTerminal(int(os.Stdout.Fd())) || term.IsTerminal(int(os.Stdin.Fd()))
}
