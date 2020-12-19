package ratelimiter

import (
	"onbio/logger"
	"onbio/redis"

	redigo "github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
	"errors"
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
	ret, err := conn.Do("get", cacheKey)
	//logger.Info("get nil ret",zap.Any("ret",test),zap.Error(err))
	if err != nil {
		logger.Warn("rate limiter get failed", zap.String("key", key), zap.Error(err))
		return
	}
	if ret == nil {
		err = errors.New("not existed")
		return 
	}
	return
}

// RateLimitAllow  限频判断
func RateLimitAllow(key string) bool {
	conn := redis.GetConn("onbio")
	defer conn.Close()

	cacheKey := "rate_limiter_" + key
	reply, err := redigo.Int(conn.Do("decr", cacheKey))
	logger.Info("decr ",zap.Int("reply",reply),zap.Error(err))
	if err != nil {
		logger.Error("err do decr",zap.Error(err))
		return false
	}

	if reply > 0 {
		return true
	}

	return false
}
