package repository

import (
	"github.com/handlename/otomo/internal/domain/entity"
)

type Session interface {
	Restore(entity.SessionID) (entity.Session, error)
	Save(entity.Session) error
}
