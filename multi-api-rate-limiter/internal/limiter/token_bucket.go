package limiter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLimiter struct {
	client *redis.Client
}

var lua = `
local key = KEYS[1]
local rate = tonumber(ARGV[1])
local capacity = tonumber(ARGV[2])
local now = tonumber(ARGV[3])

local bucket = redis.call("HMGET", key, "tokens", "last_refill")
local tokens = tonumber(bucket[1]) or capacity
local lastRefill = tonumber(bucket[2]) or now

-- Refill the bucket
local elapsed = now - lastRefill
local refillTokens = math.floor(elapsed * rate)
tokens = math.min(tokens + refillTokens, capacity)

if tokens > 0 then
	-- Consume a token
	tokens = tokens - 1
	redis.call("HMSET", key, "tokens", tokens, "last_refill", now)
	return 1 -- Allowed
else
	redis.call("HMSET", key, "tokens", tokens, "last_refill", now)
	return 0 -- Denied
end
`

func NewRedisLimiter(client *redis.Client) *RedisLimiter {
	return &RedisLimiter{client: client}
}

func (l *RedisLimiter) Allow(ctx context.Context, key string, rate float64, capacity int) (bool, error) {
	now := time.Now().Unix()

	result, err := l.client.Eval(ctx, lua, []string{key}, rate, capacity, now).Result()
	if err != nil {
		return false, err
	}

	return result.(int64) == 1, nil
}
