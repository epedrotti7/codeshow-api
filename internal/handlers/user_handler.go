package handler

import (
	"net/http"

	"github.com/epedrotti7/codeshow-api/internal/auth"
	"github.com/epedrotti7/codeshow-api/internal/errors"
	service "github.com/epedrotti7/codeshow-api/internal/services"
	"github.com/epedrotti7/codeshow-api/internal/structs"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func Create(c echo.Context) error {

	user := new(structs.User)

	// O método Bind() usa o cabeçalho "Content-Type" da requisição para determinar como ler os dados.
	// Ele suporta "application/json", "application/xml" e "application/x-www-form-urlencoded" por padrão.
	c.Bind(user)

	validate := validator.New()
	err := validate.Struct(user)

	if err != nil {
		return errors.Validate(err, c)
	}

	userResponse := service.Create(user)

	var id string
	if user.ID != nil {
		id = user.ID.String()
	}

	token, err := auth.CreateToken(id)

	userResponse.AuthToken = token

	return c.JSON(http.StatusCreated, userResponse)
}

func FindUserByID(c echo.Context) error {
	id := c.Param("id")
	userResponse := service.FindUserByID(id)
	return c.JSON(http.StatusOK, userResponse)
}
