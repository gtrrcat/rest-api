package handlers

import (
	"encoding/json"
	"net/http"
	"rest-api/internal/database"
	"rest-api/internal/models"
	"strconv"
	"strings"
)

type Handlers struct {
	store *database.TaskStore
}

func NewHandlers(store *database.TaskStore) *Handlers {
	return &Handlers{store: store}
}
func responseWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)

}
func responseWithError(w http.ResponseWriter, statusCode int, message string) {
	responseWithJSON(w, statusCode, map[string]string{"error": message})
}
func (h *Handlers) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.store.GetAll()
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Ошибка при получении задач: "+err.Error())
		return
	}
	responseWithJSON(w, http.StatusOK, tasks)
}

func (h *Handlers) GetTaskByID(w http.ResponseWriter, r *http.Request) {

	// Получаем ID из URL
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")
	idStr := pathParts[0]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Неверный ID задачи"+err.Error())
		return
	}
	task, err := h.store.GetByID(id)
	if err != nil {
		responseWithError(w, http.StatusNotFound, "Задача не найдена"+err.Error())
		return
	}
	responseWithJSON(w, http.StatusOK, task)
}

func (h *Handlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	var input models.CreateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responseWithError(w, http.StatusBadRequest, "Неверный формат данных: "+err.Error())
		return
	}
	if strings.TrimSpace(input.Title) == "" {
		responseWithError(w, http.StatusBadRequest, "Поле title не может быть пустым")
		return
	}
	task, err := h.store.Create(input)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Ошибка при создании задачи: "+err.Error())
		return
	}
	responseWithJSON(w, http.StatusCreated, task)

}
func (h *Handlers) UpdateTask(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из URL
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")
	idStr := pathParts[0]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Неверный ID задачи: "+err.Error())
		return
	}

	var input models.UpdateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responseWithError(w, http.StatusBadRequest, "Неверный формат данных: "+err.Error())
		return
	}
	task, err := h.store.Update(id, input)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Ошибка при обновлении задачи: "+err.Error())
		return
	}
	responseWithJSON(w, http.StatusOK, task)
}

func (h *Handlers) DeleteTask(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из URL
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")
	idStr := pathParts[0]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Неверный ID задачи: "+err.Error())
		return
	}
	err = h.store.Delete(id)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Ошибка при удалении задачи: "+err.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, map[string]string{"message": "Задача успешно удалена"})
}
