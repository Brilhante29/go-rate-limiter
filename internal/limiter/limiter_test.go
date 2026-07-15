package limiter

import (
	"context"
	"errors"
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewPolicyRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name  string
		rate  float64
		burst int64
	}{
		{name: "zero rate", rate: 0, burst: 1},
		{name: "NaN rate", rate: math.NaN(), burst: 1},
		{name: "infinite rate", rate: math.Inf(1), burst: 1},
		{name: "zero burst", rate: 1, burst: 0},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := NewPolicy(test.rate, test.burst); !errors.Is(err, ErrInvalidPolicy) {
				t.Fatalf("expected invalid policy, got %v", err)
			}
		})
	}
}

func TestServiceValidatesInputs(t *testing.T) {
	policy, err := NewPolicy(10, 5)
	if err != nil {
		t.Fatal(err)
	}
	service, err := NewService(NewMemoryStore(), policy)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := service.Allow(context.Background(), "bad key", 1); !errors.Is(err, ErrInvalidKey) {
		t.Fatalf("expected invalid key, got %v", err)
	}
	if _, err := service.Allow(context.Background(), "tenant-a", 6); !errors.Is(err, ErrInvalidCost) {
		t.Fatalf("expected invalid cost, got %v", err)
	}
}

func TestMemoryStoreRefillsDeterministically(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	store := newMemoryStore(func() time.Time { return now })
	policy, _ := NewPolicy(2, 2)
	service, _ := NewService(store, policy)
	ctx := context.Background()

	for i := 0; i < 2; i++ {
		decision, err := service.Allow(ctx, "tenant-a", 1)
		if err != nil || !decision.Allowed {
			t.Fatalf("request %d should be allowed: decision=%+v err=%v", i, decision, err)
		}
	}
	decision, err := service.Allow(ctx, "tenant-a", 1)
	if err != nil || decision.Allowed || decision.RetryAfter != 500*time.Millisecond {
		t.Fatalf("expected rejection with 500ms retry, decision=%+v err=%v", decision, err)
	}

	now = now.Add(500 * time.Millisecond)
	decision, err = service.Allow(ctx, "tenant-a", 1)
	if err != nil || !decision.Allowed {
		t.Fatalf("request should be allowed after refill: decision=%+v err=%v", decision, err)
	}
}

func TestMemoryStorePreservesBurstUnderConcurrency(t *testing.T) {
	fixed := time.Unix(1_700_000_000, 0)
	store := newMemoryStore(func() time.Time { return fixed })
	policy, _ := NewPolicy(1, 100)
	service, _ := NewService(store, policy)

	var accepted atomic.Int64
	var wg sync.WaitGroup
	for i := 0; i < 1_000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			decision, err := service.Allow(context.Background(), "shared", 1)
			if err != nil {
				t.Errorf("allow: %v", err)
				return
			}
			if decision.Allowed {
				accepted.Add(1)
			}
		}()
	}
	wg.Wait()
	if got := accepted.Load(); got != 100 {
		t.Fatalf("accepted=%d, want 100", got)
	}
}

func BenchmarkMemoryStoreTake(b *testing.B) {
	store := NewMemoryStore()
	policy, _ := NewPolicy(1_000_000_000, 1_000_000_000)
	service, _ := NewService(store, policy)
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := service.Allow(context.Background(), "benchmark", 1); err != nil {
				b.Fatal(err)
			}
		}
	})
}
