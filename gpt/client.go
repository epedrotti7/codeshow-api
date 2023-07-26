package gpt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/epedrotti7/codeshow-api/internal/structs"
	"github.com/labstack/echo/v4"
)

func GetQuestionChatGPT(tecnologia string, nivel string, c echo.Context) (structs.Question, error) {
	apiURL := "https://api.openai.com/v1/chat/completions"
	authToken := "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"

	questionToGPT := "Por favor, faça uma pergunta sobre " + tecnologia +
		" de nível " + nivel + ". Formate sua resposta da seguinte maneira: " +
		"Primeiro, apresente a pergunta. Em seguida, apresente quatro alternativas de resposta, " +
		"cada uma iniciada por 'a)', 'b)', 'c)' ou 'd)'. " +
		"Finalmente, indique a alternativa correta com a frase 'A resposta correta é: ' seguida da letra correspondente. " +
		"Não inclua nenhum texto após a resposta correta. " +
		"Por exemplo: 'Qual é a cor do céu? a) Verde b) Azul c) Vermelho d) Amarelo. A resposta correta é: b)'."

	requestDataQuestion := structs.ChatGPTRequest{
		Model: "gpt-3.5-turbo",
		Messages: []structs.Message{
			{
				Role:    "user",
				Content: questionToGPT,
			},
		},
	}

	re := regexp.MustCompile(`(?s)(.*?)(?:\n\n)(a\) .*?\nb\) .*?\nc\) .*?\nd\) .*?)(?:\n\n)(?:A resposta correta é: )(.*?\))`)

	for i := 0; i < 5; i++ { // Limitando para 5 tentativas

		jsonDataQuestion, err := json.Marshal(requestDataQuestion)
		if err != nil {
			return structs.Question{}, err
		}

		requestQuestion, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonDataQuestion))
		if err != nil {
			return structs.Question{}, err
		}

		requestQuestion.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
		requestQuestion.Header.Set("Content-Type", "application/json")

		client := http.Client{}
		resp, err := client.Do(requestQuestion)
		if err != nil {
			return structs.Question{}, err
		}

		defer resp.Body.Close()

		responseData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return structs.Question{}, err
		}

		var chatResponseQuestion structs.ChatGPTResponse
		err = json.Unmarshal(responseData, &chatResponseQuestion)

		if err != nil {
			return structs.Question{}, err
		}

		fmt.Println(chatResponseQuestion.Choices[0].Message.Content)

		match := re.FindStringSubmatch(chatResponseQuestion.Choices[0].Message.Content)

		if len(match) == 4 {
			question := strings.TrimSpace(match[1])
			alternatives := strings.Split(match[2], "\n")
			answer := strings.TrimSpace(match[3])

			if len(alternatives) == 4 && question != "" && answer != "" { // check if alternatives are not null and question & answer are not empty
				return structs.Question{
					Question:     question,
					Answer:       answer,
					Alternatives: alternatives,
				}, nil
			}
		}

		time.Sleep(2 * time.Second) // wait before the next attempt
	}

	return structs.Question{}, errors.New("Failed to get a valid response from GPT-3 after 5 attempts")
}
