package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Ablikon/go-1st/internal/store"
)

type TaskHandler struct {
	Store *store.Store
}

type errorResponse struct {
	Error string `json:"error"`
}

type createTaskRequest struct {
	Title string `json:"title"`
}

type updateTaskRequest struct {
	Done *bool `json:"done"`
}

func (h *TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGet(w, r)
	case http.MethodPost:
		h.handlePost(w, r)
	case http.MethodPatch:
		h.handlePatch(w, r)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
	}
}

func (h *TaskHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	idParam := strings.TrimSpace(query.Get("id"))
	if idParam != "" {
		id, err := strconv.Atoi(idParam)
		if err != nil || id <= 0 {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid id"})
			return
		}
		task, ok := h.Store.Get(id)
		if !ok {
			writeJSON(w, http.StatusNotFound, errorResponse{Error: "task not found"})
			return
		}
		writeJSON(w, http.StatusOK, task)
		return
	}

	var doneFilter *bool
	if doneParam := strings.TrimSpace(query.Get("done")); doneParam != "" {
		parsed, err := strconv.ParseBool(doneParam)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid done"})
			return
		}
		doneFilter = &parsed
	}

	tasks := h.Store.List(doneFilter)
	writeJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid title"})
		return
	}
	title := strings.TrimSpace(req.Title)
	if title == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid title"})
		return
	}

	task := h.Store.Create(title)
	writeJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) handlePatch(w http.ResponseWriter, r *http.Request) {
	idParam := strings.TrimSpace(r.URL.Query().Get("id"))
	if idParam == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid id"})
		return
	}
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid id"})
		return
	}

	var req updateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid done"})
		return
	}
	if req.Done == nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid done"})
		return
	}

	if err := h.Store.UpdateDone(id, *req.Done); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, errorResponse{Error: "task not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"updated": true})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

