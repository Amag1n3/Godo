package cmd

import (
	"fmt"
	"godo/internal/task"
	"log"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit a task",
	Long:  "Edit any existing task with the option to edit as many values.",
	Run: func(cmd *cobra.Command, args []string) {
		store, err := task.NewStore()
		if err != nil {
			log.Fatal("Error loading tasks: ", err)
		}
		if len(store.Tasks) == 0 {
			fmt.Println("No Tasks to edit, Exiting...")
			return
		}

		var options []huh.Option[string]
		taskMap := make(map[string]task.Task)

		for _, t := range store.Tasks {
			label := fmt.Sprintf("%d: %s (Due: %s %s)", t.ID, t.Name, t.DeadlineDate, t.DeadlineTime)
			IDstring := strconv.Itoa(t.ID)
			options = append(options, huh.NewOption(label, IDstring))
			taskMap[IDstring] = t
		}
		keymap := huh.NewDefaultKeyMap()
		keymap.Quit = key.NewBinding(
			key.WithKeys("q", "esc", "ctrl+c"),
			key.WithHelp("q/esc/ctrl+c", "quit"),
		)

		var selectedID string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Select a task to edit").
					Options(options...).
					Value(&selectedID),
			),
		).WithKeyMap(keymap).WithTheme(huh.ThemeCharm())

		err = form.Run()
		if err != nil {
			fmt.Println("Error running form!...Exiting")
			return
		}

		selected := taskMap[selectedID]
		name := selected.Name
		deadlinedate := selected.DeadlineDate
		deadlinetime := selected.DeadlineTime
		status := selected.Status

		editform := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Name the Task").
					Placeholder(": ").
					Value(&name),
				huh.NewInput().
					Title("Enter deadline Date: (format: DD/MM/YYYY)").
					Value(&deadlinedate),
				huh.NewInput().
					Title("Enter Deadline Time: (format: HHMM 24-hrs)").
					Value(&deadlinetime),
				huh.NewSelect[string]().
					Title("Pick task status: ").
					Options(
						huh.NewOption("Ongoing", "ongoing"),
						huh.NewOption("Completed", "completed"),
						huh.NewOption("Paused", "paused"),
					).
					Value(&status),
			),
		).WithKeyMap(keymap).WithTheme(huh.ThemeCharm())

		err = editform.Run()
		if err != nil {
			fmt.Println("Error running form!..Exiting")
			return
		}

		for i := range store.Tasks {
			if store.Tasks[i].ID == selected.ID {
				store.Tasks[i] = task.Task{
					ID:           selected.ID,
					Name:         name,
					DeadlineDate: deadlinedate,
					DeadlineTime: deadlinetime,
					Status:       status,
					CreatedAt:    selected.CreatedAt,
				}
				break

			}
		}

		err = store.Save()
		if err != nil {
			fmt.Println("Error saving file, Exiting...")
			return
		}

		fmt.Printf("Task %s has been updated successfully!\n", name)
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
