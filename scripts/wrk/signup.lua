wrk.method="POST"
wrk.headers["Content-Type"] = "application/json"

-- 导入 math 库，用于生成随机数
local random = math.random

-- 定义一个函数，用于生成 UUID
local function uuid()
    -- 定义一个 UUID 的模板字符串
    local template ='xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'
    -- 使用 string.gsub 函数替换模板中的 'x' 和 'y'
    return string.gsub(template, '[xy]', function (c)
        -- 根据字符是 'x' 还是 'y'，生成不同范围的随机数
        local v = (c == 'x') and random(0, 0xf) or random(8, 0xb)
        -- 将随机数格式化为十六进制字符串并返回
        return string.format('%x', v)
    end)
end

-- 初始化函数，在每个线程开始时调用
function init(args)
    -- 初始化计数器 cnt，每个线程都有一个独立的 cnt，因此是线程安全的
    cnt = 0
    -- 生成一个随机的前缀 prefix
    prefix = uuid()
end

-- 请求函数，构造 HTTP 请求体
function request()
    -- 使用前缀 prefix 和计数器 cnt 构造一个邮箱地址
    body=string.format('{"email":"%s%d@qq.com", "password":"hello#world123", "confirmPassword": "hello#world123"}', prefix, cnt)
    -- 每次请求后增加计数器的值
    cnt = cnt + 1
    -- 返回格式化后的 HTTP 请求
    return wrk.format('POST', wrk.path, wrk.headers, body)
end

-- 响应函数，目前为空，可用于处理响应数据
function response()

end
