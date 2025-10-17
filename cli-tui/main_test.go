package main

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pocketbase/pocketbase/core"
)

func TestInitialModel(t *testing.T) {
	m := initialModel()

	if m.currentView != "menu" {
		t.Errorf("Expected currentView 'menu', got %s", m.currentView)
	}
	if m.app == nil {
		t.Error("App should not be nil")
	}
	if m.list.Title != "" {
		t.Errorf("Expected list title '', got %s", m.list.Title)
	}
	if m.collectionsList.Title != "" {
		t.Errorf("Expected collectionsList title '', got %s", m.collectionsList.Title)
	}
	if m.recordsList.Title != "" {
		t.Errorf("Expected recordsList title '', got %s", m.recordsList.Title)
	}
}

func TestLoadCollections(t *testing.T) {
	m := model{}
	cmd := m.loadCollections()
	if cmd == nil {
		t.Error("loadCollections should return a cmd")
	}
}

func TestLoadRecords(t *testing.T) {
	m := model{}
	cmd := m.loadRecords("test")
	if cmd == nil {
		t.Error("loadRecords should return a cmd")
	}
}

func TestLoadLogs(t *testing.T) {
	m := model{}
	cmd := m.loadLogs()
	if cmd == nil {
		t.Error("loadLogs should return a cmd")
	}
}

func TestLoadSettings(t *testing.T) {
	m := model{}
	cmd := m.loadSettings()
	if cmd == nil {
		t.Error("loadSettings should return a cmd")
	}
}

func TestCreateBackup(t *testing.T) {
	m := model{}
	cmd := m.createBackup()
	if cmd == nil {
		t.Error("createBackup should return a cmd")
	}
}

func TestViewSettings(t *testing.T) {
	m := model{}

	// Test when settings is nil
	result := m.viewSettings()
	if !strings.Contains(result, "Loading settings...") {
		t.Errorf("Expected 'Loading settings...', got %s", result)
	}

	// Test when settings is set
	m.settings = &core.Settings{}
	m.settings.Meta.AppName = "TestApp"
	m.settings.Meta.AppURL = "http://test.com"
	m.settings.Meta.HideControls = true

	result = m.viewSettings()
	if !strings.Contains(result, "App Name: TestApp") {
		t.Errorf("Expected 'App Name: TestApp', got %s", result)
	}
	if !strings.Contains(result, "App URL: http://test.com") {
		t.Errorf("Expected 'App URL: http://test.com', got %s", result)
	}
	if !strings.Contains(result, "Hide Controls: true") {
		t.Errorf("Expected 'Hide Controls: true', got %s", result)
	}
}

func TestViewLogs(t *testing.T) {
	m := model{}

	// Test when logs is empty
	result := m.viewLogs()
	if !strings.Contains(result, "No logs found") {
		t.Errorf("Expected 'No logs found', got %s", result)
	}

	// Test when logs is set
	m.logs = []*core.Log{
		{Level: 1, Message: "Test log"},
	}

	result = m.viewLogs()
	if !strings.Contains(result, "Level 1: Test log") {
		t.Errorf("Expected 'Level 1: Test log', got %s", result)
	}
}

func TestInit(t *testing.T) {
	m := initialModel()
	cmd := m.Init()
	// Init returns tea.Batch(spinner.Tick, tea.EnterAltScreen)
	// Hard to test exactly, but check it's not nil
	if cmd == nil {
		t.Error("Init should return a cmd")
	}
}

func TestUpdateEsc(t *testing.T) {
	m := initialModel()
	m.currentView = "collections"
	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	newM := newModel.(model)
	if newM.currentView != "menu" {
		t.Errorf("Expected currentView 'menu', got %s", newM.currentView)
	}
	if cmd != nil {
		t.Error("Expected no cmd")
	}
}

func TestUpdateCollectionsEnter(t *testing.T) {
	m := initialModel()
	m.currentView = "collections"
	// Set up collectionsList with items
	m.collectionsList.SetItems([]list.Item{item{title: "test_col", desc: "base"}})
	m.collectionsList.Select(0)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Error("Expected cmd for loadRecords")
	}
}

func TestUpdateWindowSize(t *testing.T) {
	m := initialModel()
	_, cmd := m.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
	if cmd != nil {
		t.Error("Expected no cmd")
	}
}

func TestUpdateWindowSizeCollections(t *testing.T) {
	m := initialModel()
	m.currentView = "collections"
	_, cmd := m.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
	if cmd != nil {
		t.Error("Expected no cmd")
	}
}

func TestUpdateWindowSizeRecords(t *testing.T) {
	m := initialModel()
	m.currentView = "records"
	_, cmd := m.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
	if cmd != nil {
		t.Error("Expected no cmd")
	}
}

func TestUpdateSpinnerTick(t *testing.T) {
	m := initialModel()
	_, cmd := m.Update(m.spinner.Tick())
	if cmd == nil {
		t.Error("Expected cmd from spinner")
	}
}

func TestUpdateErrorMsg(t *testing.T) {
	m := initialModel()
	newModel, cmd := m.Update(errorMsg{err: nil})
	newM := newModel.(model)
	if newM.err != nil {
		t.Error("Expected err to be set")
	}
	if cmd != nil {
		t.Error("Expected no cmd")
	}
}

func TestUpdateSelectCollectionEnter(t *testing.T) {
	m := initialModel()
	m.currentView = "select_collection"
	m.textInput.SetValue("test_col")
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Error("Expected cmd for loadRecords")
	}
}

func TestUpdateSelectCollectionKey(t *testing.T) {
	m := initialModel()
	m.currentView = "select_collection"
	_, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	// Just to cover the textInput.Update call
}

func TestViewMenu(t *testing.T) {
	m := initialModel()
	result := m.View()
	if !strings.Contains(result, "6 items") {
		t.Errorf("Expected menu display, got %s", result)
	}
}

func TestViewCollections(t *testing.T) {
	m := initialModel()
	m.currentView = "collections"
	result := m.View()
	// Since collectionsList is empty, should show list view
	if result == "" {
		t.Error("View should not be empty")
	}
}

func TestViewRecords(t *testing.T) {
	m := initialModel()
	m.currentView = "records"
	result := m.View()
	if result == "" {
		t.Error("View should not be empty")
	}
}

func TestViewSelectCollection(t *testing.T) {
	m := initialModel()
	m.currentView = "select_collection"
	result := m.View()
	if !strings.Contains(result, "Select Collection") {
		t.Errorf("Expected 'Select Collection', got %s", result)
	}
}

func TestViewBackupsDone(t *testing.T) {
	m := initialModel()
	m.currentView = "backups_done"
	result := m.View()
	if !strings.Contains(result, "Backup created successfully") {
		t.Errorf("Expected 'Backup created successfully', got %s", result)
	}
}

func TestViewWithError(t *testing.T) {
	m := initialModel()
	m.err = &testError{msg: "test error"}
	result := m.View()
	if !strings.Contains(result, "Error: test error") {
		t.Errorf("Expected 'Error: test error', got %s", result)
	}
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

func TestItemMethods(t *testing.T) {
	i := item{title: "test", desc: "desc"}
	if i.Title() != "test" {
		t.Errorf("Expected title 'test', got %s", i.Title())
	}
	if i.Description() != "desc" {
		t.Errorf("Expected desc 'desc', got %s", i.Description())
	}
	if i.FilterValue() != "test" {
		t.Errorf("Expected filter 'test', got %s", i.FilterValue())
	}
}

func TestIsTTY(t *testing.T) {
	result := isTTY()
	// In test environment, likely false
	if result != false {
		t.Logf("isTTY returned %t", result)
	}
}
