-- 获取传入的第一个参数，即缓存的键
local key = KEYS[1]
-- 构建一个用于存储验证码尝试次数的键，格式为 "原键:cnt"
local cntKey = key.. ":cnt"
-- 获取传入的第二个参数，即准备存储的验证码值
local val = ARGV[1]

-- 调用 Redis 的 ttl 命令获取键的剩余过期时间（以秒为单位）
local ttl = tonumber(redis.call("ttl", key))

-- 根据 ttl 的值进行不同的处理
if ttl == -1 then
    -- 键存在，但没有设置过期时间
    return -2
elseif ttl == -2 or ttl < 540 then
    -- 可以发送验证码
    -- 设置验证码值并设置过期时间为 600 秒
    redis.call("set", key, val)
    redis.call("expire", key, 600)
    -- 设置尝试次数为 3 次，并设置过期时间为 600 秒
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, 600)
    return 0
else
    -- 发送验证码过于频繁
    return -1
end
