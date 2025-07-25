package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var tasks []Task

func loadTasks() {
	file, err := os.ReadFile("tasks.json")
	if err != nil {
		// If file does not exist, start with empty tasks
		if os.IsNotExist(err) {
			tasks = []Task{}
			return
		}
		fmt.Println("Error while reading the task json file:", err)
		tasks = []Task{}
		return
	}
	if err := json.Unmarshal(file, &tasks); err != nil {
		fmt.Println("Error parsing tasks.json:", err)
		tasks = []Task{}
	}
}

func saveTasks() {
	data, _ := json.MarshalIndent(tasks, "", " ")
	os.WriteFile("tasks.json", data, 0644)
}

func addTask(title string) {
	id := len(tasks) + 1
	task := Task{ID: id, Title: title, Done: false}
	tasks = append(tasks, task)
	saveTasks()
	fmt.Println("Added Task: ", title)
}

func resetTasks() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Are you sure you want to delete all tasks? (y/n): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "y" || input == "yes" {
		tasks = []Task{}
		saveTasks()
		fmt.Println("All tasks have been deleted and reset the list.")
	} else {
		fmt.Println("Reset cancelled.")
	}
}

func main() {

	loadTasks()
	fmt.Println("Welcome to Todo CLI (type 'help' for commands, 'exit' to quit)")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		if input == "exit" {
			fmt.Println("Goodbye!")
			break
		}
		if input == "help" {
			fmt.Println("Commands: \nadd <title>	(add new task) \nlist 		(list all tasks) \ndone <id> 	(mark as done) \nuncheck <id> 	(mark as not done) \ndelete <id> 	(delete the task) \nreset 		(clear the list) \nexit")
			continue
		}
		args := strings.Fields(input)
		if len(args) == 0 {
			continue
		}
		command := args[0]
		switch command {

		case "add":
			if len(args) < 2 {
				fmt.Println("Please provide a task title.")
				continue
			}
			title := strings.Join(args[1:], " ")
			addTask(title)
		case "list":
			if len(tasks) == 0 {
				fmt.Println("No tasks found.")
				continue
			}
			fmt.Println("Task List:")
			for _, task := range tasks {
				status := "[ ]"
				if task.Done {
					status = "[x]"
				}
				fmt.Printf("%d. %s %s\n", task.ID, status, task.Title)
			}
		case "reset":
			resetTasks()
			continue
		case "done":
			if len(args) < 2 {
				fmt.Println("Please provide the task ID to mark as done.")
				continue
			}
			id := args[1]
			found := false
			for i, task := range tasks {
				if fmt.Sprintf("%d", task.ID) == id {
					if task.Done {
						fmt.Printf("Task %s is already marked as done.\n", id)
					} else {
						tasks[i].Done = true
						saveTasks()
						fmt.Printf("Marked Task %s as done: %s\n", id, task.Title)
					}
					found = true
					break
				}
			}
			if !found {
				fmt.Println("Task ID not found:", id)
			}
		case "uncheck":
			if len(args) < 2 {
				fmt.Println("Please provide the task ID to uncheck")
				continue
			}
			id := args[1]
			found := false
			for i, task := range tasks {
				if fmt.Sprintf("%d", task.ID) == id {
					if !task.Done {
						fmt.Printf("Task %s is not marked as done.\n", id)
					} else {
						tasks[i].Done = false
						saveTasks()
						fmt.Printf("Change status of Task %s as not done: %s\n", id, task.Title)
					}
					found = true
					break
				}
			}
			if !found {
				fmt.Println("Task ID not found:", id)
			}
		case "delete":
			if len(args) < 2 {
				fmt.Println("Please provide the task ID to delete.")
				continue
			}
			id := args[1]
			found := false
			for i, task := range tasks {
				if fmt.Sprintf("%d", task.ID) == id {
					reader := bufio.NewReader(os.Stdin)
					fmt.Printf("Are you sure you want to delete Task %s: %s? (y/n): ", id, task.Title)
					input, _ := reader.ReadString('\n')
					input = strings.TrimSpace(strings.ToLower(input))
					if input == "y" || input == "yes" {
						tasks = append(tasks[:i], tasks[i+1:]...)
						// Reassign IDs to keep them unique and sequential
						for j := range tasks {
							tasks[j].ID = j + 1
						}
						saveTasks()
						fmt.Printf("Deleted Task %s: %s\n", id, task.Title)
					} else {
						fmt.Println("Delete cancelled.")
					}
					found = true
					break
				}
			}
			if !found {
				fmt.Println("Task ID not found:", id)
			}
		default:
			fmt.Println("Unknown command:", command)
		}
	}

}
