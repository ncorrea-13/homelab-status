package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		serverSecret   string
		headerSecret   string
		querySecret    string
		expectedStatus int
	}{
		{
			name:           "correct secret via header",
			serverSecret:   "mysecret",
			headerSecret:   "mysecret",
			querySecret:    "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "correct secret via query param",
			serverSecret:   "mysecret",
			headerSecret:   "",
			querySecret:    "mysecret",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "server has no secret configured, client sends one",
			serverSecret:   "",
			headerSecret:   "mysecret",
			querySecret:    "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "server has no secret, client sends nothing",
			serverSecret:   "",
			headerSecret:   "",
			querySecret:    "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "client sends nothing, server has secret",
			serverSecret:   "mysecret",
			headerSecret:   "",
			querySecret:    "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "incorrect secret",
			serverSecret:   "mysecret",
			headerSecret:   "wrongsecret",
			querySecret:    "",
			expectedStatus: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authService := NewAuth(tt.serverSecret)
			url := "/webhook/kuma"
			if tt.querySecret != "" {
				url = url + "?secret=" + tt.querySecret
				/* armá acá el string con el query param */
			}

			req := httptest.NewRequest(http.MethodPost, url, nil)
			if tt.headerSecret != "" {
				req.Header.Set("X-Webhook-Secret", tt.headerSecret)
			}

			rec := httptest.NewRecorder()

			nextCalled := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusOK)
			})

			authService.Middleware(next).ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
			expectedNextCalled := tt.expectedStatus == http.StatusOK
			if nextCalled != expectedNextCalled {
				t.Errorf("expected nextCalled=%v, got %v", expectedNextCalled, nextCalled)
			}
		})
	}
}
