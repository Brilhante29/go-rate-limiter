package loadbench

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
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

const workloadVersion = "token-bucket-http-v2"

type Config struct {
	Targets              []string
	Duration             time.Duration
	WarmupDuration       time.Duration
	Concurrency          int
	WarmupIterations     int
	MeasuredIterations   int
	Key                  string
	Rate                 float64
	Burst                int64
	Command              string
	FixtureDigest        string
	SourceCommit         string
	CleanTree            bool
	ImageRef             string
	ImageDigest          string
	DependencyLockDigest string
	Producer             string
	ArtifactDigest       string
	HardwareClass        string
	RedisVersion         string
}

type Metric struct {
	Name       string             `json:"name"`
	Value      float64            `json:"value"`
	Unit       string             `json:"unit"`
	Direction  string             `json:"direction"`
	Samples    []float64          `json:"samples"`
	Failures   int                `json:"failures"`
	Summary    map[string]float64 `json:"summary"`
}

type Workload struct {
	Version            string `json:"version"`
	FixtureDigest      string `json:"fixture_digest"`
	ConfigDigest       string `json:"config_digest"`
	WarmupIterations   int    `json:"warmup_iterations"`
	MeasuredIterations int    `json:"measured_iterations"`
	Concurrency        int    `json:"concurrency"`
}

type Execution struct {
	Command         string  `json:"command"`
	StartedAt       string  `json:"started_at"`
	DurationSeconds float64 `json:"duration_seconds"`
	ExitCode        int     `json:"exit_code"`
	Repeat          int     `json:"repeat"`
}

type Environment struct {
	Runtime       string `json:"runtime"`
	Architecture  string `json:"architecture"`
	HardwareClass string `json:"hardware_class"`
	OS            string `json:"os"`
	CPUCount      string `json:"cpu_count"`
	Store         string `json:"store"`
	RedisVersion  string `json:"redis_version"`
	HTTPNodes     string `json:"http_nodes"`
	Topology      string `json:"topology"`
}

type Provenance struct {
	SourceCommit         string `json:"source_commit"`
	CleanTree            bool   `json:"clean_tree"`
	ImageRef             string `json:"image_ref"`
	ImageDigest          string `json:"image_digest"`
	DependencyLockDigest string `json:"dependency_lock_digest"`
	Producer             string `json:"producer"`
	ArtifactDigest       string `json:"artifact_digest"`
}

type Result struct {
	SchemaVersion    int         `json:"schema_version"`
	RunID            string      `json:"run_id"`
	Project          string      `json:"project"`
	BenchmarkID      string      `json:"benchmark_id"`
	Workload         Workload    `json:"workload"`
	Metrics          []Metric    `json:"metrics"`
	Execution        Execution   `json:"execution"`
	Environment      Environment `json:"environment"`
	Provenance       Provenance  `json:"provenance"`
	ComparabilityKey string      `json:"comparability_key"`
}

type iterationResult struct {
	totalRPS             float64
	acceptedRPS          float64
	rejectedRPS          float64
	p95LatencyMS         float64
	nodesObserved        float64
	unexpectedErrors     float64
	globalLimitViolation float64
}

func Run(ctx context.Context, config Config) (Result, error) {
	if err := validateConfig(config); err != nil {
		return Result{}, err
	}

	startedAt := time.Now().UTC()
	transport := &http.Transport{
		MaxIdleConns:        config.Concurrency * len(config.Targets),
		MaxIdleConnsPerHost: config.Concurrency,
		IdleConnTimeout:     30 * time.Second,
	}
	client := &http.Client{Transport: transport, Timeout: 2 * time.Second}
	defer transport.CloseIdleConnections()

	for iteration := 0; iteration < config.WarmupIterations; iteration++ {
		key := fmt.Sprintf("%s-warmup-%d", config.Key, iteration+1)
		if _, err := runIteration(ctx, client, config, config.WarmupDuration, key); err != nil {
			return Result{}, fmt.Errorf("warmup %d: %w", iteration+1, err)
		}
	}

	samples := make([]iterationResult, 0, config.MeasuredIterations)
	failures := 0
	var benchmarkErrors []error
	for iteration := 0; iteration < config.MeasuredIterations; iteration++ {
		key := fmt.Sprintf("%s-measured-%d", config.Key, iteration+1)
		sample, err := runIteration(ctx, client, config, config.Duration, key)
		samples = append(samples, sample)
		if err != nil {
			failures++
			benchmarkErrors = append(benchmarkErrors, fmt.Errorf("measured iteration %d: %w", iteration+1, err))
		}
	}

	result, err := aggregate(config, startedAt, time.Since(startedAt), samples, failures)
	if err != nil {
		return Result{}, err
	}
	if len(benchmarkErrors) > 0 {
		return result, errors.Join(benchmarkErrors...)
	}
	return result, nil
}

