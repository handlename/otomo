package command

import "github.com/alecthomas/kong"

type Root struct {
	ConfigPath string `help:"Path to config file." default:"config.toml"`

	LogLevel   string           `help:"Set log level. (trace|debug|info|warn|error|panic)" default:"info"`
	LogConsole bool             `help:"Write logs for console." default:"true" negatable:""`
	Version    kong.VersionFlag `help:"Show version."`

	Server Server `cmd:"" help:"Run as server"`
}
