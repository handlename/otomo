package command

import "github.com/alecthomas/kong"

type Root struct {
	Port     int              `help:"Port for listen." default:"8080"`
	LogLevel string           `help:"Set log level. (trace|debug|info|warn|error|panic)" default:"info"`
	Version  kong.VersionFlag `help:"Show version."`
}

func (r *Root) Run(c *Context) error {
	return c.App.Run(c.Ctx)
}