func validateConfig(config Config) error {
	if len(config.Targets) < 2 {
		return errors.New("at least two HTTP targets are required")
	}
	if config.Duration <= 0 || config.WarmupDuration <= 0 || config.Concurrency < 1 || config.Rate <= 0 || config.Burst < 1 {
		return errors.New("positive duration, warmup duration, concurrency, rate, and burst are required")
	}
	if config.WarmupIterations < 1 {
		return errors.New("at least one warmup iteration is required")
	}
	if config.MeasuredIterations < 3 {
		return errors.New("at least three measured iterations are required")
	}
	return nil
}

func runIteration(ctx context.Context, client *http.Client, config Config, duration time.Duration, key string) (iterationResult, error) {
	benchmarkCtx, cancel := context.WithTimeout(ctx, duration)
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
				endpoint := target + "/v1/limits/" + url.PathEscape(key) + "/check"
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
	decisions := acceptedCount + rejectedCount
	seconds := elapsed.Seconds()
	maxAccepted := math.Ceil(float64(config.Burst)+config.Rate*seconds) + 1
	globalLimitViolation := float64(0)
	if float64(acceptedCount) > maxAccepted {
		globalLimitViolation = 1
	}

	samplesMu.Lock()
	sort.Float64s(latencies)
	nodesObserved := len(nodes)
	samplesMu.Unlock()

	result := iterationResult{
		totalRPS:             round(float64(decisions)/seconds, 2),
		acceptedRPS:          round(float64(acceptedCount)/seconds, 2),
		rejectedRPS:          round(float64(rejectedCount)/seconds, 2),
		p95LatencyMS:         round(percentile(latencies, 0.95), 3),
		nodesObserved:        float64(nodesObserved),
		unexpectedErrors:     float64(errorCount),
		globalLimitViolation: globalLimitViolation,
	}

	var iterationErrors []error
	if acceptedCount == 0 || rejectedCount == 0 || decisions == 0 {
		iterationErrors = append(iterationErrors, fmt.Errorf("accepted=%d rejected=%d errors=%d", acceptedCount, rejectedCount, errorCount))
	}
	if nodesObserved < len(config.Targets) {
		iterationErrors = append(iterationErrors, fmt.Errorf("observed %d nodes for %d targets", nodesObserved, len(config.Targets)))
	}
	if errorCount != 0 {
		iterationErrors = append(iterationErrors, fmt.Errorf("recorded %d unexpected errors", errorCount))
	}
	if globalLimitViolation != 0 {
		iterationErrors = append(iterationErrors, fmt.Errorf("global limit violated: accepted=%d maximum=%.0f", acceptedCount, maxAccepted))
	}
	return result, errors.Join(iterationErrors...)
}

