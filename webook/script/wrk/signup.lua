---
--- Created by zhang san.
--- DateTime: 2025/4/20 21:39
---

--wrk.method="POST"
--wrk.headers["Content-Type"] = "application/json"
--
---- 使用更可靠的随机数生成
--local random = math.random
--
---- 初始化随机种子
--math.randomseed(os.time() * 1000 + wrk.thread:getid())
--
---- 生成随机ID(UUID格式)
--function generate_random_id()
--    template = "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx"
--    return string.gsub(template, "[xy]", function (c)
--        local v = (c == "x") and random(0, 15) or random(8, 11)
--        return string.format("%x", v)
--    end)
--end
--
---- 生成随机邮箱
--function generate_random_email()
--    local prefix = generate_random_id()
--    return prefix .. "@webook.test"
--end
--
---- 初始化
--function init(args)
---- 每个线程都有一个 cnt，所以是线程安全的
--    cnt = 0
--end
--
---- 构建请求
--function request()
--    cnt = cnt + 1
--
--    -- 为每一个请求生成唯一的邮箱
--    local email = generate_random_email()
--    local body = string.format(
--        '{"email":"%s", "password":"hello#world123", "confirmPassword": "hello#world123"}',
--        email
--    )
--    return wrk.format('POST', wrk.path, wrk.headers, body)
--end
--
---- 响应处理
--function response(status, headers, body)
--    -- 可以在这里添加响应验证逻辑
--    if status >= 400 then
--        local err_msg = "HTTP" .. tostring(status)
--        print(string.format("请求失败：%s - %s", err_msg, body))
--    end
--end

wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"

local random = math.random

-- 兼容单/多线程的随机种子初始化
local thread_id = 0
if wrk and wrk.thread then
    thread_id = wrk.thread:getid()
else
    thread_id = 0
end

-- 在初始化时设置随机种子
math.randomseed(os.time() * 1000 + thread_id)

-- 生成随机ID（UUID格式）
function generate_random_id()
    template = "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx"
    return string.gsub(template, "[xy]", function(c)
        local v = (c == "x") and random(0, 15) or random(8, 11)
        return string.format("%x", v)
    end)
end

-- 生成随机邮箱
function generate_random_email()
    local prefix = generate_random_id()
    return prefix .. "@webook.test"
end

-- 初始化计数器
function init(args)
    cnt = 0
end

-- 构建请求
function request()
    cnt = cnt + 1

    local email = generate_random_email()

    local body = string.format(
        '{"email":"%s", "password":"hello#world123", "confirmPassword": "hello#world123"}',
        email
    )

    return wrk.format('POST', wrk.path, wrk.headers, body)
end

-- 响应处理（可选）
function response(status, headers, body)
    -- 可选：验证响应状态码
    if status >= 400 then
        local err_msg = "HTTP " .. tostring(status)
        print(string.format("请求失败: %s - %s", err_msg, body))
    end
end
