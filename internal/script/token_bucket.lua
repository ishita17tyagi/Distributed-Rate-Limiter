local key = KEYS[1]
local capacity = tonumber(ARGV[1])
local refill_rate = tonumber(ARGV[2])
local now_ms = tonumber(ARGV[3])

-- Get the bucket data
local bucket = redis.call('HMGET', key, 'tokens', 'last_refill_ms')
local tokens = tonumber(bucket[1])
local last_refill_ms = tonumber(bucket[2])

-- INITIALIZATION: If key is missing (nil), initialize it
if tokens == nil then
    tokens = capacity
    last_refill_ms = now_ms
end

-- Math
local elapsed_sec = math.max(0, (now_ms - last_refill_ms) / 1000.0)
tokens = math.min(capacity, tokens + (elapsed_sec * refill_rate))

local allowed = 0
local retry_after = 0

if tokens >= 1 then
    tokens = tokens - 1
    allowed = 1
else
    local missing = 1 - tokens
    retry_after = math.ceil(missing / refill_rate)
end

-- Save
redis.call('HSET', key, 'tokens', tokens, 'last_refill_ms', now_ms)
redis.call('EXPIRE', key, math.ceil(capacity / refill_rate) * 2)

return {allowed, math.floor(tokens), retry_after}