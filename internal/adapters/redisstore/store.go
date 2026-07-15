package redisstore

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Brilhante29/go-rate-limiter/internal/limiter"
	"github.com/redis/go-redis/v9"
)

const tokenBucketScript = `
local rate = tonumber(ARGV[1])
local burst = tonumber(ARGV[2])
local cost = tonumber(ARGV[3])
local server_time = redis.call('TIME')
local now_ms = tonumber(server_time[1]) * 1000 + math.floor(tonumber(server_time[2]) / 1000)
local state = redis.call('HMGET', KEYS[1], 'tokens', 'updated_at_ms')
local tokens = tonumber(state[1])
local updated_at_ms = tonumber(state[2])

if tokens == nil then
  tokens = burst
  updated_at_ms = now_ms
end

local elapsed_ms = math.max(0, now_ms - updated_at_ms)
tokens = math.min(burst, tokens + (elapsed_ms * rate / 1000))

local allowed = 0
local retry_after_ms = 0
if tokens >= cost then
  tokens = tokens - cost
  allowed = 1
else
  retry_after_ms = math.ceil(((cost - tokens) / rate) * 1000)
end

redis.call('HSET', KEYS[1], 'tokens', tokens, 'updated_at_ms', now_ms)
local ttl_ms = math.max(1000, math.ceil((burst / rate) * 2000))
redis.call('PEXPIRE', KEYS[1], ttl_ms)

return {allowed, math.floor(tokens), retry_after_ms}
`

type Store struct {
	client *redis.Client
	prefix string
	script *redis.Script
}

func New(addr, prefix string) *Store {
	return &Store{
		client: redis.NewClient(&redis.Options{
			Addr:         addr,
			DialTimeout:  2 * time.Second,
			ReadTimeout:  2 * time.Second,
			WriteTimeout: 2 * time.Second,
			PoolSize:     64,
		}),
		prefix: strings.TrimSuffix(prefix, ":") + ":",
		script: redis.NewScript(tokenBucketScript),
	}
}

func (s *Store) Take(ctx context.Context, key string, policy limiter.Policy, cost int64) (limiter.Decision, error) {
	result, err := s.script.Run(
		ctx,
		s.client,
		[]string{s.redisKey(key)},
		strconv.FormatFloat(policy.RatePerSecond, 'f', -1, 64),
		policy.Burst,
		cost,
	).Slice()
	if err != nil {
		return limiter.Decision{}, fmt.Errorf("execute token bucket script: %w", err)
	}
	if len(result) != 3 {
		return limiter.Decision{}, fmt.Errorf("unexpected script result length: %d", len(result))
	}

	allowed, err := asInt64(result[0])
	if err != nil {
		return limiter.Decision{}, fmt.Errorf("decode allowed: %w", err)
	}
	remaining, err := asInt64(result[1])
	if err != nil {
		return limiter.Decision{}, fmt.Errorf("decode remaining: %w", err)
	}
	retryAfterMS, err := asInt64(result[2])
	if err != nil {
		return limiter.Decision{}, fmt.Errorf("decode retry after: %w", err)
	}
	return limiter.Decision{
		Allowed:    allowed == 1,
		Remaining:  remaining,
		RetryAfter: time.Duration(retryAfterMS) * time.Millisecond,
	}, nil
}

func (s *Store) Ping(ctx context.Context) error {
	return s.client.Ping(ctx).Err()
}

func (s *Store) Reset(ctx context.Context, key string) error {
	return s.client.Del(ctx, s.redisKey(key)).Err()
}

func (s *Store) Close() error {
	return s.client.Close()
}

func (s *Store) redisKey(key string) string {
	return s.prefix + key
}

func asInt64(value any) (int64, error) {
	switch typed := value.(type) {
	case int64:
		return typed, nil
	case string:
		return strconv.ParseInt(typed, 10, 64)
	case []byte:
		return strconv.ParseInt(string(typed), 10, 64)
	default:
		return 0, fmt.Errorf("unsupported value type %T", value)
	}
}
