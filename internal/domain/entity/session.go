package entity

// Session represents a series of exchanges between User(s) and an Otomo.
type Session struct {
	context Context
}

func (s *Session) Context() Context { return s.context }

func NewSession() *Session {
	return &Session{
		context: nil,
	}
}

func (s *Session) UpdateContext(context Context) {
	s.context = context
}
