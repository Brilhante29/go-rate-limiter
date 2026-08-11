package redisstore

import (
	"context"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Brilhante29/go-rate-limiter/internal/limiter"
)

func TestAsInt64(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  int64
	}{
		{name: "integer", input: int64(7), want: 7},
		{name: "string", input: "8", want: 8},
		{name: "bytes", input: []byte("9"), want: 9},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := asInt64(test.input)
			if err != nil || got != test.want {
				t.Fatalf("got=%d err=%v, want=%d", got, err, test.want)
			}
		})
	}
	if _, err := asInt64(1.5); err == nil {
		t.Fatal("expected unsupported type error")
	}
}

func TestRedisStoresPreserveOneBurstUnderConcurrency(t *testing.T) {
	addr := os.Getenv("REDIS_TEST_ADDR")
	if addr == "" {
		t.Skip("REDIS_TEST_ADDR is required for the Redis integration test")
	}

	prefix := "rate-limiter-test-" + time.Now().UTC().Format("20060102150405.000000000")
	storeA := New(addr, prefix)
	storeB := New(addr, prefix)
	t.Cleanup(func() {
		_ = storeA.Close()
		_ = storeB.Close()
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := storeA.Ping(ctx); err != nil {
		t.Fatalf("ping Redis: %v", err)
	}
	if err := storeA.Reset(ctx, "shared"); err != nil {
		t.Fatalf("reset bucket: %v", err)
	}

	policy, err := limiter.NewPolicy(0.001, 100)
	if err != nil {
		t.Fatal(err)
	}
	stores := []*Store{storeA, storeB}
	var accepted atomic.Int64
	var failures atomic.Int64
	var workers sync.WaitGroup
	for index := 0; index < 1_000; index++ {
		workers.Add(1)
		go func(store *Store) {
			defer workers.Done()
			decision, takeErr := store.Take(ctx, "shared", policy, 1)
			if takeErr != nil {
				failures.Add(1)
				return
			}
			if decision.Allowed {
				accepted.Add(1)
			}
		}(stores[index%len(stores)])
	}
	workers.Wait()

	if got := failures.Load(); got != 0 {
		t.Fatalf("Redis calls failed=%d", got)
	}
	if got := accepted.Load(); got != 100 {
		t.Fatalf("accepted=%d, want exactly one shared burst of 100", got)
	}
}

func TestRedisStoreFailsClosedWhenUnavailable(t *testing.T) {
	store := New("127.0.0.1:1", "unavailable-test")
	t.Cleanup(func() { _ = store.Close() })
	policy, err := limiter.NewPolicy(10, 10)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()
	if _, err := store.Take(ctx, "shared", policy, 1); err == nil {
		t.Fatal("expected an error when Redis is unavailable")
	}
}
