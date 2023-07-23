package router

import (
	handler "github.com/epedrotti7/codeshow-api/internal/handlers"
	"github.com/labstack/echo/v4"
)

func StartRouter(e *echo.Echo) {
	route := e.Group("/api/v1")

	route.POST("/user", handler.Create)
	route.GET("/user/:id", handler.FindUserByID)

	route.POST("/question/:userId", handler.CreateQuestionByUserId)

	route.POST("/answer/:id", handler.CompareAnswerById)
}
