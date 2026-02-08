package store

import (
	"errors"
	"sync"

	"github.com/Ablikon/go-1st/internal/models"
)

var ErrNotFound = errors.New("task not found")

type Store struct {
	mu     sync.Mutex
	nextID int
	tasks  map[int]models.Task
}

func New() *Store {
	return &Store{
		nextID: 1,
		tasks:  make(map[int]models.Task),
	}
}

func (s *Store) Create(title string) models.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := models.Task{
		ID:    s.nextID,
		Title: title,
		Done:  false,
	}
	s.tasks[task.ID] = task
	s.nextID++
	return task
}

func (s *Store) Get(id int) (models.Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[id]
	return task, ok
}

func (s *Store) List(doneFilter *bool) []models.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]models.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		if doneFilter != nil && task.Done != *doneFilter {
			continue
		}
		result = append(result, task)
	}
	return result
}

func (s *Store) UpdateDone(id int, done bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[id]
	if !ok {
		return ErrNotFound
	}
	task.Done = done
	s.tasks[id] = task
	return nil
}
