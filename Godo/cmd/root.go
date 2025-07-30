package cmd

import (
	"fmt"
	"godo/internal/task"
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var purge bool
var purgeall bool

var rootCmd = &cobra.Command{
	Use:   "godo",
	Short: "A CLI task manager with TUI",
	Long:  "A CLI tool built in Go which allows you to manage your tasks with a beautiful terminal-UI",
	Run: func(cmd *cobra.Command, args []string) {
		store, err := task.NewStore()
		if err != nil {
			fmt.Println("Error accessing store, Exiting...")
			return
		}

		if purgeall {
			store.Tasks = nil
			err = store.Save()
			if err != nil {
				fmt.Println("Error saving after purgeall")
			}
			fmt.Println("All tasks deleted.")
			return
		}

		if purge {
			filtered := []task.Task{}
			for _, t := range store.Tasks {
				if t.Status != "completed" {
					filtered = append(filtered, t)
				}
			}
			store.Tasks = filtered
			err = store.Save()
			if err != nil {
				fmt.Println("Error saving after purge")
			}
			fmt.Println("Completed tasks deleted.")
			return
		}
		var choice string
		for {
			keymap := huh.NewDefaultKeyMap()
			keymap.Quit = key.NewBinding(key.WithKeys("q"))
			keymap.Quit = key.NewBinding(key.WithKeys("esc"))
			keymap.Quit = key.NewBinding(key.WithKeys("ctrl+c"))
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Title("Choose an Action (press ctrl+C to exit!)").
						Options(
							huh.NewOption("Add", "add"),
							huh.NewOption("Delete", "delete"),
							huh.NewOption("List", "list"),
							huh.NewOption("Edit", "edit"),
							huh.NewOption("Quit", "quit"),
						).
						Value(&choice),
				),
			).WithKeyMap(keymap).WithTheme(huh.ThemeCharm())

			err := form.Run()
			if err != nil {
				fmt.Println("Exiting...")
				return
			}
			switch choice {
			case "add":
				addCmd.Run(cmd, args)
			case "delete":
				delCmd.Run(cmd, args)
			case "list":
				listCmd.Run(cmd, args)
			case "edit":
				editCmd.Run(cmd, args)
			case "quit":
				fmt.Println("Exiting!...")
				return
			}
		}

	},
}

func Execute() {
	//Flags need to be registered before execute....shouldve been obvious🤦
	rootCmd.PersistentFlags().BoolVar(&purge, "purge", false, "Delete completed tasks")
	rootCmd.PersistentFlags().BoolVar(&purgeall, "purgeall", false, "Delete all tasks")
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}

}
