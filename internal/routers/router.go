package router

import (
	"github.com/epedrotti7/codeshow-api/internal/auth"
	handler "github.com/epedrotti7/codeshow-api/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func StartRouter(e *echo.Echo) {
	route := e.Group("/api/v1")

	route.POST("/session", auth.CreateSession)

	route.POST("/user", handler.Create)

	route.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte("f6df78ee7edd79b057024c2c922b128a"),
	}))

	route.GET("/user/:id", handler.FindUserByID)

	route.POST("/question", handler.CreateQuestionByUserId)

	route.POST("/answer/:id", handler.CompareAnswerById)
}
