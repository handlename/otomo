package mcp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/zerolog/log"
)

type Server struct {
	port int
	p    *tea.Program
	srv  *http.Server
}

type PostMessageInput struct {
	Message string `json:"message" jsonschema:"The message to send to the otomo chat room"`
}

type PostMessageOutput struct {
	Response string `json:"response" jsonschema:"The response from otomo"`
}

type MCPResponse struct {
	Response string
	Error    error
}

type MCPRequestMsg struct {
	Prompt    string
	ReplyChan chan MCPResponse
}

func NewServer(port int, p *tea.Program) *Server {
	return &Server{
		port: port,
		p:    p,
	}
}

func (s *Server) Start(ctx context.Context) error {
	impl := &mcp.Implementation{
		Name:    "otomo-mcp",
		Version: "1.0.0",
	}
	server := mcp.NewServer(impl, nil)

	// Add the post_message tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "post_message",
		Description: "Post a message to the active otomo chat TUI session and retrieve the response.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input PostMessageInput) (*mcp.CallToolResult, PostMessageOutput, error) {
		replyChan := make(chan MCPResponse, 1)

		s.p.Send(MCPRequestMsg{
			Prompt:    input.Message,
			ReplyChan: replyChan,
		})

		select {
		case reply := <-replyChan:
			if reply.Error != nil {
				return nil, PostMessageOutput{}, reply.Error
			}
			return nil, PostMessageOutput{Response: reply.Response}, nil
		case <-ctx.Done():
			return nil, PostMessageOutput{}, ctx.Err()
		case <-time.After(2 * time.Minute):
			return nil, PostMessageOutput{}, fmt.Errorf("timeout waiting for otomo response")
		}
	})

	handler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		return server
	}, nil)

	s.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: handler,
	}

	go func() {
		log.Info().Msgf("Starting MCP Server at http://localhost:%d/sse", s.port)
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("MCP Server failed")
		}
	}()

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = s.srv.Shutdown(shutdownCtx)
	}()

	return nil
}
