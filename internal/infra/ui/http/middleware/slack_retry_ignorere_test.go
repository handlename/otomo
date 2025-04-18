package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SlackRetryIgnorere(t *testing.T) {
	tests := []struct {
		name                    string
		retryNumHeader          string
		retryReasonHeader       string
		expectNextHandlerCalled bool
		expectResponseBody      string
	}{
		{
			name:                    "with retry num header, next handler is not called",
			retryNumHeader:          "1",
			retryReasonHeader:       "timeout",
			expectNextHandlerCalled: false,
			expectResponseBody:      "ignored",
		},
		{
			name:                    "without retry num header, next handler is called",
			retryNumHeader:          "",
			retryReasonHeader:       "",
			expectNextHandlerCalled: true,
			expectResponseBody:      "next handler called",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare handler
			nextHandlerCalled := false
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextHandlerCalled = true
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("next handler called"))
			})
			middleware := NewSlackRetryIgnorere()
			handler := middleware.Wrap(nextHandler)

			// Prepare request
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			req.Header.Set("x-slack-retry-num", tt.retryNumHeader)
			req.Header.Set("x-slack-retry-reason", tt.retryReasonHeader)
			rw := httptest.NewRecorder()

			// Call the handler
			handler.ServeHTTP(rw, req)

			// Check
			assert.Equal(t, http.StatusOK, rw.Code, "status code should match expected value")
			assert.Equal(t, tt.expectNextHandlerCalled, nextHandlerCalled, "nextHandler should be called correctly")
			assert.Equal(t, tt.expectResponseBody, rw.Body.String(), "response body should match expected value")
		})
	}
}
