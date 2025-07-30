package cmd

import (
	"fmt"
	"godo/internal/task"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var delCmd = &cobra.Command{
	Use:   "del",
	Short: "Delete a task",
	Long:  "Select a task and hit enter to delete",
	Run: func(cmd *cobra.Command, args []string) {
		store, err := task.NewStore()
		if err != nil {
			fmt.Println("Failed to load tasks")
		}
		if len(store.Tasks) == 0 {
			fmt.Println("No Tasks to delete.")
			return
		}
		var options []huh.Option[string]
		//var taskIDs []int

		for _, t := range store.Tasks {
			label := fmt.Sprintf("%d: %s (Due date: %s %s)", t.ID, t.Name, t.DeadlineDate, t.DeadlineTime)
			options = append(options, huh.NewOption(label, t.Name))
			//taskIDs = append(taskIDs, t.ID)
		}

		// Custom keymap for select form
		keymap := huh.NewDefaultKeyMap()
		keymap.Quit = key.NewBinding(
			key.WithKeys("q", "esc", "ctrl+c"),
			key.WithHelp("q/esc/ctrl+c", "quit"),
		)

		var selectedname string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Select a Task to DELETE").
					Options(options...).
					Value(&selectedname),
			),
		).WithKeyMap(keymap).WithTheme(huh.ThemeCharm())

		err = form.Run()
		if err != nil {
			fmt.Println("Exiting...")
			return
		}

		var confirm bool
		confirmForm := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Confirm Task deletion").
					Value(&confirm).
					Affirmative("Yes!").
					Negative("nah, I'm good!"),
			),
		).WithKeyMap(keymap).WithTheme(huh.ThemeCharm())

		err = confirmForm.Run()
		if err != nil || !confirm {
			fmt.Println("Cancelled, Exiting...")
			return
		}
		var newTasks []task.Task
		var deletedID int

		for _, t := range store.Tasks {
			if t.Name == selectedname {
				deletedID = t.ID
			} else {
				newTasks = append(newTasks, t)
			}
		}
		store.Tasks = newTasks
		err = store.Save()
		if err != nil {
			fmt.Println("Failed to save file")
			return
		}

		fmt.Printf("Task %d: %s deleted successfully\n", deletedID, selectedname)
	},
}

func init() {
	rootCmd.AddCommand(delCmd)
}
