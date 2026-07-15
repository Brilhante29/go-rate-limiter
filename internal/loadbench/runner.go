package loadbench

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Config struct {
	Targets     []string
	Duration    time.Duration
	Concurrency int
	Key         string
	Rate        float64
	Burst       int64
	Command     string
}

type Result struct {
	Project              string             `json:"project"`
	Metric               string             `json:"metric"`
	Value                float64            `json:"value"`
	Unit                 string             `json:"unit"`
	Timestamp            string             `json:"timestamp"`
	Command              string             `json:"command"`
	Repeat               int                `json:"repeat"`
	Summary              map[string]float64 `json:"summary"`
	Environment          map[string]string  `json:"environment"`
	DurationSeconds      float64            `json:"duration_seconds"`
	Concurrency          int                `json:"concurrency"`
	TotalRequests        int64              `json:"total_requests"`
	Accepted             int64              `json:"accepted"`
	Rejected             int64              `json:"rejected"`
	Errors               int64              `json:"errors"`
	Nodes                []string           `json:"nodes"`
	GlobalLimitPreserved bool               `json:"global_limit_preserved"`
}

func Run(ctx context.Context, config Config) (Result, error) {
	if len(config.Targets) < 1 || config.Duration <= 0 || config.Concurrency < 1 || config.Rate <= 0 || config.Burst < 1 {
		return Result{}, errors.New("targets, positive duration/concurrency/rate, and burst are required")
	}

	transport := &http.Transport{
		MaxIdleConns:        config.Concurrency * len(config.Targets),
		MaxIdleConnsPerHost: config.Concurrency,
		IdleConnTimeout:     30 * time.Second,
	}
	client := &http.Client{Transport: transport, Timeout: 2 * time.Second}
	defer transport.CloseIdleConnections()

	benchmarkCtx, cancel := context.WithTimeout(ctx, config.Duration)
	defer cancel()

	var accepted atomic.Int64
	var rejected atomic.Int64
	var requestErrors atomic.Int64
	var targetIndex atomic.Uint64
	var samplesMu sync.Mutex
	latencies := make([]float64, 0, 100_000)
	nodes := make(map[string]struct{})
	payload := []byte(`{"cost":1}`)
	start := time.Now()

	var workers sync.WaitGroup
	for i := 0; i < config.Concurrency; i++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			for benchmarkCtx.Err() == nil {
				index := targetIndex.Add(1) - 1
				target := strings.TrimRight(config.Targets[index%uint64(len(config.Targets))], "/")
				endpoint := target + "/v1/limits/" + url.PathEscape(config.Key) + "/check"
				request, err := http.NewRequestWithContext(benchmarkCtx, http.MethodPost, endpoint, bytes.NewReader(payload))
				if err != nil {
					requestErrors.Add(1)
					return
				}
				request.Header.Set("Content-Type", "application/json")
				requestStart := time.Now()
				response, err := client.Do(request)
				latencyMS := float64(time.Since(requestStart).Microseconds()) / 1000
				if err != nil {
					if benchmarkCtx.Err() != nil {
						return
					}
					requestErrors.Add(1)
					continue
				}
				_, _ = io.Copy(io.Discard, response.Body)
				_ = response.Body.Close()

				samplesMu.Lock()
				latencies = append(latencies, latencyMS)
				if node := response.Header.Get("X-Rate-Limiter-Node"); node != "" {
					nodes[node] = struct{}{}
				}
				samplesMu.Unlock()

				switch response.StatusCode {
				case http.StatusOK:
					accepted.Add(1)
				case http.StatusTooManyRequests:
					rejected.Add(1)
				default:
					requestErrors.Add(1)
				}
			}
		}()
	}
	workers.Wait()
	elapsed := time.Since(start)

	acceptedCount := accepted.Load()
	rejectedCount := rejected.Load()
	errorCount := requestErrors.Load()
	total := acceptedCount + rejectedCount + errorCount
	if acceptedCount == 0 || rejectedCount == 0 || total == 0 {
		return Result{}, fmt.Errorf("benchmark needs accepted and rejected requests: accepted=%d rejected=%d errors=%d", acceptedCount, rejectedCount, errorCount)
	}

	samplesMu.Lock()
	sort.Float64s(latencies)
	nodeNames := make([]string, 0, len(nodes))
	for node := range nodes {
		nodeNames = append(nodeNames, node)
	}
	samplesMu.Unlock()
	sort.Strings(nodeNames)

	seconds := elapsed.Seconds()
	acceptedRPS := float64(acceptedCount) / seconds
	rejectedRPS := float64(rejectedCount) / seconds
	totalRPS := float64(total) / seconds
	p95 := percentile(latencies, 0.95)
	maxAccepted := math.Ceil(float64(config.Burst)+config.Rate*seconds) + 1
	globalLimitPreserved := float64(acceptedCount) <= maxAccepted

	result := Result{
		Project:   "go-rate-limiter",
		Metric:    "total_rps",
		Value:     round(totalRPS, 2),
		Unit:      "requests_per_second",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Command:   config.Command,
		Repeat:    1,
		Summary: map[string]float64{
			"accepted_rps":          round(acceptedRPS, 2),
			"rejected_rps":          round(rejectedRPS, 2),
			"total_rps":             round(totalRPS, 2),
			"p95_latency_ms":        round(p95, 3),
			"max_expected_accepted": maxAccepted,
			"nodes_observed":        float64(len(nodeNames)),
		},
		Environment: map[string]string{
			"go":    runtime.Version(),
			"os":    runtime.GOOS,
			"arch":  runtime.GOARCH,
			"cpus":  strconv.Itoa(runtime.NumCPU()),
			"store": "redis",
		},
		DurationSeconds:      round(seconds, 3),
		Concurrency:          config.Concurrency,
		TotalRequests:        total,
		Accepted:             acceptedCount,
		Rejected:             rejectedCount,
		Errors:               errorCount,
		Nodes:                nodeNames,
		GlobalLimitPreserved: globalLimitPreserved,
	}
	if len(nodeNames) < len(config.Targets) {
		return result, fmt.Errorf("observed %d nodes for %d targets", len(nodeNames), len(config.Targets))
	}
	if errorCount != 0 {
		return result, fmt.Errorf("benchmark recorded %d unexpected errors", errorCount)
	}
	if !globalLimitPreserved {
		return result, fmt.Errorf("global limit violated: accepted=%d maximum=%.0f", acceptedCount, maxAccepted)
	}
	return result, nil
}

func Write(path string, result Result) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create result directory: %w", err)
	}
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal result: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write result: %w", err)
	}
	return nil
}

func JSON(result Result) ([]byte, error) {
	return json.MarshalIndent(result, "", "  ")
}

func percentile(sorted []float64, quantile float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	index := int(math.Ceil(quantile*float64(len(sorted)))) - 1
	if index < 0 {
		index = 0
	}
	if index >= len(sorted) {
		index = len(sorted) - 1
	}
	return sorted[index]
}

func round(value float64, places int) float64 {
	factor := math.Pow10(places)
	return math.Round(value*factor) / factor
}
