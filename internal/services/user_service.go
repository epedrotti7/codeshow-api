package service

import (
	"log"

	repository "github.com/epedrotti7/codeshow-api/internal/repositories"
	"github.com/epedrotti7/codeshow-api/internal/structs"
)

func Create(user *structs.User) *structs.User {
	userResponse := repository.Create(user)
	return userResponse
}

func FindUserByID(id string) *structs.User {

	userResponse, err := repository.FindUserByID(id)
	questions, err := repository.FindQuestionsByUserId(id)

	userResponse.Questions = questions

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return userResponse
}
