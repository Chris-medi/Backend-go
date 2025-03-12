package handlers

import (
	"backend/FakeBackend"

	"backend/types"
	"backend/verification"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type TestId struct {
	name       string
	id         string
	wantStatus int
	wantMsg    string
}

type TestInput struct {
	name       string
	input      string
	wantStatus int
	wantMsg    string
}

func TestHealthCheck(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test
	if assert.NoError(t, HealthCheck(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var response map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "healthy", response["status"])
	}
}

func TestGetAllTask(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test
	if assert.NoError(t, GetAllTask(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var response types.TasksResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, types.MessageTasksRetrieved, response.Message)
	}
}

func TestGetTaskByIdOk(t *testing.T) {
	e := echo.New()
	task := FakeBackend.CreateTask(types.Task{Title: "Test Task"})

	req := httptest.NewRequest(http.MethodGet, "/task/"+task.Id.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/task/:id")
	c.SetParamNames("id")
	c.SetParamValues(task.Id.String())

	if assert.NoError(t, GetTaskById(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Task retrieved successfully")
	}
}

func TestGetTaskById(t *testing.T) {
	tests := []TestId{
		{
			name:       "Empty ID",
			id:         "",
			wantStatus: http.StatusBadRequest,
			wantMsg:    types.MessageIDRequired,
		},
		{
			name:       "Invalid UUID",
			id:         "invalid-uuid",
			wantStatus: http.StatusBadRequest,
			wantMsg:    types.MessageInvalidIDFormat,
		},
		{
			name:       "Valid UUID but Task Not Found",
			id:         uuid.New().String(),
			wantStatus: http.StatusNotFound,
			wantMsg:    types.MessageTaskNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)

			err := GetTaskById(c)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)

			var response types.TaskResponse
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantMsg, response.Message)
		})
	}
}

func TestCreateTask(t *testing.T) {
	tests := []TestInput{
		{
			name:       "Invalid JSON",
			input:      `{"title": }`,
			wantStatus: http.StatusBadRequest,
			wantMsg:    "code=400, message=Syntax error",
		},
		{
			name:       "Empty Title",
			input:      `{"title": ""}`,
			wantStatus: http.StatusBadRequest,
			wantMsg:    "Key: 'Task.Title' Error:Field validation for 'Title' failed on the 'required' tag",
		},
		{
			name:       "Valid Task",
			input:      `{"title": "Test Task"}`,
			wantStatus: http.StatusCreated,
			wantMsg:    types.MessageTaskCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = &verification.TaskValidator{Validator: validator.New()}

			req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(tt.input))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := CreateTask(c)
			if tt.wantStatus == http.StatusCreated {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantStatus, rec.Code)
				var response types.TaskResponse
				err = json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantMsg, response.Message)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantMsg)
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {

	task := FakeBackend.CreateTask(types.Task{Title: "Test Task"})

	tests := []TestId{
		{
			name:       "Empty ID",
			id:         "",
			wantStatus: http.StatusBadRequest,
			wantMsg:    types.MessageIDRequired,
		},
		{
			name:       "Invalid UUID",
			id:         "invalid-uuid",
			wantStatus: http.StatusBadRequest,
			wantMsg:    types.MessageInvalidIDFormat,
		},
		{
			name:       "Task Not Found",
			id:         uuid.New().String(),
			wantStatus: http.StatusNotFound,
			wantMsg:    types.MessageTaskNotFound,
		},
		{
			name:       "Task Deleted successfully",
			id:         task.Id.String(),
			wantStatus: http.StatusOK,
			wantMsg:    types.MessageTaskDeleted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)

			err := DeleteTask(c)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)

			var response types.TaskResponse
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantMsg, response.Message)
		})
	}
}

func TestEditTask(t *testing.T) {
	tests := []TestInput{
		{
			name:       "Invalid JSON",
			input:      `{"title": }`,
			wantStatus: http.StatusBadRequest,
			wantMsg:    "code=400, message=Syntax error",
		},
		{
			name:       "Empty Title",
			input:      `{"title": "", "id": "` + uuid.New().String() + `"}`,
			wantStatus: http.StatusBadRequest,
			wantMsg:    "Key: 'Task.Title' Error:Field validation for 'Title' failed on the 'required' tag",
		},
		{
			name:       "Missing ID",
			input:      `{"title": "Test Task"}`,
			wantStatus: http.StatusBadRequest,
			wantMsg:    types.MessageIDRequired,
		},
		{
			name:       "Task Not Found",
			input:      `{"title": "Test Task", "id": "` + uuid.New().String() + `"}`,
			wantStatus: http.StatusNotFound,
			wantMsg:    types.MessageTaskNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = &verification.TaskValidator{Validator: validator.New()}

			req := httptest.NewRequest(http.MethodPut, "/tasks", strings.NewReader(tt.input))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := EditTask(c)
			if tt.wantStatus == http.StatusOK {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantStatus, rec.Code)
				var response types.TaskResponse
				err = json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantMsg, response.Message)
			} else {
				if err != nil {
					assert.Contains(t, err.Error(), tt.wantMsg)
				} else {
					var response types.TaskResponse
					err = json.Unmarshal(rec.Body.Bytes(), &response)
					assert.NoError(t, err)
					assert.Equal(t, tt.wantMsg, response.Message)
				}
			}
		})
	}
}
