package task

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/AkashGit21/task-ms/lib/persistence"
	"github.com/AkashGit21/task-ms/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const DEFAULT_PAGE_SIZE = 12

type task struct {
	ID                  string    `json:"id"`
	Title               string    `json:"title"`
	Content             string    `json:"content"`
	HTMLStylizedContent string    `json:"stylized_content"`
	Status              string    `json:"status"`
	CreatedAt           time.Time `json:"created_at"`
	CreatedBy           string    `json:"created_by"`
	ModifiedAt          time.Time `json:"modified_at"`
	ModifiedBy          string    `json:"modified_by"`
}
type createTaskRequest struct {
	Title               *string `json:"title,omitempty"`
	Content             string  `json:"content"`
	HTMLStylizedContent *string `json:"stylized_content,omitempty"`
	Status              *string `json:"status,omitempty"`
}
type patchTaskRequest struct {
	Title               *string `json:"title,omitempty"`
	Content             *string `json:"content,omitempty"`
	HTMLStylizedContent *string `json:"stylized_content,omitempty"`
	Status              *string `json:"status,omitempty"`
}
type listTasksResponse struct {
	Tasks      []persistence.Task `json:"tasks"`
	NextCursor *string            `json:"next_cursor,omitempty"`
}

/** Creates a new Task **/
func (tvh *TaskHandler) createTask(w http.ResponseWriter, r *http.Request) {
	utils.DebugLog("Inside createTask")
	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Content == "" {
		http.Error(w, "Content is mandatory", http.StatusBadRequest)
		return
	}

	status, ok := persistence.TaskStatus_value[*req.Status]
	if !ok {
		// Set the deafult Task status as TODO
		status = persistence.STATUS_TODO
	}

	title, stylizedContent := "", ""
	if req.Title != nil {
		title = *req.Title
	}
	if req.HTMLStylizedContent != nil {
		stylizedContent = *req.HTMLStylizedContent
	}

	taskID := uuid.New().String()
	currentTime := time.Now().UTC()
	taskObject := persistence.Task{
		ID:                  taskID,
		Title:               title,
		Content:             req.Content,
		HTMLStylizedContent: stylizedContent,
		Status:              persistence.TaskStatus(status),

		Discarded:  false,
		CreatedAt:  currentTime,
		ModifiedAt: currentTime,
		// TODO: Add user id after basic authn is implemented
		CreatedBy:  "SYSTEM",
		ModifiedBy: "SYSTEM",
	}
	rowsAffected, err := tvh.TaskOps.SaveRecord(taskObject)
	if rowsAffected <= 0 || err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := task{
		ID:                  taskObject.ID,
		Title:               taskObject.Title,
		Content:             taskObject.Content,
		HTMLStylizedContent: taskObject.HTMLStylizedContent,
		Status:              taskObject.Status.String(),
		CreatedAt:           taskObject.CreatedAt,
		CreatedBy:           taskObject.CreatedBy,
		ModifiedAt:          taskObject.ModifiedAt,
		ModifiedBy:          taskObject.ModifiedBy,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/** Get the Task with ID **/
func (tvh *TaskHandler) getTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["taskID"]
	if id == "" {
		http.Error(w, "Task ID is mandatory", http.StatusBadRequest)
		return
	}

	taskRecord, err := tvh.TaskOps.GetRecord(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if taskRecord == nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	resp := task{
		ID:                  taskRecord.ID,
		Title:               taskRecord.Title,
		Content:             taskRecord.Content,
		HTMLStylizedContent: taskRecord.HTMLStylizedContent,
		Status:              taskRecord.Status.String(),
		CreatedAt:           taskRecord.CreatedAt,
		CreatedBy:           taskRecord.CreatedBy,
		ModifiedAt:          taskRecord.ModifiedAt,
		ModifiedBy:          taskRecord.ModifiedBy,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/** Patch certain fields of Task with given ID **/
func (tvh *TaskHandler) patchTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["taskID"]

	if id == "" {
		http.Error(w, "Task ID is mandatory", http.StatusBadRequest)
		return
	}

	taskRecord, err := tvh.TaskOps.GetRecord(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if taskRecord == nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	var req patchTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedRecord, err := patchMutableFields(taskRecord, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedTasks, err := tvh.TaskOps.UpdateRecord(id, updatedRecord)
	if updatedTasks < 1 || err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := task{
		ID:                  updatedRecord.ID,
		Title:               updatedRecord.Title,
		Content:             updatedRecord.Content,
		HTMLStylizedContent: updatedRecord.HTMLStylizedContent,
		Status:              updatedRecord.Status.String(),
		CreatedAt:           updatedRecord.CreatedAt,
		CreatedBy:           updatedRecord.CreatedBy,
		ModifiedAt:          updatedRecord.ModifiedAt,
		ModifiedBy:          updatedRecord.ModifiedBy,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

/** Perform soft delete of Task with given ID **/
func (tvh *TaskHandler) deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["taskID"]

	if id == "" {
		http.Error(w, "Task ID is mandatory", http.StatusBadRequest)
		return
	}

	err := tvh.TaskOps.DeactivateRecord(id)
	if err != nil {
		if strings.EqualFold(err.Error(), "no rows affected") {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

/** List the available Tasks with filters in paginated view **/
func (tvh *TaskHandler) listTasks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	statusStr := query.Get("status")
	cursor := query.Get("cursor")
	limitStr := query.Get("limit")

	var status *persistence.TaskStatus
	if statusStr != "" {
		statusVal, ok := persistence.TaskStatus_value[statusStr]
		if !ok {
			// Ignoring incorrect status - APIs should not break
			return
		}
		status = new(persistence.TaskStatus)
		*status = persistence.TaskStatus(statusVal)
	}

	limit := DEFAULT_PAGE_SIZE
	if limitStr != "" {
		limitInt, err := strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}
		limit = limitInt
	}

	params := persistence.TaskFilterParams{
		Status: status,
		Cursor: &cursor,
		Limit:  limit,
	}

	tasks, err := tvh.TaskOps.FetchRecords(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := listTasksResponse{Tasks: tasks}

	if len(tasks) == limit {
		lastTask := tasks[len(tasks)-1]

		cursorData := persistence.CursorData{
			ModifiedAt: lastTask.ModifiedAt,
			ID:         lastTask.ID,
		}

		cursorBytes, err := json.Marshal(cursorData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cursorStr := base64.StdEncoding.EncodeToString(cursorBytes)
		resp.NextCursor = &cursorStr
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func patchMutableFields(task *persistence.Task, req patchTaskRequest) (persistence.Task, error) {
	if req.Content != nil {
		task.Content = *req.Content
	}

	if req.HTMLStylizedContent != nil {
		task.HTMLStylizedContent = *req.HTMLStylizedContent
	}

	if req.Status != nil {
		if status, ok := persistence.TaskStatus_value[*req.Status]; !ok {
			return *task, errors.New("invalid status")
		} else {
			task.Status = persistence.TaskStatus(status)
		}
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	return *task, nil
}
