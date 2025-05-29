package repository

import (
	"github.com/handlename/otomo/internal/domain/entity"
	drepo "github.com/handlename/otomo/internal/domain/repository"
	"github.com/handlename/otomo/internal/domain/valueobject"
)

var _ drepo.Session = (*VolatileSession)(nil)

type VolatileSession struct {
	session entity.Session
}

// Restore implements repository.Session.
func (v *VolatileSession) Restore(valueobject.SessionID) (entity.Session, error) {
	if v.session == nil {
		v.session = entity.NewSession()
	}

	return v.session, nil
}

// Save implements repository.Session.
func (v *VolatileSession) Save(sess entity.Session) error {
	v.session = sess
	return nil
}
