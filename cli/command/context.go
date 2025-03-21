package command

import (
	"context"

	otomo "github.com/handlename/otomo"
)

type Context struct {
	Ctx context.Context
	App *otomo.App
}
