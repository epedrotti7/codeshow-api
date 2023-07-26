package service

import (
	"regexp"
	"strings"

	"github.com/epedrotti7/codeshow-api/gpt"
	repository "github.com/epedrotti7/codeshow-api/internal/repositories"
	"github.com/epedrotti7/codeshow-api/internal/structs"
	"github.com/labstack/echo/v4"
)

func CreateQuestionByUserId(question *structs.QuestionRequest, userId string, c echo.Context) (*structs.Question, error) {

	questionResponseGPT, _ := gpt.GetQuestionChatGPT(question.Tecnologia, question.Nivel, c)

	questionFormatted := structs.Question{
		UserId:       userId,
		Question:     questionResponseGPT.Question,
		Answer:       questionResponseGPT.Answer,
		Alternatives: questionResponseGPT.Alternatives,
	}

	questionSaved := repository.CreateQuestionByUserId(&questionFormatted)

	return questionSaved, nil

}

func CompareAnswerById(answerReq *structs.Answer, id string, userId string) (*structs.QuestionResponse, error) {
	questionAnswer, err := repository.FindQuestionById(id)

	if err != nil {
		return nil, err
	}

	repository.UpdateQuestionById(answerReq, id)

	correctAnswer := extractAnswer(questionAnswer.Answer)
	if correctAnswer == strings.ToLower(answerReq.Answer) {

		repository.UpdateUserScore(userId)

		return &structs.QuestionResponse{
			Message:     "resposta correta",
			CorrectCode: 1,
		}, nil
	}

	return &structs.QuestionResponse{
		Message:     "resposta incorreta",
		CorrectCode: 0,
	}, nil
}

func extractAnswer(answer string) string {
	re := regexp.MustCompile(`[a-dA-D]\)?`)
	matches := re.FindStringSubmatch(answer)
	if len(matches) > 0 {
		letter := matches[0]
		if strings.Contains(letter, ")") {
			return strings.ToLower(letter[0:1])
		} else {
			return strings.ToLower(letter)
		}
	}
	return ""
}
