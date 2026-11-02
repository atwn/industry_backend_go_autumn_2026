package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Task struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	UpdatedAt time.Time `json:"updatedAt"`
}
type TaskRepo interface {
	Create(title string) (Task, error)
	Get(id string) (Task, bool)
	List(done bool) []Task
	SetDone(id string, done bool) (Task, error)
}
type Clock interface{ Now() time.Time }

var (
	ErrNotFound     = errors.New("task not found")
	ErrInvalidTitle = errors.New("invalid title")
)

type inMemoryTaskRepo struct {
	mu    sync.RWMutex
	clock Clock
	seq   uint64
	tasks map[string]Task
}

func NewInMemoryTaskRepo(clock Clock) TaskRepo {
	return &inMemoryTaskRepo{
		mu:    sync.RWMutex{},
		clock: clock,
		seq:   0,
		tasks: make(map[string]Task)}
}

func (r *inMemoryTaskRepo) Create(title string) (Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	title = strings.TrimSpace(title)
	if len(title) == 0 {
		return Task{}, ErrInvalidTitle
	}

	id := fmt.Sprintf("%v", r.seq)
	r.tasks[id] = Task{ID: id, Title: title, Done: false, UpdatedAt: r.clock.Now()}
	r.seq++
	return r.tasks[id], nil
}

func (r *inMemoryTaskRepo) Get(id string) (Task, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	v, ok := r.tasks[id]
	return v, ok
}

func (r *inMemoryTaskRepo) List(done bool) []Task {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]Task, 0)
	for _, task := range r.tasks {
		if task.Done == done {
			list = append(list, task)
		}
	}

	return list
}

func (r *inMemoryTaskRepo) SetDone(id string, done bool) (Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	task, ok := r.tasks[id]
	if !ok {
		return Task{}, ErrNotFound
	}

	if task.Done != done {
		task.Done = done
		task.UpdatedAt = r.clock.Now()
		r.tasks[id] = task
	}

	return task, nil
}

type httpHandler struct{ repo TaskRepo }

func NewHTTPHandler(repo TaskRepo) http.Handler {
	panic("TODO: implement")
}
func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("TODO: implement")
}
func (h *httpHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	panic("TODO: implement")
}
func (h *httpHandler) handleGet(w http.ResponseWriter, r *http.Request, id string) {
	panic("TODO: implement")
}
func (h *httpHandler) handleList(w http.ResponseWriter, r *http.Request) {
	panic("TODO: implement")
}
func (h *httpHandler) handlePatch(w http.ResponseWriter, r *http.Request, id string) {
	panic("TODO: implement")
}
func decodeStrictJSON(r io.Reader, v any) error {
	panic("TODO: implement")
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	panic("TODO: implement")
}
