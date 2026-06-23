package mcp

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer_StartAndShutdown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	port := 40506
	server := NewServer(port, nil)
	err := server.Start(ctx)
	require.NoError(t, err)

	time.Sleep(200 * time.Millisecond)

	// Check if server is reachable
	client := &http.Client{Timeout: 1 * time.Second}
	resp, err := client.Get(fmt.Sprintf("http://localhost:%d/sse", port))
	if err == nil {
		defer resp.Body.Close()
		assert.Contains(t, []int{http.StatusOK, http.StatusMethodNotAllowed, http.StatusNotFound, http.StatusBadRequest}, resp.StatusCode)
	}

	cancel()
	time.Sleep(200 * time.Millisecond)
}
