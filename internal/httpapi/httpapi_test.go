package httpapi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Brilhante29/go-rate-limiter/internal/limiter"
)

type stubStore struct {
	decision limiter.Decision
	err      error
}

func (s stubStore) Take(context.Context, string, limiter.Policy, int64) (limiter.Decision, error) {
	return s.decision, s.err
}

func TestCheckAllowed(t *testing.T) {
	service, _ := limiter.NewService(stubStore{decision: limiter.Decision{Allowed: true, Remaining: 9}}, limiter.Policy{RatePerSecond: 10, Burst: 10})
	handler := New(service, nil, "node-a", "memory").Handler()
	request := httptest.NewRequest(http.MethodPost, "/v1/limits/tenant-a/check", strings.NewReader(`{"cost":1}`))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
	}
	if got := response.Header().Get("X-Rate-Limiter-Node"); got != "node-a" {
		t.Fatalf("node header=%q", got)
	}
}

func TestCheckRejected(t *testing.T) {
	service, _ := limiter.NewService(stubStore{decision: limiter.Decision{RetryAfter: 250 * time.Millisecond}}, limiter.Policy{RatePerSecond: 10, Burst: 10})
	handler := New(service, nil, "node-b", "redis").Handler()
	request := httptest.NewRequest(http.MethodPost, "/v1/limits/shared/check", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusTooManyRequests {
		t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
	}
	if got := response.Header().Get("Retry-After"); got != "1" {
		t.Fatalf("retry-after=%q", got)
	}
}

func TestCheckRejectsInvalidInput(t *testing.T) {
	service, _ := limiter.NewService(stubStore{}, limiter.Policy{RatePerSecond: 10, Burst: 10})
	handler := New(service, nil, "node-a", "memory").Handler()
	tests := []struct {
		name string
		url  string
		body string
	}{
		{name: "key", url: "/v1/limits/bad%20key/check", body: `{}`},
		{name: "cost", url: "/v1/limits/good/check", body: `{"cost":11}`},
		{name: "json", url: "/v1/limits/good/check", body: `{`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.url, strings.NewReader(test.body))
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != http.StatusBadRequest {
				t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
			}
		})
	}
}

func TestCheckFailsClosedWhenStoreIsUnavailable(t *testing.T) {
	service, _ := limiter.NewService(stubStore{err: errors.New("Redis unavailable")}, limiter.Policy{RatePerSecond: 10, Burst: 10})
	handler := New(service, nil, "node-a", "redis").Handler()
	request := httptest.NewRequest(http.MethodPost, "/v1/limits/shared/check", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
	}
	if strings.Contains(response.Body.String(), "Redis unavailable") {
		t.Fatal("internal store error leaked in the HTTP response")
	}
}
