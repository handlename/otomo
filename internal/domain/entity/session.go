package entity

type SessionID string

// Session represents a series of exchanges between User(s) and an Otomo.
type Session interface {
	UpdateContext(Context)
}

var _ Session = (*session)(nil)

type session struct {
	context Context
}

func (s *session) Context() Context { return s.context }

func NewSession() Session {
	return &session{
		context: nil,
	}
}

func (s *session) UpdateContext(context Context) {
	s.context = context
}
