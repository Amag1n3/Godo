package cmd

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"godo/internal/task"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var sortBy string
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Run: func(cmd *cobra.Command, args []string) {
		err := runList()
		if err != nil {
			fmt.Println("Error printing table: ", err)
			return

		}

		//sorting task functionality is built-in into list as a flag "--sort" or "-s"

	},
}

func runList() error {
	store, err := task.NewStore()
	if err != nil {
		return err
	}

	if len(store.Tasks) == 0 {
		fmt.Println("No tasks found..")
		return nil
	}

	switch sortBy {
	case "id":
		sort.Slice(store.Tasks, func(i, j int) bool {
			return store.Tasks[i].ID < store.Tasks[j].ID
		})
	case "name":
		sort.Slice(store.Tasks, func(i, j int) bool {
			return store.Tasks[i].Name < store.Tasks[j].Name
		})
	case "created":
		sort.Slice(store.Tasks, func(i, j int) bool {
			return store.Tasks[i].CreatedAt < store.Tasks[j].CreatedAt
		})
	case "deadline":
		sort.Slice(store.Tasks, func(i, j int) bool {
			di := store.Tasks[i].DeadlineDate + store.Tasks[i].DeadlineTime
			dj := store.Tasks[j].DeadlineDate + store.Tasks[j].DeadlineTime
			return di < dj
		})
	case "status":
		sort.Slice(store.Tasks, func(i, j int) bool {
			return store.Tasks[i].Status < store.Tasks[j].Status
		})
	}

	var rows []table.Row

	for _, t := range store.Tasks {
		created, _ := time.Parse("02/01/2006 1504", t.CreatedAt)
		createdDate := created.Format("02/01/2006")

		deadlineStr := t.DeadlineDate + " " + t.DeadlineTime
		deadline, parseErr := time.Parse("02/01/2006 1504", deadlineStr)
		deadlineDisplay := deadlineStr
		if parseErr == nil {
			deadlineDisplay = deadline.Format("02/01/2006 1504")
		}

		var timeLeft string
		if parseErr != nil {
			timeLeft = "Invalid"
		} else {
			left := time.Until(deadline)
			if left < 0 {
				timeLeft = "Task Overdue"
			} else {
				leftDays := int(left.Hours()) / 24
				leftHours := int(left.Hours()) % 24
				leftMins := int(left.Minutes()) % 60
				timeLeft = fmt.Sprintf("%dd %dh %dm", leftDays, leftHours, leftMins)
			}
		}

		rows = append(rows, table.Row{
			strconv.Itoa(t.ID),
			t.Name,
			createdDate,
			deadlineDisplay,
			timeLeft,
		})
	}

	columns := []table.Column{
		{Title: "ID", Width: 6},
		{Title: "Name", Width: 25},
		{Title: "Created", Width: 12},
		{Title: "Deadline", Width: 18},
		{Title: "Time Left", Width: 15},
	}
	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(rows)+1), // +1 for header
	)

	// Apply basic styles
	s := table.DefaultStyles()
	s.Header = s.Header.Bold(true)
	tbl.SetStyles(s)

	// Run Bubble Tea app
	m := model{tbl: tbl}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		return err
	}

	return nil
}

// Bubble Tea model
type model struct {
	tbl table.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			fmt.Println("Exiting..")
			return m, tea.Quit
		}
	}

	m.tbl, cmd = m.tbl.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return "\n" + m.tbl.View() + "\nPress 'q' to quit\n"
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&sortBy, "sort", "s", "", "Sort by: id | name | created | deadline | status")

}
