package errors

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ValidationErrorResponse struct {
	Errors []ValidationError `json:"errors"`
}

type ValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

func Validate(err error, c echo.Context) error {
	var validationErrors []ValidationError

	// Verificando se é um erro de validação
	if _, ok := err.(validator.ValidationErrors); ok {
		for _, validationErr := range err.(validator.ValidationErrors) {
			ve := ValidationError{
				Field: validationErr.Field(),
				Error: validationErr.Error(),
			}
			validationErrors = append(validationErrors, ve)
		}
	} else {
		// Isso significa que é um erro normal, não um erro de validação
		ve := ValidationError{
			Field: "Erro Personalizado",
			Error: err.Error(),
		}
		validationErrors = append(validationErrors, ve)
	}

	response := ValidationErrorResponse{
		Errors: validationErrors,
	}

	return c.JSON(http.StatusBadRequest, response)
}
