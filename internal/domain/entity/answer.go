package entity

// Answer is the result of Otomo thinking based on the Instruction.
type Answer struct {
	body string
}

func (ans *Answer) Body() string { return ans.body }

func NewAnswer(instruction string) *Answer {
	return &Answer{
		body: string(instruction),
	}
}
