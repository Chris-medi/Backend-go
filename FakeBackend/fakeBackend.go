package FakeBackend

import (
	"backend/types"

	"github.com/google/uuid"
)

var data []types.Task

func GetAllData() []types.Task {
	return data
}

func GetTaskById(id uuid.UUID) *types.Task {
	for _, task := range data {
		if task.Id == id {
			return &task
		}
	}
	return nil
}

func CreateTask(task types.Task) types.Task {

	newId, err := uuid.NewUUID()
	if err != nil {
		// In a real application, you would handle this error properly
		// For now, we'll just use a random UUID as fallback
		newId = uuid.New()
	}
	task.Id = newId
	status := "pending"
	task.Status = &status

	data = append(data, task)
	return task
}

func DeleteTask(id uuid.UUID) bool {
	for i, task := range data {
		if task.Id == id {
			data = append(data[:i], data[i+1:]...)
			return true
		}
	}
	return false
}

func EditTask(task types.Task) bool {
	for i, _ := range data {
		if data[i].Id == task.Id {
			data[i] = task
			status := "completed"
			data[i].Status = &status
			return true
		}
	}
	return false
}
