package demo

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/Brilhante29/go-rate-limiter/internal/adapters/redisstore"
	"github.com/Brilhante29/go-rate-limiter/internal/httpapi"
	"github.com/Brilhante29/go-rate-limiter/internal/limiter"
	"github.com/Brilhante29/go-rate-limiter/internal/loadbench"
)

type Config struct {
	RedisBinary         string
	Output              string
	Duration            time.Duration
	WarmupDuration      time.Duration
	Concurrency         int
	WarmupIterations    int
	MeasuredIterations  int
	Rate                float64
	Burst               int64
	Command             string
	FixtureDigest       string
	SourceCommit        string
	CleanTree           bool
	ImageRef            string
	ImageDigest         string
	DependencyLockDigest string
	Producer            string
	ArtifactDigest      string
	HardwareClass       string
	RedisVersion        string
}

func Run(ctx context.Context, config Config) (loadbench.Result, error) {
	redisAddr, err := freeAddress()
	if err != nil {
		return loadbench.Result{}, fmt.Errorf("reserve Redis address: %w", err)
	}
	redisCommand := exec.Command(
		config.RedisBinary,
		"--bind", "127.0.0.1",
		"--port", portOf(redisAddr),
		"--save", "",
		"--appendonly", "no",
		"--protected-mode", "yes",
		"--loglevel", "warning",
	)
	redisCommand.Stdout = os.Stderr
	redisCommand.Stderr = os.Stderr
	if err := redisCommand.Start(); err != nil {
		return loadbench.Result{}, fmt.Errorf("start Redis: %w", err)
	}
	defer stopProcess(redisCommand)

	storeA := redisstore.New(redisAddr, "rate-limiter")
	defer storeA.Close()
	storeB := redisstore.New(redisAddr, "rate-limiter")
	defer storeB.Close()
	if err := waitForRedis(ctx, storeA); err != nil {
		return loadbench.Result{}, err
	}

	policy, err := limiter.NewPolicy(config.Rate, config.Burst)
	if err != nil {
		return loadbench.Result{}, err
	}
	serviceA, _ := limiter.NewService(storeA, policy)
	serviceB, _ := limiter.NewService(storeB, policy)

	listenerA, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return loadbench.Result{}, fmt.Errorf("listen node-a: %w", err)
	}
	listenerB, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		_ = listenerA.Close()
		return loadbench.Result{}, fmt.Errorf("listen node-b: %w", err)
	}

	serverA := newServer(httpapi.New(serviceA, storeA.Ping, "node-a", "redis").Handler())
	serverB := newServer(httpapi.New(serviceB, storeB.Ping, "node-b", "redis").Handler())
	defer shutdownServers(serverA, serverB)
	serverErrors := make(chan error, 2)
	go serve(serverA, listenerA, serverErrors)
	go serve(serverB, listenerB, serverErrors)

	const benchmarkKey = "benchmark-shared"
	if err := storeA.Reset(ctx, benchmarkKey); err != nil {
		return loadbench.Result{}, fmt.Errorf("reset benchmark bucket: %w", err)
	}
	result, benchmarkErr := loadbench.Run(ctx, loadbench.Config{
		Targets: []string{
			"http://" + listenerA.Addr().String(),
			"http://" + listenerB.Addr().String(),
		},
		Duration:             config.Duration,
		WarmupDuration:       config.WarmupDuration,
		Concurrency:          config.Concurrency,
		WarmupIterations:     config.WarmupIterations,
		MeasuredIterations:   config.MeasuredIterations,
		Key:                  benchmarkKey,
		Rate:                 config.Rate,
		Burst:                config.Burst,
		Command:              config.Command,
		FixtureDigest:        config.FixtureDigest,
		SourceCommit:         config.SourceCommit,
		CleanTree:            config.CleanTree,
		ImageRef:             config.ImageRef,
		ImageDigest:          config.ImageDigest,
		DependencyLockDigest: config.DependencyLockDigest,
		Producer:             config.Producer,
		ArtifactDigest:       config.ArtifactDigest,
		HardwareClass:        config.HardwareClass,
		RedisVersion:         config.RedisVersion,
	})

	select {
	case serverErr := <-serverErrors:
		if serverErr != nil && benchmarkErr == nil {
			benchmarkErr = serverErr
		}
	default:
	}
	if result.Project != "" {
		if err := loadbench.Write(config.Output, result); err != nil {
			return result, err
		}
	}
	return result, benchmarkErr
}

func newServer(handler http.Handler) *http.Server {
	return &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       3 * time.Second,
		WriteTimeout:      3 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
}

func serve(server *http.Server, listener net.Listener, errorsChannel chan<- error) {
	if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		errorsChannel <- err
	}
}

func shutdownServers(servers ...*http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	for _, server := range servers {
		_ = server.Shutdown(ctx)
	}
}

func freeAddress() (string, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", err
	}
	defer listener.Close()
	return listener.Addr().String(), nil
}

func portOf(address string) string {
	_, port, err := net.SplitHostPort(address)
	if err != nil {
		return "0"
	}
	return port
}

func waitForRedis(ctx context.Context, store *redisstore.Store) error {
	deadline := time.NewTimer(5 * time.Second)
	defer deadline.Stop()
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	for {
		pingCtx, cancel := context.WithTimeout(ctx, 250*time.Millisecond)
		err := store.Ping(pingCtx)
		cancel()
		if err == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-deadline.C:
			return fmt.Errorf("Redis did not become ready")
		case <-ticker.C:
		}
	}
}

func stopProcess(command *exec.Cmd) {
	if command.Process == nil {
		return
	}
	_ = command.Process.Kill()
	_ = command.Wait()
}
