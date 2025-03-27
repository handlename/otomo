package entity

// Context is contextual knowledge linked to a Session used by Otomo when interpreting Instructions.
type Context struct {
	body string
}

func NewContext(body string) *Context {
	return &Context{
		body: body,
	}
}
