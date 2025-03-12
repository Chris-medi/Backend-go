package handlers

import (
	"backend/FakeBackend"
	"backend/types"
	"math/rand"
	"net/http"
	"time"

	"sync"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// HealthCheck handles the health check endpoint
func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "healthy",
	})
}

// WorkSimulation simulates work by sleeping for a random duration
// Tambien podemos usar un canal
func WorkSimulation(wg *sync.WaitGroup) {
	defer wg.Done()
	duration := time.Duration(rand.Intn(3)+1) * time.Second
	time.Sleep(duration)
}

// GetAllTask handles the retrieval of all tasks
func GetAllTask(c echo.Context) error {
	data := FakeBackend.GetAllData()
	return c.JSON(http.StatusOK, types.TasksResponse{
		Message: types.MessageTasksRetrieved,
		Data:    data,
	})
}

// GetTaskById handles the retrieval of a task by ID
func GetTaskById(c echo.Context) error {
	idStr := c.Param("id")
	if idStr == "" {
		return c.JSON(http.StatusBadRequest, types.TaskResponse{
			Message: types.MessageIDRequired,
			Data:    nil,
		})
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.TaskResponse{
			Message: types.MessageInvalidIDFormat,
			Data:    nil,
		})
	}

	task := FakeBackend.GetTaskById(id)
	if task == nil {
		return c.JSON(http.StatusNotFound, types.TaskResponse{
			Message: types.MessageTaskNotFound,
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, types.TaskResponse{
		Message: types.MessageTaskRetrieved,
		Data:    task,
	})
}

// CreateTask handles the creation of a new task
func CreateTask(c echo.Context) error {
	u := new(types.Task)
	if err := c.Bind(u); err != nil {
		return err
	}
	if err := c.Validate(u); err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go WorkSimulation(&wg)
	wg.Wait()

	data := FakeBackend.CreateTask(*u)
	return c.JSON(http.StatusCreated, types.TaskResponse{
		Message: types.MessageTaskCreated,
		Data:    data,
	})
}

// DeleteTask handles the deletion of a task by ID
func DeleteTask(c echo.Context) error {
	idStr := c.Param("id")
	if idStr == "" {
		return c.JSON(http.StatusBadRequest, types.TaskResponse{
			Message: types.MessageIDRequired,
			Data:    nil,
		})
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.TaskResponse{
			Message: types.MessageInvalidIDFormat,
			Data:    nil,
		})
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go WorkSimulation(&wg)
	wg.Wait()

	if deleted := FakeBackend.DeleteTask(id); !deleted {
		return c.JSON(http.StatusNotFound, types.TaskResponse{
			Message: types.MessageTaskNotFound,
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, types.TaskResponse{
		Message: types.MessageTaskDeleted,
		Data:    nil,
	})
}

// EditTask handles the updating of a task
func EditTask(c echo.Context) error {
	u := new(types.Task)
	if err := c.Bind(u); err != nil {
		return err
	}
	if err := c.Validate(u); err != nil {
		return err
	}

	// Check if the UUID is valid (not zero)
	if u.Id == uuid.Nil {
		return c.JSON(http.StatusBadRequest, types.TaskResponse{
			Message: types.MessageIDRequired,
			Data:    nil,
		})
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go WorkSimulation(&wg)
	wg.Wait()
	if updated := FakeBackend.EditTask(*u); !updated {
		return c.JSON(http.StatusNotFound, types.TaskResponse{
			Message: types.MessageTaskNotFound,
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, types.TaskResponse{
		Message: types.MessageTaskUpdated,
		Data:    u,
	})
}
