package reasoning

// Answer is a value object representing the outcome of Brain reasoning.
type Answer struct {
	body string
}

func (ans *Answer) Body() string { return ans.body }

func NewAnswer(instruction string) *Answer {
	return &Answer{
		body: instruction,
	}
}
