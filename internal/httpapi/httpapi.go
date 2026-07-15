package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/Brilhante29/go-rate-limiter/internal/limiter"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type HealthFunc func(context.Context) error

type API struct {
	service   *limiter.Service
	health    HealthFunc
	nodeID    string
	storeName string
}

func New(service *limiter.Service, health HealthFunc, nodeID, storeName string) *API {
	return &API{service: service, health: health, nodeID: nodeID, storeName: storeName}
}

func (a *API) Handler() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Get("/healthz", a.healthz)
	router.Post("/v1/limits/{key}/check", a.check)
	return router
}

func (a *API) healthz(w http.ResponseWriter, r *http.Request) {
	if a.health != nil {
		ctx, cancel := context.WithTimeout(r.Context(), 500*time.Millisecond)
		defer cancel()
		if err := a.health(ctx); err != nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]any{
				"status": "unavailable",
				"node":   a.nodeID,
				"store":  a.storeName,
			})
			return
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"status": "ok",
		"node":   a.nodeID,
		"store":  a.storeName,
	})
}

type checkRequest struct {
	Cost int64 `json:"cost"`
}

type checkResponse struct {
	Allowed       bool    `json:"allowed"`
	Remaining     int64   `json:"remaining"`
	RetryAfterMS  int64   `json:"retry_after_ms"`
	RatePerSecond float64 `json:"rate_per_second"`
	Burst         int64   `json:"burst"`
	Node          string  `json:"node"`
}

func (a *API) check(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Rate-Limiter-Node", a.nodeID)

	request := checkRequest{Cost: 1}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<10)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil && !errors.Is(err, io.EOF) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}
	if err := ensureSingleJSONValue(decoder); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "body must contain one JSON value"})
		return
	}

	decision, err := a.service.Allow(r.Context(), chi.URLParam(r, "key"), request.Cost)
	if err != nil {
		status := http.StatusServiceUnavailable
		message := "rate limiter unavailable"
		if errors.Is(err, limiter.ErrInvalidKey) || errors.Is(err, limiter.ErrInvalidCost) {
			status = http.StatusBadRequest
			message = err.Error()
		}
		writeJSON(w, status, map[string]string{"error": message})
		return
	}

	policy := a.service.Policy()
	w.Header().Set("X-RateLimit-Limit", strconv.FormatInt(policy.Burst, 10))
	w.Header().Set("X-RateLimit-Remaining", strconv.FormatInt(decision.Remaining, 10))
	status := http.StatusOK
	if !decision.Allowed {
		status = http.StatusTooManyRequests
		retrySeconds := int64(math.Ceil(decision.RetryAfter.Seconds()))
		if retrySeconds < 1 {
			retrySeconds = 1
		}
		w.Header().Set("Retry-After", strconv.FormatInt(retrySeconds, 10))
	}
	writeJSON(w, status, checkResponse{
		Allowed:       decision.Allowed,
		Remaining:     decision.Remaining,
		RetryAfterMS:  decision.RetryAfter.Milliseconds(),
		RatePerSecond: policy.RatePerSecond,
		Burst:         policy.Burst,
		Node:          a.nodeID,
	})
}

func ensureSingleJSONValue(decoder *json.Decoder) error {
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err == nil {
			return errors.New("multiple JSON values")
		}
		return err
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
