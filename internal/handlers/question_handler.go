package handler

import (
	"fmt"
	"net/http"

	"github.com/epedrotti7/codeshow-api/internal/errors"
	service "github.com/epedrotti7/codeshow-api/internal/services"
	"github.com/epedrotti7/codeshow-api/internal/structs"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func CreateQuestionByUserId(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"].(string)

	question := new(structs.QuestionRequest)

	// O método Bind() usa o cabeçalho "Content-Type" da requisição para determinar como ler os dados.
	// Ele suporta "application/json", "application/xml" e "application/x-www-form-urlencoded" por padrão.
	c.Bind(question)

	validate := validator.New()
	err := validate.Struct(question)

	if err != nil {
		return errors.Validate(err, c)
	}

	questionSavedResponse, err := service.CreateQuestionByUserId(question, userId, c)

	if err != nil {
		return errors.Validate(err, c)
	}

	return c.JSON(http.StatusCreated, questionSavedResponse)
}

func CompareAnswerById(c echo.Context) error {

	answer := new(structs.Answer)

	c.Bind(answer)

	id := c.Param("id")
	userId := c.Request().Header.Get("X-User-Id")

	validate := validator.New()
	err := validate.Struct(answer)

	if err != nil {
		return errors.Validate(err, c)
	}

	answerResponse, err := service.CompareAnswerById(answer, id, userId)

	if err != nil {
		// manipule o erro aqui
		fmt.Println(err)
		return nil
	}

	return c.JSON(http.StatusOK, answerResponse)
}
