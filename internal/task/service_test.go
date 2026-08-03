package task

import (
	"errors"
	"testing"
)

type fakeRepository struct {
	createCalled bool
}

func (r *fakeRepository) Create(task Task) Task {
	r.createCalled = true
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
