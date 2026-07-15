package limiter

import (
	"context"
	"errors"
	"fmt"
	"math"
	"regexp"
	"time"
)

var keyPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$`)

var (
	ErrInvalidKey    = errors.New("key must match [A-Za-z0-9][A-Za-z0-9._:-]{0,127}")
	ErrInvalidCost   = errors.New("cost must be greater than zero and no greater than burst")
	ErrInvalidPolicy = errors.New("rate must be greater than zero and burst must be greater than zero")
)

type Policy struct {
	RatePerSecond float64
	Burst         int64
}

func NewPolicy(ratePerSecond float64, burst int64) (Policy, error) {
	if ratePerSecond <= 0 || math.IsNaN(ratePerSecond) || math.IsInf(ratePerSecond, 0) || burst <= 0 {
		return Policy{}, ErrInvalidPolicy
	}
	return Policy{RatePerSecond: ratePerSecond, Burst: burst}, nil
}

type Decision struct {
	Allowed    bool
	Remaining  int64
	RetryAfter time.Duration
}

type BucketStore interface {
	Take(ctx context.Context, key string, policy Policy, cost int64) (Decision, error)
}

type Service struct {
	store  BucketStore
	policy Policy
}

func NewService(store BucketStore, policy Policy) (*Service, error) {
	if store == nil {
		return nil, errors.New("bucket store is required")
	}
	if _, err := NewPolicy(policy.RatePerSecond, policy.Burst); err != nil {
		return nil, err
	}
	return &Service{store: store, policy: policy}, nil
}

func (s *Service) Allow(ctx context.Context, key string, cost int64) (Decision, error) {
	if !keyPattern.MatchString(key) {
		return Decision{}, ErrInvalidKey
	}
	if cost <= 0 || cost > s.policy.Burst {
		return Decision{}, ErrInvalidCost
	}
	decision, err := s.store.Take(ctx, key, s.policy, cost)
	if err != nil {
		return Decision{}, fmt.Errorf("take tokens: %w", err)
	}
	return decision, nil
}

func (s *Service) Policy() Policy {
	return s.policy
}
