package repository

import (
	"github.com/handlename/otomo/internal/domain/entity"
	vo "github.com/handlename/otomo/internal/domain/valueobject"
)

type Session interface {
	Restore(vo.SessionID) (*entity.Session, error)
	Save(*entity.Session) error
}
