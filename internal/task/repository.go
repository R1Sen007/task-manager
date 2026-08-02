package task

type Repository struct {
	tasks  map[int64]Task
	nextID int64
}

func NewRepository() *Repository {
	return &Repository{
		tasks:  make(map[int64]Task),
		nextID: 1,
	}
}

var _ TaskRepository = (*Repository)(nil)

func (r *Repository) Create(task Task) Task {
	task.ID = r.nextID
	r.nextID++

	r.tasks[task.ID] = task

	return task
}

func (r *Repository) GetByID(id int64) (Task, bool) {
	task, exists := r.tasks[id]
	return task, exists
}

func (r *Repository) GetAll() []Task {
	result := make([]Task, 0, len(r.tasks))

	for _, value := range r.tasks {
		result = append(result, value)
	}

	return result
}

func (r *Repository) Delete(id int64) bool {
	if _, exists := r.tasks[id]; !exists {
		return false
	}

	delete(r.tasks, id)
	return true
}

func (r *Repository) Update(updatedTask Task) bool {
	if _, exists := r.tasks[updatedTask.ID]; !exists {
		return false
	}
	r.tasks[updatedTask.ID] = updatedTask
	return true
}
