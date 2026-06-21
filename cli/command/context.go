package command

import (
	"context"

	"github.com/handlename/otomo/internal/infra/app"
)

type Context struct {
	Ctx context.Context
	App *app.App
}
