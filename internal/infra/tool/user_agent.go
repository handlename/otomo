package tool

import (
	"fmt"

	"github.com/handlename/otomo"
)

func UserAgent() string {
	return fmt.Sprintf("otomo-bot/%s", otomo.Version)
}
