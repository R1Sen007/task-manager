package task

import (
	"errors"
	"testing"
)

type fakeRepository struct {
	createCalled bool
	createdTask  Task

	getByIDTask   Task
	getByIDExists bool

	updateCalled bool
	updatedTask  Task
	updateResult bool
}

func (r *fakeRepository) Create(task Task) Task {
	r.createCalled = true
	r.createdTask = task

	task.ID = 1
	return task
}

func (r *fakeRepository) GetByID(id int64) (Task, bool) {
	return r.getByIDTask, r.getByIDExists
}

func (r *fakeRepository) Update(task Task) bool {
	r.updateCalled = true
	r.updatedTask = task
	return r.updateResult
}

func (r *fakeRepository) Delete(id int64) bool {
	return false
}

var _ TaskRepository = (*fakeRepository)(nil)

func TestService_CreateTask(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		description string
		wantError   bool
		wantCreate  bool
	}{
		{
			name:        "success",
			title:       "Learn Go",
			description: "Understand tests",
			wantError:   false,
			wantCreate:  true,
		},
		{
			name:        "empty title",
			title:       "",
			description: "Invalid task",
			wantError:   true,
			wantCreate:  false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			fakeRepo := &fakeRepository{}
			service := NewService(fakeRepo)

			createdTask, err := service.CreateTask(testCase.title, testCase.description)

			gotError := err != nil

			if gotError != testCase.wantError {
				t.Fatalf(
					"expected error=%v, got error=%v: %v",
					testCase.wantError,
					gotError,
					err,
				)
			}

			if fakeRepo.createCalled != testCase.wantCreate {
				t.Fatalf(
					"expected repository Create called=%v, got %v",
					testCase.wantCreate,
					fakeRepo.createCalled,
				)
			}

			if gotError {
				var validationErr ValidationError

				if !errors.As(err, &validationErr) {
					t.Fatalf("expected ValidationError, got %T: %v", err, err)
				}

				if validationErr.Field != "title" {
					t.Fatalf(
						"expected validation field %q, got %q",
						"title",
						validationErr.Field,
					)
				}

				if validationErr.Message != "can't be empty" {
					t.Fatalf(
						"expected validation message %q, got %q",
						"can't be empty",
						validationErr.Message,
					)
				}
				return
			}

			if fakeRepo.createdTask.Title != testCase.title {
				t.Fatalf(
					"expected repository task title %q, got %q",
					testCase.title,
					fakeRepo.createdTask.Title,
				)
			}

			if fakeRepo.createdTask.Description != testCase.description {
				t.Fatalf(
					"expected repository task description %q, got %q",
					testCase.description,
					fakeRepo.createdTask.Description,
				)
			}

			if fakeRepo.createdTask.Done {
				t.Fatal("expected repository task Done to be false")
			}

			if fakeRepo.createdTask.ID != 0 {
				t.Fatalf(
					"service must not assign ID, got %d",
					fakeRepo.createdTask.ID,
				)
			}

			if createdTask.ID != 1 {
				t.Fatalf(
					"expected created task ID 1, got %d",
					createdTask.ID,
				)
			}

		})
	}
}

func TestService_UpdateTask(t *testing.T) {
	newTitle := ""
	doneFalse := false

	tests := []struct {
		name           string
		id             int64
		input          UpdateTaskInput
		repoTask       Task
		repoTaskExists bool
		updateResult   bool

		wantError       bool
		wantNotFound    bool
		wantValidation  bool
		wantUpdateCall  bool
		wantUpdatedTask Task
	}{
		{
			name:           "task not found",
			id:             10,
			input:          UpdateTaskInput{},
			repoTaskExists: false,
			wantError:      true,
			wantNotFound:   true,
			wantUpdateCall: false,
		},
		{
			name: "empty title",
			id:   10,
			input: UpdateTaskInput{
				Title: &newTitle,
			},
			repoTask: Task{
				ID:          10,
				Title:       "old title",
				Description: "old description",
				Done:        false,
			},
			repoTaskExists: true,
			wantError:      true,
			wantValidation: true,
			wantUpdateCall: false,
		},
		{
			name: "update done true to false",
			id:   10,
			input: UpdateTaskInput{
				Done: &doneFalse,
			},
			repoTask: Task{
				ID:          10,
				Title:       "old title",
				Description: "old description",
				Done:        true,
			},
			repoTaskExists: true,
			updateResult:   true,
			wantError:      false,
			wantUpdateCall: true,
			wantUpdatedTask: Task{
				ID:          10,
				Title:       "old title",
				Description: "old description",
				Done:        false,
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			fakeRepo := &fakeRepository{
				getByIDTask:   testCase.repoTask,
				getByIDExists: testCase.repoTaskExists,
				updateResult:  testCase.updateResult,
			}

			service := NewService(fakeRepo)

			updatedTask, err := service.UpdateTask(testCase.id, testCase.input)

			gotError := err != nil
			if gotError != testCase.wantError {
				t.Fatalf(
					"expected error=%v, got error=%v: %v",
					testCase.wantError,
					gotError,
					err,
				)
			}

			if fakeRepo.updateCalled != testCase.wantUpdateCall {
				t.Fatalf(
					"expected Update called=%v, got %v",
					testCase.wantUpdateCall,
					fakeRepo.updateCalled,
				)
			}

			if testCase.wantNotFound {
				if !errors.Is(err, ErrTaskNotFound) {
					t.Fatalf(
						"expected ErrTaskNotFound, got %T: %v",
						err,
						err,
					)
				}

				return
			}

			if testCase.wantValidation {
				var validationErr ValidationError

				if !errors.As(err, &validationErr) {
					t.Fatalf(
						"expected ValidationError, got %T: %v",
						err,
						err,
					)
				}

				if validationErr.Field != "title" {
					t.Fatalf(
						"expected validation field %q, got %q",
						"title",
						validationErr.Field,
					)
				}

				return
			}

			if fakeRepo.updatedTask != testCase.wantUpdatedTask {
				t.Fatalf(
					"expected repository updated task %+v, got %+v",
					testCase.wantUpdatedTask,
					fakeRepo.updatedTask,
				)
			}

			if updatedTask != testCase.wantUpdatedTask {
				t.Fatalf(
					"expected returned task %+v, got %+v",
					testCase.wantUpdatedTask,
					updatedTask,
				)
			}
		})
	}
}
