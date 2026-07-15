package limiter

import (
	"context"
	"math"
	"sync"
	"time"
)

type memoryBucket struct {
	tokens float64
	last   time.Time
}

type MemoryStore struct {
	mu      sync.Mutex
	buckets map[string]memoryBucket
	now     func() time.Time
}

func NewMemoryStore() *MemoryStore {
	return newMemoryStore(time.Now)
}

func newMemoryStore(now func() time.Time) *MemoryStore {
	return &MemoryStore{buckets: make(map[string]memoryBucket), now: now}
}

func (s *MemoryStore) Take(_ context.Context, key string, policy Policy, cost int64) (Decision, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := s.now()
	bucket, found := s.buckets[key]
	if !found {
		bucket = memoryBucket{tokens: float64(policy.Burst), last: now}
	}

	elapsed := now.Sub(bucket.last).Seconds()
	if elapsed < 0 {
		elapsed = 0
	}
	bucket.tokens = math.Min(float64(policy.Burst), bucket.tokens+elapsed*policy.RatePerSecond)
	bucket.last = now

	decision := Decision{}
	if bucket.tokens >= float64(cost) {
		bucket.tokens -= float64(cost)
		decision.Allowed = true
	} else {
		seconds := (float64(cost) - bucket.tokens) / policy.RatePerSecond
		decision.RetryAfter = time.Duration(math.Ceil(seconds * float64(time.Second)))
	}
	decision.Remaining = int64(math.Floor(bucket.tokens))
	s.buckets[key] = bucket
	return decision, nil
}
