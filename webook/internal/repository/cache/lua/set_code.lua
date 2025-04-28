---
--- Created by zhang san.
--- DateTime: 2025/4/26 19:43
---

local key = KEYS[1]
local val = ARGV[1]
local cntKey = key..":cnt"
local ttl = tonumber(redis.call("ttl", key))

-- 强制转换为整数类型
local expireSeconds = 600 -- 明确整数类型

-- 首次发送或验证码已过期时存储新验证码
if ttl == -2 or ttl < 540 then
    -- 强制覆盖旧数据
    redis.call("SET", key, tostring(val))
    redis.call("EXPIRE", key, expireSeconds)
    -- 初始化验证次数为 3（字符串类型）
    redis.call("SET", cntKey, "3")
    redis.call("EXPIRE", cntKey, expireSeconds)
    return 0
else
    -- 验证码仍有效，拒绝重复发送
    return -1
end



--你的验证码在 Redis 上的 key
-- phone_code:login:189xxxxxxxx
--local key = KEYS[1]
---- 验证次数，我们一个验证码，最多重复三次，这个记录还可以验证几次
---- phone_code:login:189xxxxxxxx:cnt
--local cntKey = key..":cnt"
---- 你的验证码 123456
--local val= ARGV[1]
---- 过期时间
--local ttl = tonumber(redis.call("ttl", key))
--if ttl == -1 then
--    --    key 存在，但是没有过期时间
--    -- 系统错误，你的同事手贱，手动设置了这个 key，但是没给过期时间
--    return -2
--    --    540 = 600-60 九分钟
--elseif ttl == -2 or ttl < 540 then
--    redis.call("set", key, tostring(val))
--    redis.call("expire", key, 600)
--    redis.call("set", cntKey, 3)
--    redis.call("expire", cntKey, 600)
--    -- 完美，符合预期
--    return 0
--else
--    -- 发送太频繁
--    return -1
--end