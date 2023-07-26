package structs

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	AuthToken string              `json:"authToken" bson:"id,omitempty"`
	Name      string              `json:"name" validate:"required"`
	Email     string              `json:"email" validate:"required,email"`
	Password  string              `json:"password" validate:"required,min=8"`
	Score     int64               `json:"score"`
	Questions []*Question         `json:"questions"`
}

type Question struct {
	Id              string   `json:"id,omitempty" bson:"id,omitempty"`
	UserId          string   `json:"userId,omitempty"`
	Question        string   `json:"question,omitempty"`
	Answer          string   `json:"answer,omitempty"`
	Alternatives    []string `json:"alternatives" bson:"-"`
	Message         string   `json:"message"`
	UserAlternative string   `json:"userAlternative,omitempty" bson:"userAlternative,omitempty"`
}

type Answer struct {
	Answer string `json:"answer"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatGPTRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ChatGPTResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type QuestionRequest struct {
	Nivel      string `json:"nivel"`
	Tecnologia string `json:"tecnologia"`
}

type QuestionResponse struct {
	Message     string `json:"message"`
	CorrectCode int    `json:"correctCode"`
}
