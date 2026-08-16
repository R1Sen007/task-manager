package task

import (
	"errors"
	"fmt"
)

var ErrTaskNotFound = errors.New("task not found")

type ValidationError struct {
	Field   string
	Message string
}

func (err ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", err.Field, err.Message)
}

type Service struct {
	repo TaskRepository
}

type TaskRepository interface {
	Create(task Task) Task
	GetByID(id int64) (Task, bool)
	Update(updatedTask Task) bool
	Delete(id int64) bool
}

type UpdateTaskInput struct {
	Title       *string
	Description *string
	Done        *bool
}

func NewService(repo TaskRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateTask(title, description string) (Task, error) {
	if title == "" {
		return Task{},
			fmt.Errorf(
				"can't create task: %w",
				ValidationError{Field: "title", Message: "can't be empty"},
			)
	}

	return s.repo.Create(
		Task{
			Title:       title,
			Description: description,
			Done:        false,
		},
	), nil
}

func (s *Service) GetTask(id int64) (Task, bool) {
	return s.repo.GetByID(id)
}

func (s *Service) UpdateTask(
	id int64,
	input UpdateTaskInput,
) (Task, error) {
	task, exists := s.repo.GetByID(id)
	if !exists {
		return Task{}, ErrTaskNotFound
	}

	if input.Title != nil {
		if *input.Title == "" {
			return Task{}, fmt.Errorf(
				"can't update task: %w",
				ValidationError{
					Field:   "title",
					Message: "can't be empty",
				},
			)
		}

		task.Title = *input.Title
	}

	if input.Description != nil {
		task.Description = *input.Description
	}

	if input.Done != nil {
		task.Done = *input.Done
	}

	if updated := s.repo.Update(task); !updated {
		return Task{}, ErrTaskNotFound
	}

	return task, nil
}

func (s *Service) DeleteTask(id int64) bool {
	return s.repo.Delete(id)
}
