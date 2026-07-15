package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/Brilhante29/go-rate-limiter/internal/adapters/redisstore"
	"github.com/Brilhante29/go-rate-limiter/internal/demo"
	"github.com/Brilhante29/go-rate-limiter/internal/httpapi"
	"github.com/Brilhante29/go-rate-limiter/internal/limiter"
	"github.com/Brilhante29/go-rate-limiter/internal/loadbench"
)

const dockerBenchmarkCommand = "docker run --rm go-rate-limiter"

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return errors.New("expected one of: serve, benchmark, demo")
	}
	switch args[0] {
	case "serve":
		return runServe(args[1:])
	case "benchmark":
		return runBenchmark(args[1:])
	case "demo":
		return runDemo(args[1:])
	case "help", "-h", "--help":
		printUsage()
		return nil
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func runServe(args []string) error {
	flags := flag.NewFlagSet("serve", flag.ContinueOnError)
	listen := flags.String("listen", envString("LISTEN_ADDR", ":8080"), "HTTP listen address")
	nodeID := flags.String("node", envString("NODE_ID", hostname()), "node identifier")
	storeKind := flags.String("store", envString("STORE", "redis"), "redis or memory")
	redisAddr := flags.String("redis-addr", envString("REDIS_ADDR", "127.0.0.1:6379"), "Redis address")
	prefix := flags.String("prefix", envString("REDIS_PREFIX", "rate-limiter"), "Redis key prefix")
	rate := flags.Float64("rate", envFloat("RATE_PER_SECOND", 1_000), "tokens added per second")
	burst := flags.Int64("burst", envInt64("BURST", 1_000), "maximum bucket size")
	if err := flags.Parse(args); err != nil {
		return err
	}

	policy, err := limiter.NewPolicy(*rate, *burst)
	if err != nil {
		return err
	}
	var store limiter.BucketStore
	var health httpapi.HealthFunc
	closeStore := func() error { return nil }
	storeName := strings.ToLower(*storeKind)
	switch storeName {
	case "redis":
		redisStore := redisstore.New(*redisAddr, *prefix)
		store = redisStore
		health = redisStore.Ping
		closeStore = redisStore.Close
		pingCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := redisStore.Ping(pingCtx); err != nil {
			_ = redisStore.Close()
			return fmt.Errorf("connect to Redis: %w", err)
		}
	case "memory":
		store = limiter.NewMemoryStore()
	default:
		return fmt.Errorf("unsupported store %q", *storeKind)
	}
	defer closeStore()

	service, err := limiter.NewService(store, policy)
	if err != nil {
		return err
	}
	server := &http.Server{
		Addr:              *listen,
		Handler:           httpapi.New(service, health, *nodeID, storeName).Handler(),
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       3 * time.Second,
		WriteTimeout:      3 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	serverErrors := make(chan error, 1)
	go func() {
		slog.Info("rate limiter listening", "address", *listen, "node", *nodeID, "store", storeName, "rate", *rate, "burst", *burst)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}

func runBenchmark(args []string) error {
	flags := flag.NewFlagSet("benchmark", flag.ContinueOnError)
	targetsRaw := flags.String("targets", "http://127.0.0.1:8081,http://127.0.0.1:8082", "comma-separated node URLs")
	duration := flags.Duration("duration", 10*time.Second, "benchmark duration")
	concurrency := flags.Int("concurrency", 64, "parallel workers")
	key := flags.String("key", "benchmark-shared", "shared rate-limit key")
	rate := flags.Float64("rate", 1_000, "configured token rate")
	burst := flags.Int64("burst", 1_000, "configured burst")
	output := flags.String("output", "benchmarks/results/rate-limiter-baseline.json", "result JSON path")
	command := flags.String("command", dockerBenchmarkCommand, "reproduction command recorded in JSON")
	if err := flags.Parse(args); err != nil {
		return err
	}
	result, benchmarkErr := loadbench.Run(context.Background(), loadbench.Config{
		Targets:     splitTargets(*targetsRaw),
		Duration:    *duration,
		Concurrency: *concurrency,
		Key:         *key,
		Rate:        *rate,
		Burst:       *burst,
		Command:     *command,
	})
	if result.Project != "" {
		if err := loadbench.Write(*output, result); err != nil {
			return err
		}
		if err := printResult(result); err != nil {
			return err
		}
	}
	return benchmarkErr
}

func runDemo(args []string) error {
	flags := flag.NewFlagSet("demo", flag.ContinueOnError)
	redisBinary := flags.String("redis-binary", "redis-server", "redis-server executable")
	output := flags.String("output", "/tmp/rate-limiter-baseline.json", "result JSON path")
	duration := flags.Duration("duration", 10*time.Second, "benchmark duration")
	concurrency := flags.Int("concurrency", 64, "parallel workers")
	rate := flags.Float64("rate", 1_000, "global token rate")
	burst := flags.Int64("burst", 1_000, "global burst")
	if err := flags.Parse(args); err != nil {
		return err
	}
	result, err := demo.Run(context.Background(), demo.Config{
		RedisBinary: *redisBinary,
		Output:      *output,
		Duration:    *duration,
		Concurrency: *concurrency,
		Rate:        *rate,
		Burst:       *burst,
		Command:     dockerBenchmarkCommand,
	})
	if result.Project != "" {
		if printErr := printResult(result); printErr != nil {
			return printErr
		}
	}
	return err
}

func splitTargets(raw string) []string {
	parts := strings.Split(raw, ",")
	targets := make([]string, 0, len(parts))
	for _, part := range parts {
		if target := strings.TrimSpace(part); target != "" {
			targets = append(targets, target)
		}
	}
	return targets
}

func printResult(result loadbench.Result) error {
	data, err := loadbench.JSON(result)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func envString(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

func envFloat(name string, fallback float64) float64 {
	value, err := strconv.ParseFloat(os.Getenv(name), 64)
	if err != nil {
		return fallback
	}
	return value
}

func envInt64(name string, fallback int64) int64 {
	value, err := strconv.ParseInt(os.Getenv(name), 10, 64)
	if err != nil {
		return fallback
	}
	return value
}

func hostname() string {
	value, err := os.Hostname()
	if err != nil || value == "" {
		return "node-unknown"
	}
	return value
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "usage: go-rate-limiter <serve|benchmark|demo> [flags]")
}
