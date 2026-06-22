package terminal

import (
	"context"
	"fmt"

	"github.com/handlename/otomo/internal/domain/chat"
	"github.com/handlename/otomo/internal/domain/reasoning"
)

func StartChatLoop(ctx context.Context, otomo *chat.Otomo, tools []reasoning.Tool) error {
	fmt.Println("Chat loop started!")
	return nil
}
