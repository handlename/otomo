package usecase

import (
	"context"
	"fmt"

	"github.com/handlename/otomo/internal/domain/chat"
	"github.com/handlename/otomo/internal/domain/reasoning"
	"github.com/handlename/otomo/internal/errorcode"
	"github.com/morikuni/failure/v2"
	"github.com/handlename/otomo/internal/infra/trace"
	"go.opentelemetry.io/otel"
)

func ExecuteToolLoop(ctx context.Context, otomo *chat.Otomo, c *reasoning.Context, tools []reasoning.Tool) (*reasoning.Answer, error) {
	ctx, span := otel.Tracer("otomo").Start(ctx, "ExecuteToolLoop")
	defer span.End()

	turns := 0
	for {
		if err := ctx.Err(); err != nil {
			trace.RecordError(span, err)
			return nil, failure.Wrap(err)
		}

		if !reasoning.ShouldContinueToUseTool(turns) {
			err := failure.New(errorcode.ErrInternal, failure.Message("too many tool execution turns"))
			trace.RecordError(span, err)
			return nil, err
		}
		turns++

		ans, err := otomo.Think(ctx, c)
		if err != nil {
			err = failure.Wrap(err,
				failure.WithCode(errorcode.ErrInternal),
				failure.Message("failed to think"),
			)
			trace.RecordError(span, err)
			return nil, err
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
					err = failure.Wrap(err, failure.WithCode(errorcode.ErrInternal), failure.Message("failed to create tool result"))
					trace.RecordError(span, err)
					return nil, err
				}
				results = append(results, tr)
				continue
			}

			toolCtx, toolSpan := otel.Tracer("otomo").Start(ctx, "Tool Execute: "+tc.Name().Value())
			out, executeErr := tool.Execute(toolCtx, tc.InputJSON())
			trace.RecordError(toolSpan, executeErr)
			toolSpan.End()

			if executeErr != nil {
				tr, err := reasoning.NewToolResult(
					tc.ID(),
					fmt.Sprintf("error executing tool: %v", executeErr),
					reasoning.ToolResultError,
				)
				if err != nil {
					err = failure.Wrap(err, failure.WithCode(errorcode.ErrInternal), failure.Message("failed to create tool result"))
					trace.RecordError(span, err)
					return nil, err
				}
				results = append(results, tr)
			} else {
				tr, err := reasoning.NewToolResult(
					tc.ID(),
					out,
					reasoning.ToolResultSuccess,
				)
				if err != nil {
					err = failure.Wrap(err, failure.WithCode(errorcode.ErrInternal), failure.Message("failed to create tool result"))
					trace.RecordError(span, err)
					return nil, err
				}
				results = append(results, tr)
			}
		}

		if err := c.AddToolUseResponse(string(ans.Body()), ans.ToolCalls()); err != nil {
			err = failure.Wrap(err, failure.WithCode(errorcode.ErrInternal), failure.Message("failed to update context with tool calls"))
			trace.RecordError(span, err)
			return nil, err
		}
		if err := c.AddToolResults(results); err != nil {
			err = failure.Wrap(err, failure.WithCode(errorcode.ErrInternal), failure.Message("failed to update context with tool results"))
			trace.RecordError(span, err)
			return nil, err
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
