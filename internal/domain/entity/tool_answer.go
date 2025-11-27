package entity

type ToolAnswer interface {
	Body() string
}

type toolAnswer struct {
	body string
}

func NewToolAnswer(body string) ToolAnswer {
	return &toolAnswer{body: body}
}

func (t *toolAnswer) Body() string {
	return t.body
}
