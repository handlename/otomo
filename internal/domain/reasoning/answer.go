package reasoning

import "fmt"

// Answer is a value object representing the outcome of Brain reasoning.
type Answer struct {
	body string
}

func (ans *Answer) Body() string { return ans.body }

func NewAnswer(body string) (*Answer, error) {
	if body == "" {
		return nil, fmt.Errorf("answer body is required")
	}
	return &Answer{
		body: body,
	}, nil
}
