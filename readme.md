# 📝 TUI Notes

A minimalist, high-efficiency terminal-based note-taking application. Designed for users who live in the terminal and want a keyboard-centric workflow for managing thoughts, snippets, and documentation.

Built with **Go** and the **Bubble Tea** (Charm.sh) framework.

---

## ✨ Features

- **Vim-like Navigation**: Use `j` and `k` for fluid list browsing.
- **Integrated Search**: Instantly narrow down notes with a real-time `/` filter.
- **Direct Management**: Rename, delete, and create notes without leaving the interface.
- **Save & Sync**: Quick `Ctrl+s` saving with local storage.
- **Visual Clarity**: Modern, color-coded UI with a dedicated help bar for ease of use.

---

## ⌨️ Shortcuts & Controls

### Main Navigation (List View)

| Key          | Action                   |
| :----------- | :----------------------- |
| `j` / `↓`    | Move selection down      |
| `k` / `↑`    | Move selection up        |
| `g` / `Home` | Go to top of list        |
| `G` / `End`  | Go to bottom of list     |
| `/`          | **Filter** notes by name |
| `Ctrl + r`   | **Rename** selected note |
| `Ctrl + d`   | **Delete** selected note |
| `?`          | Toggle Help menu         |

### File Operations

| Key        | Action                  |
| :--------- | :---------------------- |
| `Ctrl + n` | Create a **New** note   |
| `Ctrl + l` | Return to **List** view |
| `Ctrl + s` | **Save** current note   |
| `Esc`      | Go **Back** / Cancel    |
| `Ctrl + q` | **Quit** application    |

---

## 🚀 Installation

### Prerequisites

- [Go](https://golang.org/doc/install) (1.18 or higher)

### Build from Source

1. Clone the repository:
   ```bash
   git clone [https://github.com/DevSatyamCollab/tui-notes.git](https://github.com/DevSatyamCollab/tui-notes.git)
   cd tui-notes
   ```
2. Build the binary

   ```bash
   go build -o notes main.go
   ```

3. (Optional) Move to your path:
   ```bash
   sudo mv notes /usr/local/bin/
   ```
