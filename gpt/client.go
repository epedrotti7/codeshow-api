package gpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/epedrotti7/codeshow-api/internal/errors"
	"github.com/epedrotti7/codeshow-api/internal/structs"
	"github.com/labstack/echo/v4"
)

func GetQuestionChatGPT(tecnologia string, nivel string, c echo.Context) (<-chan structs.Question, <-chan error) {
	// Criamos dois canais, um para a questão e outro para um erro potencial
	questionCh := make(chan structs.Question, 1)
	errCh := make(chan error, 1)

	go func() {
		// Esta é a mesma lógica que você tinha antes
		// Só que agora estamos em uma goroutine separada

		apiURL := "https://api.openai.com/v1/chat/completions"
		authToken := "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"

		questionToGPT := "Faça uma pergunta sobre " + tecnologia +
			" de nível " + nivel + " com 4 alternativas de resposta " +
			"e retorne APENAS a letra da alternativa correta."

		requestDataQuestion := structs.ChatGPTRequest{
			Model: "gpt-3.5-turbo",
			Messages: []structs.Message{
				{
					Role:    "user",
					Content: questionToGPT,
				},
			},
		}

		jsonDataQuestion, err := json.Marshal(requestDataQuestion)
		if err != nil {
			errCh <- err
			return
		}

		requestQuestion, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonDataQuestion))

		if err != nil {
			errCh <- err
			return
		}

		requestQuestion.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
		requestQuestion.Header.Set("Content-Type", "application/json")

		client := http.Client{}
		resp, err := client.Do(requestQuestion)

		if err != nil {
			errCh <- err
			return
		}

		defer resp.Body.Close()

		responseData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errCh <- err
			return
		}

		// Converte a resposta JSON para a estrutura de resposta
		var chatResponseQuestion structs.ChatGPTResponse
		err = json.Unmarshal(responseData, &chatResponseQuestion)
		if err != nil {
			errCh <- err
			return
		}

		if len(chatResponseQuestion.Choices) > 0 {

			var question, answer string
			var alternatives []string

			// Expressão regular que captura a palavra "correta".
			re := regexp.MustCompile(`(?i)resposta`)

			parts := re.Split(chatResponseQuestion.Choices[0].Message.Content, -1)

			if len(parts) < 2 {
				err := fmt.Errorf("Ocorreu um erro inesperado, por favor tente novamente.")
				errCh <- errors.Validate(err, c) // supondo que c é seu contexto do echo
				return
			}

			// Divide a primeira parte (pergunta e alternativas) por parágrafos.
			paragraphs := strings.Split(parts[0], "\n\n")

			if len(paragraphs) > 1 {
				question = strings.TrimSpace(paragraphs[0])
				alternatives = paragraphs[1:]
			}

			// A segunda parte contém a resposta. Nós a limpamos para remover espaços em branco e caracteres não desejados.
			answer = strings.TrimSpace(parts[1])
			answer = strings.Trim(answer, ": ")

			// Junta as alternativas em uma única string.
			alternativesJoined := strings.Join(alternatives, "\n\n")

			questionCh <- structs.Question{
				Question:     question,
				Answer:       answer,
				Alternatives: alternativesJoined,
			}
			return
		}

		errCh <- fmt.Errorf("Nenhuma pergunta retornada pela API do GPT-3")
	}()

	// Retornamos os canais. O chamador pode selecionar nos canais para receber a questão ou um erro.
	return questionCh, errCh
}

// func GetQuestionChatGPT(tecnologia string, nivel string, c echo.Context) (structs.Question, error) {

// 	apiURL := "https://api.openai.com/v1/chat/completions"
// 	authToken := "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"

// 	questionToGPT := "Faça uma pergunta sobre " + tecnologia +
// 		" de nível " + nivel + " com 4 alternativas de resposta " +
// 		"e retorne APENAS a letra da alternativa correta."

// 	requestDataQuestion := structs.ChatGPTRequest{
// 		Model: "gpt-3.5-turbo",
// 		Messages: []structs.Message{
// 			{
// 				Role:    "user",
// 				Content: questionToGPT,
// 			},
// 		},
// 	}

// 	jsonDataQuestion, err := json.Marshal(requestDataQuestion)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	requestQuestion, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonDataQuestion))

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	requestQuestion.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
// 	requestQuestion.Header.Set("Content-Type", "application/json")

// 	client := http.Client{}
// 	resp, err := client.Do(requestQuestion)

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer resp.Body.Close()

// 	responseData, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Converte a resposta JSON para a estrutura de resposta
// 	var chatResponseQuestion structs.ChatGPTResponse
// 	err = json.Unmarshal(responseData, &chatResponseQuestion)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if len(chatResponseQuestion.Choices) > 0 {

// 		var question, answer string
// 		var alternatives []string

// 		// Expressão regular que captura a palavra "correta".
// 		re := regexp.MustCompile(`(?i)resposta`)

// 		parts := re.Split(chatResponseQuestion.Choices[0].Message.Content, -1)

// 		if len(parts) < 2 {
// 			err := fmt.Errorf("Ocorreu um erro inesperado, por favor tente novamente.")
// 			return structs.Question{}, errors.Validate(err, c) // supondo que c é seu contexto do echo
// 		}

// 		// Divide a primeira parte (pergunta e alternativas) por parágrafos.
// 		paragraphs := strings.Split(parts[0], "\n\n")

// 		if len(paragraphs) > 1 {
// 			question = strings.TrimSpace(paragraphs[0])
// 			alternatives = paragraphs[1:]
// 		}

// 		// A segunda parte contém a resposta. Nós a limpamos para remover espaços em branco e caracteres não desejados.
// 		answer = strings.TrimSpace(parts[1])
// 		answer = strings.Trim(answer, ": ")

// 		// Junta as alternativas em uma única string.
// 		alternativesJoined := strings.Join(alternatives, "\n\n")

// 		return structs.Question{
// 			Question:     question,
// 			Answer:       answer,
// 			Alternatives: alternativesJoined,
// 		}, nil
// 	}

// 	return structs.Question{}, nil
// }
