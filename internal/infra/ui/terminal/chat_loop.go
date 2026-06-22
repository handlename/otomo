package terminal

import (
	"context"
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/handlename/otomo/internal/domain/chat"
	"github.com/handlename/otomo/internal/domain/core"
	"github.com/handlename/otomo/internal/domain/reasoning"
	"github.com/handlename/otomo/internal/errorcode"
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
		ans, err := executeToolLoop(ctx, otomo, c, tools)
		// Clear thinking text
		fmt.Print("\r\033[K")

		if err != nil {
			fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(fmt.Sprintf("Error: %v", err)))
			continue
		}

		// Output AI response
		fmt.Println(aiStyle.Render("Otomo:"))
		fmt.Println(ans.Body())
		fmt.Println()

		// Record to memory history
		userMsg, _ := core.NewMessage(core.RoleUser, core.UserID{}, core.MessageBody(userInput))
		aiMsg, _ := core.NewMessage(core.RoleAssistant, core.UserID{}, core.MessageBody(ans.Body()))
		history = append(history, userMsg, aiMsg)
	}
}

// Copy local tool runner execution from usecase
func executeToolLoop(ctx context.Context, otomo *chat.Otomo, c *reasoning.Context, tools []reasoning.Tool) (*reasoning.Answer, error) {
	turns := 0
	for {
		if err := ctx.Err(); err != nil {
			return nil, failure.Wrap(err)
		}
		if !reasoning.ShouldContinueToUseTool(turns) {
			return nil, failure.New(errorcode.ErrInternal, failure.Message("too many tool execution turns"))
		}
		turns++

		ans, err := otomo.Think(ctx, c)
		if err != nil {
			return nil, failure.Wrap(err)
		}

		if !ans.HasToolCalls() {
			return ans, nil
		}

		var results []reasoning.ToolResult
		for _, tc := range ans.ToolCalls() {
			tool, ok := findTool(tools, tc.Name())
			if !ok {
				tr, _ := reasoning.NewToolResult(tc.ID(), "tool not found", reasoning.ToolResultError)
				results = append(results, tr)
				continue
			}
			out, err := tool.Execute(ctx, tc.InputJSON())
			if err != nil {
				tr, _ := reasoning.NewToolResult(tc.ID(), err.Error(), reasoning.ToolResultError)
				results = append(results, tr)
			} else {
				tr, _ := reasoning.NewToolResult(tc.ID(), out, reasoning.ToolResultSuccess)
				results = append(results, tr)
			}
		}

		_ = c.AddToolUseResponse(string(ans.Body()), ans.ToolCalls())
		_ = c.AddToolResults(results)
	}
}

func findTool(tools []reasoning.Tool, name reasoning.ToolName) (reasoning.Tool, bool) {
	for _, t := range tools {
		if t.Name().Equals(name) {
			return t, true
		}
	}
	return nil, false
}
