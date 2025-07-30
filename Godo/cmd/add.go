package cmd

import (
	"fmt"
	"godo/internal/task"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

func isValidTime(ttime string) bool {
	if len(ttime) != 4 {
		return false
	}
	for _, c := range ttime {
		if c < '0' || c > '9' {
			return false
		}
	}
	hours := ttime[0:2]
	mins := ttime[2:4]

	h, _ := strconv.Atoi(hours)
	m, _ := strconv.Atoi(mins)

	return h >= 0 && h <= 23 && m >= 0 && m <= 59
}
func isValidDate(date string) bool {
	// Check format: DD/MM/YYYY
	parts := strings.Split(date, "/")
	if len(parts) != 3 {
		return false
	}

	dayStr, monthStr, yearStr := parts[0], parts[1], parts[2]

	// Must be 2/2/4 digits
	if len(dayStr) != 2 || len(monthStr) != 2 || len(yearStr) != 4 {
		return false
	}

	// Must be all digits
	for _, s := range []string{dayStr, monthStr, yearStr} {
		for _, c := range s {
			if c < '0' || c > '9' {
				return false
			}
		}
	}

	// Parse integers
	day, _ := strconv.Atoi(dayStr)
	month, _ := strconv.Atoi(monthStr)
	year, _ := strconv.Atoi(yearStr)

	// Check range
	if year < 1 || year > 9999 {
		return false
	}
	if month < 1 || month > 12 {
		return false
	}

	// Days in each month (index 1 = Jan)
	daysInMonth := []int{0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

	// Leap year check
	if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
		daysInMonth[2] = 29
	}

	if day < 1 || day > daysInMonth[month] {
		return false
	}

	return true
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Use to add a new task",
	Run: func(cmd *cobra.Command, args []string) {
		store, err := task.NewStore()
		if err != nil {
			fmt.Println("Failed to load task storage")
		}
		var name, date, ttime, status string

		for {
			keymap := huh.NewDefaultKeyMap()
			keymap.Quit = key.NewBinding(
				key.WithKeys("q", "esc", "ctrl+c"),
				key.WithHelp("q/esc/ctrl+c", "quit"),
			)

			form := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Name the Task").
						Placeholder(": ").
						Value(&name),
					huh.NewInput().
						Title("Enter deadline Date: (format: DD/MM/YYYY)").
						Value(&date),
					huh.NewInput().
						Title("Enter Deadline Time: (format: HHMM 24-hrs)").
						Value(&ttime),
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

			err = form.Run()
			if err != nil {
				fmt.Println("Error running form")
				return
			}

			if name == "" || ttime == "" || date == "" {
				fmt.Println("All fields are required")
				continue
			}
			if !isValidTime(ttime) {
				fmt.Println("Invalid time format")
				continue
			}
			if !isValidDate(date) {
				fmt.Println("Invalid Date format")
				continue
			}
			break
		}

		if status == "" {
			status = "ongoing"
		}

		t := task.Task{
			ID:           store.GenerateID(),
			Name:         name,
			CreatedAt:    time.Now().Format("02/01/2006 1504"),
			DeadlineDate: date,
			DeadlineTime: ttime,
			Status:       status,
		}
		store.Tasks = append(store.Tasks, t)
		err = store.Save()
		if err != nil {
			fmt.Println("Error saving file")
			return
		}

		fmt.Printf("Task %v: %s Saved Sucessfully\n", t.ID, t.Name)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
