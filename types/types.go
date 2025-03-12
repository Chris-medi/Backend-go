package types

import (
	"github.com/google/uuid"
)

type Task struct {
	Id          uuid.UUID `json:"id"`
	Title       string    `json:"title" validate:"required"`
	Description *string   `json:"description"`
	Status      *string   `json:"status"`
}

type TasksResponse struct {
	Message string `json:"message"`
	Data    []Task `json:"data,omitempty"`
}

type TaskResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
