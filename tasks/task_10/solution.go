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
	title = strings.TrimSpace(title)
	if len(title) == 0 {
		return Task{}, ErrInvalidTitle
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.seq++
	id := fmt.Sprintf("%v", r.seq) // works only if repository is a singleton
	task := Task{
		ID:        id,
		Title:     title,
		Done:      false,
		UpdatedAt: r.clock.Now()}
	r.tasks[id] = task
	return task, nil
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

	// task is not found:
	task, ok := r.tasks[id]
	if !ok {
		return Task{}, ErrNotFound
	}

	// change is not required:
	if task.Done == done {
		return task, nil
	}

	// do the change:
	task.Done = done
	task.UpdatedAt = r.clock.Now()
	r.tasks[id] = task

	return task, nil
}

type httpHandler struct {
	mux  *http.ServeMux
	repo TaskRepo
}

func NewHTTPHandler(repo TaskRepo) http.Handler {
	mux := http.NewServeMux()
	h := httpHandler{mux: mux, repo: repo}

	mux.HandleFunc("POST /tasks", h.handleCreate)
	mux.HandleFunc("GET /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		h.handleGet(w, r, id)
	})
	mux.HandleFunc("GET /tasks", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		values, ok := query["done"]
		if ok && len(values) == 1 && (values[0] == "true" || values[0] == "false") { // ?done=true или ?done=false
			h.handleList(w, r)
			return
		}

		if !ok { // параметр 'done' не задан, обрабатывать как ?done=false
			//query.Set("done", "false")
			//r.URL.RawQuery = query.Encode()
			h.handleList(w, r)
			return
		}

		http.Error(w, "'done' must be 'true' or 'false' when set", http.StatusBadRequest) // 'done' задан неверно
	})
	mux.HandleFunc("PATCH /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		h.handlePatch(w, r, id)
	})
	// явно прописать обработчик для HEAD метода чтобы вернуть 405, а не обрабатывать его как GET
	// (см. https://pkg.go.dev/net/http#ServeMux, подсказал AI - сам бы не догадался)
	mux.HandleFunc("HEAD /tasks", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
	mux.HandleFunc("HEAD /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
	return &h
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
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
