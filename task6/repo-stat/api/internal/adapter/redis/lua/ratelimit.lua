local key = KEYS[1]
local rps = tonumber(ARGV[1])
local burst = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local requested = 1

local res = redis.call("HMGET", key, "tokens", "last_time")
local last_tokens = tonumber(res[1]) or burst
local last_time = tonumber(res[2]) or now

local delta = math.max(0, now - last_time)
local refill = delta * rps
local current_tokens = math.min(burst, last_tokens + refill)

redis.call("EXPIRE", key, 60)

if current_tokens >= requested then
    current_tokens = current_tokens - requested
    redis.call("HMSET", key, "tokens", current_tokens, "last_time", now)
    return 1
else
    return 0
end
