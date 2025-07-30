# 🧠 godo — Your Smart CLI Task Companion

**godo** is a sleek, terminal-based todo manager built with Go — combining the power of the CLI with a smooth TUI using Charmbracelet's `huh` for an interactive experience.

> ⚡ No clutter. No distractions. Just productivity.

---

## ✨ Features

- 📋 Add, delete, edit, and list tasks with ease
- ⏰ Set deadlines and track remaining time
- ✅ Mark tasks as completed
- 🧹 `--purge` to remove completed tasks
- 💣 `--purgeall` to wipe the slate clean
- 🎨 Beautiful terminal UI powered by `huh`
- 💾 Data stored locally in JSON format @`~/.godo/tasks.json`

---

## 🚀 Getting Started

```bash
git clone https://github.com/Amag1n3/Godo.git
cd Godo
go build -o godo
./godo
```


## ⚙️ Functionality
```
godo             # Launches the TUI
godo add         # Adds a new task via prompt
godo list        # Shows all tasks
godo edit        # Edits a selected task
godo delete      # Deletes a selected task
godo --purge     # Deletes only completed tasks
godo --purgeall  # Deletes ALL tasks
```

## 🔧 Tech Stack
## 🛠️ Tech Stack

- [Go](https://golang.org/)
- [Cobra](https://github.com/spf13/cobra)
- [tablewriter](https://github.com/olekukonko/tablewriter)
- [huh](https://github.com/charmbracelet/huh)
- [Bubbletea](https://github.com/charmbracelet/bubbletea)

