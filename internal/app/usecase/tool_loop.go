package usecase

import (
	"context"
	"fmt"

	"github.com/handlename/otomo/internal/domain/chat"
	"github.com/handlename/otomo/internal/domain/reasoning"
	"github.com/handlename/otomo/internal/errorcode"
	"github.com/morikuni/failure/v2"
)

func ExecuteToolLoop(ctx context.Context, otomo *chat.Otomo, c *reasoning.Context, tools []reasoning.Tool) (*reasoning.Answer, error) {
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
			return nil, failure.Wrap(err,
				failure.WithCode(errorcode.ErrInternal),
				failure.Message("failed to think"),
			)
		}

		if !ans.HasToolCalls() {
			return ans, nil
		}

		var results []reasoning.ToolResult
		for _, tc := range ans.ToolCalls() {
			tool, ok := findToolInList(tools, tc.Name())
			if !ok {
				tr, err := reasoning.NewToolResult(
					tc.ID(),
					fmt.Sprintf("error: tool '%s' not found", tc.Name().Value()),
					reasoning.ToolResultError,
				)
				if err != nil {
					return nil, failure.Wrap(err, failure.WithCode(errorcode.ErrInternal), failure.Message("failed to create tool result"))
				}
				results = append(results, tr)
				continue
			}

			out, err := tool.Execute(ctx, tc.InputJSON())
			if err != nil {
				tr, err := reasoning.NewToolResult(
					tc.ID(),
					fmt.Sprintf("error executing tool: %v", err),
					reasoning.ToolResultError,
				)
				if err != nil {
					return nil, failure.Wrap(err, failure.WithCode(errorcode.ErrInternal), failure.Message("failed to create tool result"))
				}
				results = append(results, tr)
			} else {
				tr, err := reasoning.NewToolResult(
					tc.ID(),
					out,
					reasoning.ToolResultSuccess,
				)
				if err != nil {
					return nil, failure.Wrap(err, failure.WithCode(errorcode.ErrInternal), failure.Message("failed to create tool result"))
				}
				results = append(results, tr)
			}
		}

		if err := c.AddToolUseResponse(string(ans.Body()), ans.ToolCalls()); err != nil {
			return nil, failure.Wrap(err, failure.WithCode(errorcode.ErrInternal), failure.Message("failed to update context with tool calls"))
		}
		if err := c.AddToolResults(results); err != nil {
			return nil, failure.Wrap(err, failure.WithCode(errorcode.ErrInternal), failure.Message("failed to update context with tool results"))
		}
	}
}

func findToolInList(tools []reasoning.Tool, name reasoning.ToolName) (reasoning.Tool, bool) {
	for _, t := range tools {
		if t.Name().Equals(name) {
			return t, true
		}
	}
	return nil, false
}
