package ratelimiter

import (
	"onbio/logger"
	"onbio/redis"

	redigo "github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
)

// NewRateLimiter 新建一个redis的限流
func NewRateLimiter(key string, timeout int, limit int) (err error) {
	conn := redis.GetConn("onbio")
	defer conn.Close()

	cacheKey := "rate_limiter_" + key

	_, err = conn.Do("set", cacheKey, limit, "EX", timeout, "NX")
	if err != nil {
		logger.Warn("rate limiter set failed", zap.String("key", key), zap.Error(err))
		return
	}
	return
}

func IsRateLimiterExisted(key string) (err error) {
	conn := redis.GetConn("onbio")
	defer conn.Close()

	cacheKey := "rate_limiter_" + key
	_, err = conn.Do("get", cacheKey)
	if err != nil {
		logger.Warn("rate limiter get failed", zap.String("key", key), zap.Error(err))
		return
	}
	return
}

// RateLimitAllow  限频判断
func RateLimitAllow(key string) bool {
	conn := redis.GetConn("onbio")
	defer conn.Close()

	cacheKey := "rate_limiter_" + key

	luaScript := `
		local key =  KEYS[1];
		local num = redis.call("decr", key);
		if(num > 0) then
			return num;
		end

		local ttl = redis.call("ttl", key);
		if(ttl == -1) then
			redis.call("del", key);
		end
		return 0;
	`

	script := redigo.NewScript(1, luaScript)
	reply, err := redigo.Int(script.Do(conn, cacheKey))
	if err != nil {
		return false
	}

	if reply > 0 {
		return true
	}

	return false
}
