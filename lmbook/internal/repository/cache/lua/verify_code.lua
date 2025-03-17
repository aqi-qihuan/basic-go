-- 获取传入的第一个参数，即缓存的键
local key = KEYS[1]
-- 构建一个用于存储验证码尝试次数的键，格式为 "原键:cnt"
local cntKey = key.. ":cnt"
-- 用户输入的验证码
local expectedCode = ARGV[1]

-- 从 Redis 中获取验证码尝试次数
local cnt = tonumber(redis.call("get", cntKey))
-- 从 Redis 中获取缓存的验证码
local code = redis.call("get", key)

-- 如果尝试次数为空或小于等于 0，则表示验证次数耗尽
if cnt == nil or cnt <= 0 then
    -- 验证次数耗尽了
    return -1
end

-- 如果缓存的验证码与用户输入的验证码相等，则表示验证通过
if code == expectedCode then
    -- 将尝试次数设置为 0，表示验证通过
    redis.call("set", cntKey, 0)
    return 0
else
    -- 如果验证码不匹配，则减少一次尝试次数
    redis.call("decr", cntKey)
    -- 不相等，用户输错了
    return -2
end