func aggregate(config Config, startedAt time.Time, elapsed time.Duration, iterations []iterationResult, failures int) (Result, error) {
	runID, err := newUUID()
	if err != nil {
		return Result{}, fmt.Errorf("create run id: %w", err)
	}
	configDigest, err := digestConfig(config)
	if err != nil {
		return Result{}, fmt.Errorf("digest benchmark config: %w", err)
	}

	totalRPS := collect(iterations, func(sample iterationResult) float64 { return sample.totalRPS })
	acceptedRPS := collect(iterations, func(sample iterationResult) float64 { return sample.acceptedRPS })
	rejectedRPS := collect(iterations, func(sample iterationResult) float64 { return sample.rejectedRPS })
	p95Latency := collect(iterations, func(sample iterationResult) float64 { return sample.p95LatencyMS })
	nodesObserved := collect(iterations, func(sample iterationResult) float64 { return sample.nodesObserved })
	unexpectedErrors := collect(iterations, func(sample iterationResult) float64 { return sample.unexpectedErrors })
	globalLimitViolations := collect(iterations, func(sample iterationResult) float64 { return sample.globalLimitViolation })

	hardwareClass := config.HardwareClass
	if hardwareClass == "" {
		hardwareClass = fmt.Sprintf("local-docker-%dcpu", runtime.NumCPU())
	}
	redisVersion := config.RedisVersion
	if redisVersion == "" {
		redisVersion = "8.8.0"
	}
	producer := config.Producer
	if producer == "" {
		producer = "local"
	}

	return Result{
		SchemaVersion: 2,
		RunID:         runID,
		Project:       "go-rate-limiter",
		BenchmarkID:   "shared-token-bucket",
		Workload: Workload{
			Version:            workloadVersion,
			FixtureDigest:      digestOrPlaceholder(config.FixtureDigest, workloadVersion),
			ConfigDigest:       configDigest,
			WarmupIterations:   config.WarmupIterations,
			MeasuredIterations: config.MeasuredIterations,
			Concurrency:        config.Concurrency,
		},
		Metrics: []Metric{
			newMetric("total_rps", "requests_per_second", "higher_is_better", totalRPS, failures),
			newMetric("accepted_rps", "requests_per_second", "target", acceptedRPS, failures),
			newMetric("rejected_rps", "requests_per_second", "target", rejectedRPS, failures),
			newMetric("p95_latency_ms", "milliseconds", "lower_is_better", p95Latency, failures),
			newMetric("nodes_observed", "count", "target", nodesObserved, countFailures(nodesObserved, func(value float64) bool { return value < float64(len(config.Targets)) })),
			newMetric("unexpected_errors", "count", "target", unexpectedErrors, countFailures(unexpectedErrors, func(value float64) bool { return value != 0 })),
			newMetric("global_limit_violations", "count", "target", globalLimitViolations, countFailures(globalLimitViolations, func(value float64) bool { return value != 0 })),
		},
		Execution: Execution{
			Command:         config.Command,
			StartedAt:       startedAt.Format(time.RFC3339),
			DurationSeconds: round(elapsed.Seconds(), 3),
			ExitCode:        0,
			Repeat:          config.MeasuredIterations,
		},
		Environment: Environment{
			Runtime:       runtime.Version(),
			Architecture:  runtime.GOARCH,
			HardwareClass: hardwareClass,
			OS:            runtime.GOOS,
			CPUCount:      strconv.Itoa(runtime.NumCPU()),
			Store:         "redis",
			RedisVersion:  redisVersion,
			HTTPNodes:     strconv.Itoa(len(config.Targets)),
			Topology:      "two-http-processes-one-redis-process",
		},
		Provenance: Provenance{
			SourceCommit:         valueOrPlaceholder(config.SourceCommit, 40),
			CleanTree:            config.CleanTree,
			ImageRef:             stringOrDefault(config.ImageRef, "go-rate-limiter:local"),
			ImageDigest:          digestOrPlaceholder(config.ImageDigest, "image"),
			DependencyLockDigest: digestOrPlaceholder(config.DependencyLockDigest, "go.sum"),
			Producer:             producer,
			ArtifactDigest:       digestOrPlaceholder(config.ArtifactDigest, "go-rate-limiter-binary"),
		},
		ComparabilityKey: comparabilityKey(config, redisVersion),
	}, nil
}

func newMetric(name, unit, direction string, samples []float64, failures int) Metric {
	return Metric{
		Name:      name,
		Value:     median(samples),
		Unit:      unit,
		Direction: direction,
		Samples:   samples,
		Failures:  failures,
		Summary:   summarize(samples),
	}
}

