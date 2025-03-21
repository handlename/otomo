package main

import (
	"os"

	"github.com/handlename/otomo/cli"
)

func main() {
	os.Exit(int(cli.Run()))
}
