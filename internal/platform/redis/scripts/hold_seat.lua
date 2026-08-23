-- KEYS[1]: The seat hold key  -> "hold:{eventId}:{seatNo}"
-- KEYS[2]: The user hold key  -> "user:hold:{userId}"
-- ARGV[1]: The user ID        -> "6"

local seatHoldKey = KEYS[1]
local userHoldKey = KEYS[2]
local userId = ARGV[1]

-- 1. Check if the seat is already held
if redis.call("EXISTS", seatHoldKey) == 1 then
    return 0 -- Code 0: Seat is already held by someone else
end

-- 2. Check if this user already holds a seat (Anti-Hoarding)
if redis.call("EXISTS", userHoldKey) == 1 then
    return -1 -- Code -1: User already has an active hold
end

-- 3. Both checks passed! Set the hold with a 600-second (10-minute) TTL
redis.call("SET", seatHoldKey, userId, "EX", 600)
redis.call("SET", userHoldKey, seatHoldKey, "EX", 600)

return 1 -- Code 1: Success! Seat successfully held