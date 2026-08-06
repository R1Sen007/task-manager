package task

import (
	"errors"
	"testing"
)

type fakeRepository struct {
	createCalled bool
	createdTask  Task
}

func (r *fakeRepository) Create(task Task) Task {
	r.createCalled = true
	r.createdTask = task

	task.ID = 1
	return task
}

func (r *fakeRepository) GetByID(id int64) (Task, bool) {
	return Task{}, false
}

func (r *fakeRepository) Update(task Task) bool {
	return false
}

func (r *fakeRepository) Delete(id int64) bool {
	return false
}

var _ TaskRepository = (*fakeRepository)(nil)

func TestService_CreateTask_EmptyTitle(t *testing.T) {
	fakeRepo := &fakeRepository{}
	service := NewService(fakeRepo)

	_, err := service.CreateTask("", "descr")

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

	if fakeRepo.createCalled {
		t.Fatal("repository Create must not be called")
	}

}

func TestService_CreateTask_Success(t *testing.T) {
	fakeRepo := &fakeRepository{}
	service := NewService(fakeRepo)

	createdTask, err := service.CreateTask("smTitle", "descr")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !fakeRepo.createCalled {
		t.Fatal("expected repository Create to be called")
	}

	if fakeRepo.createdTask.Title != "smTitle" {
		t.Fatalf(
			"expected repository task title %q, got %q",
			"smTitle",
			fakeRepo.createdTask.Title,
		)
	}

	if fakeRepo.createdTask.Description != "descr" {
		t.Fatalf(
			"expected repository task description %q, got %q",
			"descr",
			fakeRepo.createdTask.Description,
		)
	}

	if fakeRepo.createdTask.Done {
		t.Fatal("expected repository task Done to be false")
	}

	if createdTask.ID != 1 {
		t.Fatalf("expected created task ID 1, got %d", createdTask.ID)
	}
}
