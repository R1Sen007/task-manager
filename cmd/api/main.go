package main

import (
	"fmt"

	"github.com/R1Sen007/task-manager/internal/task"
)

type CreateTaskInput struct {
	Title       string
	Description string
}

func main() {
	repo := task.NewRepository()
	service := task.NewService(repo)

	inputs := []CreateTaskInput{
		{Title: "Learn Go", Description: "Understand errors"},
		{Title: "", Description: "This should fail"},
		{Title: "Learn interfaces", Description: "Go deeper"},
	}

	for _, task_input := range inputs {
		tsk, err := service.CreateTask(task_input.Title, task_input.Description)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(tsk.Title, tsk.Description)
	}
}
