package task

import (
	"fmt"
)

// var ErrTitleRequired = errors.New("title is required")

type ValidationError struct {
	Field   string
	Message string
}

func (err ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", err.Field, err.Message)
}

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateTask(title, description string) (Task, error) {
	if title == "" {
		return Task{},
			fmt.Errorf(
				"can`t create task: %w",
				ValidationError{Field: "title", Message: "can`t be empty"},
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

func (s *Service) DeleteTask(id int64) bool {
	return s.repo.Delete(id)
}
