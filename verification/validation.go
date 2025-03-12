package verification

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type (
	Task struct {
		Title string `json:"title" validate:"required"`
	}

	TaskValidator struct {
		Validator *validator.Validate
	}
)

func (cv *TaskValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		// Return validation errors
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
