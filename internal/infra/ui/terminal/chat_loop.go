package terminal

import (
	"context"
	"errors"
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/handlename/otomo/internal/app/usecase"
	"github.com/handlename/otomo/internal/domain/chat"
	"github.com/handlename/otomo/internal/domain/core"
	"github.com/handlename/otomo/internal/domain/reasoning"
	"github.com/morikuni/failure/v2"
)

func StartChatLoop(ctx context.Context, otomo *chat.Otomo, tools []reasoning.Tool) error {
	var history []*core.Message

	userStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true) // Green
	aiStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Bold(true)   // Blue
	infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Italic(true) // Gray

	fmt.Println(infoStyle.Render("Otomo interactive chat session started. Type exit or empty string to quit."))

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		var userInput string
		err := huh.NewInput().
			Title(userStyle.Render("You")).
			Prompt("> ").
			Value(&userInput).
			Run()
		if err != nil {
			if err == huh.ErrUserAborted {
				fmt.Println(infoStyle.Render("Exited chat session."))
				return nil
			}
			return failure.Wrap(err)
		}

		if userInput == "" || userInput == "exit" || userInput == "quit" {
			fmt.Println(infoStyle.Render("Exited chat session."))
			return nil
		}

		// Prepare Reasoning Context
		c := reasoning.NewContext()
		if len(history) > 0 {
			if err := c.SetMessages(history); err != nil {
				return failure.Wrap(err)
			}
		}
		c.SetUserPrompt(core.PromptBody(userInput))
		c.SetTools(tools)

		fmt.Print(infoStyle.Render("Otomo is thinking..."))

		// Execute Bedrock Reasoning
		ans, err := usecase.ExecuteToolLoop(ctx, otomo, c, tools)
		// Clear thinking text
		fmt.Print("\r\033[K")

		if err != nil {
			if failure.Is(err, context.Canceled) || errors.Is(err, context.Canceled) {
				return err
			}
			fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(fmt.Sprintf("Error: %v", err)))
			continue
		}

		// Output AI response
		fmt.Println(aiStyle.Render("Otomo:"))
		fmt.Println(ans.Body())
		fmt.Println()

		// Record to memory history
		userMsg, err := core.NewMessage(core.RoleUser, core.UserID{}, core.MessageBody(userInput))
		if err != nil {
			return failure.Wrap(err)
		}
		aiMsg, err := core.NewMessage(core.RoleAssistant, core.UserID{}, core.MessageBody(ans.Body()))
		if err != nil {
			return failure.Wrap(err)
		}
		history = append(history, userMsg, aiMsg)
	}
}