func collect(iterations []iterationResult, value func(iterationResult) float64) []float64 {
	result := make([]float64, 0, len(iterations))
	for _, iteration := range iterations {
		result = append(result, value(iteration))
	}
	return result
}

func countFailures(samples []float64, failed func(float64) bool) int {
	count := 0
	for _, sample := range samples {
		if failed(sample) {
			count++
		}
	}
	return count
}

func summarize(samples []float64) map[string]float64 {
	if len(samples) == 0 {
		return map[string]float64{"min": 0, "median": 0, "mean": 0, "max": 0}
	}
	ordered := append([]float64(nil), samples...)
	sort.Float64s(ordered)
	total := float64(0)
	for _, sample := range ordered {
		total += sample
	}
	return map[string]float64{
		"min":    ordered[0],
		"median": median(ordered),
		"mean":   round(total/float64(len(ordered)), 3),
		"max":    ordered[len(ordered)-1],
	}
}

func median(samples []float64) float64 {
	if len(samples) == 0 {
		return 0
	}
	ordered := append([]float64(nil), samples...)
	sort.Float64s(ordered)
	middle := len(ordered) / 2
	if len(ordered)%2 == 1 {
		return ordered[middle]
	}
	return round((ordered[middle-1]+ordered[middle])/2, 3)
}

func digestConfig(config Config) (string, error) {
	canonical := struct {
		Version              string  `json:"version"`
		Nodes                int     `json:"nodes"`
		DurationMilliseconds int64   `json:"duration_ms"`
		WarmupMilliseconds   int64   `json:"warmup_ms"`
		Concurrency          int     `json:"concurrency"`
		WarmupIterations     int     `json:"warmup_iterations"`
		MeasuredIterations   int     `json:"measured_iterations"`
		Rate                 float64 `json:"rate_per_second"`
		Burst                int64   `json:"burst"`
		Cost                 int64   `json:"cost"`
	}{
		Version:              workloadVersion,
		Nodes:                len(config.Targets),
		DurationMilliseconds: config.Duration.Milliseconds(),
		WarmupMilliseconds:   config.WarmupDuration.Milliseconds(),
		Concurrency:          config.Concurrency,
		WarmupIterations:     config.WarmupIterations,
		MeasuredIterations:   config.MeasuredIterations,
		Rate:                 config.Rate,
		Burst:                config.Burst,
		Cost:                 1,
	}
	data, err := json.Marshal(canonical)
	if err != nil {
		return "", err
	}
	return digestBytes(data), nil
}

func comparabilityKey(config Config, redisVersion string) string {
	rate := strings.ReplaceAll(strconv.FormatFloat(config.Rate, 'f', -1, 64), ".", "_")
	return fmt.Sprintf(
		"go-rate-limiter:%s:redis-%s:nodes-%d:rate-%s:burst-%d:duration-%dms:concurrency-%d",
		workloadVersion,
		strings.ReplaceAll(redisVersion, ".", "_"),
		len(config.Targets),
		rate,
		config.Burst,
		config.Duration.Milliseconds(),
		config.Concurrency,
	)
}

func digestOrPlaceholder(value, fallback string) string {
	if strings.HasPrefix(value, "sha256:") && len(value) == len("sha256:")+64 {
		return strings.ToLower(value)
	}
	return digestBytes([]byte(fallback))
}

func valueOrPlaceholder(value string, length int) string {
	if len(value) == length {
		return strings.ToLower(value)
	}
	return strings.Repeat("0", length)
}

func stringOrDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func digestBytes(data []byte) string {
	hash := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(hash[:])
}

func newUUID() (string, error) {
	data := make([]byte, 16)
	if _, err := rand.Read(data); err != nil {
		return "", err
	}
	data[6] = (data[6] & 0x0f) | 0x40
	data[8] = (data[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", data[0:4], data[4:6], data[6:8], data[8:10], data[10:16]), nil
}

func Write(path string, result Result) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create result directory: %w", err)
	}
	data, err := JSON(result)
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
