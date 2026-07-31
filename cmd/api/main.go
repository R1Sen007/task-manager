package main

import (
	"errors"
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

	for _, taskInput := range inputs {
		tsk, err := service.CreateTask(taskInput.Title, taskInput.Description)

		var validationErr task.ValidationError
		if errors.As(err, &validationErr) {
			fmt.Printf(
				"validation error in field %q: %s\n",
				validationErr.Field,
				validationErr.Message,
			)
			continue
		}
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(tsk.Title, tsk.Description)
	}
}
