local key= KEYS[1]
local now=tonumber(ARGV[1])
local windowstart=tonumber(ARGV[2])
local limit=tonumber(ARGV[3])
local windowSeconds=tonumber(ARGV[4])
local member =ARGV[5]
--delete timestamps before windowstart
redis.call("ZREMRANGEBYSCORE", key, "-inf", windowstart)

local currentCount=redis.call("ZCARD",key)
if currentCount<limit then
    redis.call("ZADD",key,now,member)
    redis.call("EXPIRE",key,windowSeconds)
    return 1
end
if currentCount>=limit then
    return 0 
end

