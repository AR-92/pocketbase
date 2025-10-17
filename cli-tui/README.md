# PocketBase CLI TUI

A terminal user interface for managing PocketBase instances using Bubble Tea.

## Features

- **Collections Management**: View, create, and modify database collections
- **Records CRUD**: Browse, create, update, and delete records
- **Settings**: Configure application settings, mail, and storage
- **Backups**: Create, restore, and manage backups
- **Logs**: View application logs and statistics
- **Interactive Interface**: Full keyboard navigation and TUI experience

## Usage

```bash
cd cli-tui
go run main.go
```

## Controls

- `↑/↓`: Navigate menus
- `Enter`: Select item
- `Esc`: Go back
- `Ctrl+C` or `q`: Quit

## Requirements

- Terminal with TTY support
- Go 1.25.1+
- PocketBase data directory (set via `POCKETBASE_DATA_DIR` env var or defaults to `./pb_data`)

## Architecture

Built with:
- **Bubble Tea**: Terminal app framework
- **Bubbles**: UI components
- **Lipgloss**: Styling
- **PocketBase**: Backend library

The TUI provides a complete interface for PocketBase management without requiring a web browser.