package repository

import (
	"github.com/handlename/otomo/internal/domain/entity"
)

type Context interface {
	New() entity.Context
}
