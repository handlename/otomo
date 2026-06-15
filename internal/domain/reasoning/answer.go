package reasoning

import "fmt"

type AnswerBody string

// Answer is a value object representing the outcome of Brain reasoning.
type Answer struct {
	body AnswerBody
}

func (ans *Answer) Body() AnswerBody { return ans.body }

func NewAnswer(body AnswerBody) (*Answer, error) {
	if body == "" {
		return nil, fmt.Errorf("answer body is required")
	}
	return &Answer{
		body: body,
	}, nil
}
