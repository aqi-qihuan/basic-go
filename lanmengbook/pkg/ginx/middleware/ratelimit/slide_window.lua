-- 1, 2, 3, 4, 5, 6, 7 这是你的元素
-- ZREMRANGEBYSCORE key1 0 6
-- 7 执行完之后

-- 限流对象
local key = KEYS[1]
-- 窗口大小
local window = tonumber(ARGV[1])
-- 阈值
local threshold = tonumber( ARGV[2])
local now = tonumber(ARGV[3])
-- 窗口的起始时间
local min = now - window

-- 移除窗口之外的元素
redis.call('ZREMRANGEBYSCORE', key, '-inf', min)
-- 统计窗口内的元素数量
local cnt = redis.call('ZCOUNT', key, '-inf', '+inf')
-- local cnt = redis.call('ZCOUNT', key, min, '+inf')
-- 如果元素数量超过阈值，则执行限流
if cnt >= threshold then
    -- 执行限流
    return "true"
else
    -- 把 score 和 member 都设置成 now
    redis.call('ZADD', key, now, now)
    -- 设置 key 的过期时间为窗口大小
    redis.call('PEXPIRE', key, window)
    -- 未达到限流条件
    return "false"
end
