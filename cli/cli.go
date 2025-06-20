package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/alecthomas/kong"
	"github.com/handlename/otomo"
	"github.com/handlename/otomo/cli/command"
	"github.com/handlename/otomo/config"
	"github.com/morikuni/failure/v2"
	"github.com/rs/zerolog/log"
)

type ExitCode int

const (
	ExitCodeOK    ExitCode = 0
	ExitCodeError ExitCode = 1
)

func Run() ExitCode {
	var root command.Root
	ktx := kong.Parse(&root, kong.Vars{"version": fmt.Sprintf("otomo v%s", otomo.Version)})

	if err := config.Load(root.ConfigPath); err != nil {
		handleError(err)
		return ExitCodeError
	}

	otomo.InitLogger(root.LogLevel, root.LogConsole)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := ktx.Run(&command.Context{Ctx: ctx, App: &otomo.App{}}); err != nil {
		if errors.Is(err, context.Canceled) {
			log.Error().Msg("canceled")
		} else {
			handleError(err)
		}

		return ExitCodeError
	}

	return ExitCodeOK
}

func handleError(err error) {
	fmt.Println("======== error ========")

	code := failure.CodeOf(err)
	fmt.Printf("code = %s\n", code)

	msg := failure.MessageOf(err)
	fmt.Printf("message = %s\n", msg)

	cs := failure.CallStackOf(err)
	fmt.Printf("callstack = %s\n", cs)

	fmt.Printf("cause = %s\n", failure.CauseOf(err))

	fmt.Println()
	fmt.Println("======== detail ========")
	fmt.Printf("%+v\n", err)
}
