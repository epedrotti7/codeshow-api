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

	// Chamada a função e obtém os canais
	questionCh, errCh := gpt.GetQuestionChatGPT(question.Tecnologia, question.Nivel, c)

	// Seleciona o primeiro canal que tem um retorno
	select {
	case questionResponseGPT := <-questionCh:
		// Se o canal de perguntas retornar primeiro, processamos a pergunta
		questionFormatted := structs.Question{
			UserId:       userId,
			Question:     questionResponseGPT.Question,
			Answer:       questionResponseGPT.Answer,
			Alternatives: questionResponseGPT.Alternatives,
		}

		questionSaved := repository.CreateQuestionByUserId(&questionFormatted)

		return questionSaved, nil

	case err := <-errCh:
		// Se o canal de erro retornar primeiro, retornamos o erro
		return nil, err
	}
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
