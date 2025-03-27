package entity

// Session represents a series of exchanges between User(s) and an Otomo.
type Session struct {
	context      Context
	instructions []*Instruction
}

func (s *Session) Context() Context { return s.context }

func NewSession() *Session {
	return &Session{
		context:      Context{},
		instructions: make([]*Instruction, 0),
	}
}

func (s *Session) UpdateContext(context Context) {
	s.context = context
}

func (s *Session) AddInstruction(instruction *Instruction) {
	s.instructions = append(s.instructions, instruction)
}
