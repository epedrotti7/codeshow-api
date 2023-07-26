package auth

import (
	"net/http"
	"time"

	"github.com/epedrotti7/codeshow-api/internal/errors"
	repository "github.com/epedrotti7/codeshow-api/internal/repositories"
	"github.com/epedrotti7/codeshow-api/internal/structs"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func CreateSession(c echo.Context) error {

	user := new(structs.User)
	c.Bind(user)

	user, err := repository.FindUserByEmail(user.Email)

	if err != nil {
		return err
	}

	token, err := CreateToken(user.ID.Hex())

	if err != nil {
		return errors.Validate(err, c)
	}

	user.AuthToken = token

	return c.JSON(http.StatusCreated, user)
}

func CreateToken(userID string) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userID
	atClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	token, err := at.SignedString([]byte("f6df78ee7edd79b057024c2c922b128a"))
	if err != nil {
		return "", err
	}
	return token, nil
}
