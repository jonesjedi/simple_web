package redis

import (
	"fmt"
	"sync"
	"time"

	"onbio/conf"

	"github.com/garyburd/redigo/redis"
)

const (
	redisMaxIdle        = 10  //最大空闲连接数
	redisIdleTimeout    = 240 //最大空闲连接时间
	redisConnectTimeout = 1
	redisReadTimeout    = 1
	redisWriteTimeout   = 1
)

var RedisPools sync.Map
var RedisPool *redis.Pool

// RedisPoolsMap 连接池服务
type redisPoolsMap struct {
	Pools map[string]*redis.Pool
}

// GRedisPools 全局的redis池
var GRedisPools = &redisPoolsMap{}

// NewRedisPool 返回redis连接池
func NewRedisPool(instance string) *redis.Pool {
	RedisURL, RedisPassword := getRedisConfig(instance)
	fmt.Printf("init redis conn, dsn:%s\n", RedisURL)
	return &redis.Pool{
		MaxIdle:     redisMaxIdle,
		IdleTimeout: redisIdleTimeout * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", RedisURL,
				redis.DialConnectTimeout(time.Duration(redisConnectTimeout)*time.Second),
				redis.DialReadTimeout(time.Duration(redisReadTimeout)*time.Second),
				redis.DialWriteTimeout(time.Duration(redisWriteTimeout)*time.Second))
			if err != nil {
				return nil, err
			}
			if RedisPassword != "" {
				if _, authErr := c.Do("AUTH", RedisPassword); authErr != nil {
					return nil, fmt.Errorf("redis auth password error: %s", authErr)
				}
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				return fmt.Errorf("ping redis error: %s", err)
			}
			return nil
		},
	}
}

// Init 初始化连接池
func Init(instance string) {

	if GRedisPools.Pools == nil {
		GRedisPools.Pools = make(map[string]*redis.Pool)
	}
	con, ok := RedisPools.Load(instance)
	if ok == false {
		RedisPool = NewRedisPool(instance)
		GRedisPools.Pools[instance] = RedisPool
		RedisPools.Store(instance, RedisPool)
	} else {
		//切换到当前实例
		RedisPool = con.(*redis.Pool)
		GRedisPools.Pools[instance] = RedisPool
	}

}

// 获取Conn
func GetConn(instance string) redis.Conn {
	return GRedisPools.Pools[instance].Get()
}

func getRedisConfig(instance string) (string, string) {
	DBConf := conf.GetRedisConfig()
	dsn := fmt.Sprintf("%s:%d", DBConf.RedisHost, DBConf.RedisPort)
	return dsn, DBConf.RedisPwd
}

func getError(err error) error {
	return err
	//return errors.New("redis exception")
}

// Set string
func (p *redisPoolsMap) Set(instance string, key string, value string) error {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		return getError(err)
	}

	return err
}

func (p *redisPoolsMap) Del(instance string, keys []string) (int, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	param := []interface{}{}
	for _, v := range keys {
		param = append(param, v)
	}
	rs, err := redis.Int(conn.Do("DEL", param...))
	if err != nil {
		return 0, getError(err)
	}

	return rs, nil
}

/**
 * 第一个参数为是否设置成功，0失败1成功
 */
func (p *redisPoolsMap) SetNx(instance string, key string, value string) (int, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rs, err := redis.Int(conn.Do("SETNX", key, value))
	if err != nil {
		return rs, getError(err)
	}

	return rs, err
}

func (p *redisPoolsMap) SetEx(instance string, key string, seconds int, value string) error {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	_, err := conn.Do("SETEX", key, seconds, value)
	if err != nil {
		return getError(err)
	}

	return err
}

func (p *redisPoolsMap) PSetEx(instance string, key string, milliseconds int, value string) error {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	_, err := conn.Do("PSETEX", key, milliseconds, value)
	if err != nil {
		return getError(err)
	}

	return err
}

//返回OK表示设置成功
//返回""空字符串表示键名已存在
func (p *redisPoolsMap) SetPxNx(instance string, key string, value interface{}, milliseconds int) (string, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rs, err := conn.Do("SET", key, value, "PX", milliseconds, "NX")
	if err != nil {
		return "", getError(err)
	}
	if rs == nil {
		return "", nil
	}

	rsStr, err := redis.String(rs, err)

	return rsStr, err
}

func (p *redisPoolsMap) Get(instance string, key string) (value string, err error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rs, err := conn.Do("GET", key)
	if err != nil {
		value = ""
		err = getError(err)
		return
	}
	if rs == nil {
		value = ""
	} else {
		value, err = redis.String(rs, nil)
	}

	return
}

func (p *redisPoolsMap) GetSet(instance string, key string, value string) (oldValue string, err error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rs, err := conn.Do("GETSET", key, value)
	if err != nil {
		value = ""
		err = getError(err)
		return
	}
	if rs == nil {
		oldValue = ""
	} else {
		oldValue, err = redis.String(rs, nil)
	}

	return
}

func (p *redisPoolsMap) Incr(instance string, key string) (int, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rs, err := redis.Int(conn.Do("INCR", key))
	if err != nil {
		return rs, getError(err)
	}

	return rs, err
}

func (p *redisPoolsMap) IncrBy(instance string, key string, increment int) (int, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rs, err := redis.Int(conn.Do("INCRBY", key, increment))
	if err != nil {
		return rs, getError(err)
	}

	return rs, err
}

func (p *redisPoolsMap) Decr(instance string, key string) (int, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rs, err := redis.Int(conn.Do("DECR", key))
	if err != nil {
		return rs, getError(err)
	}

	return rs, err
}

func (p *redisPoolsMap) DecrBy(instance string, key string, decrement int) (int, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rs, err := redis.Int(conn.Do("DECRBy", key, decrement))
	if err != nil {
		return rs, getError(err)
	}

	return rs, err
}

func (p *redisPoolsMap) MSet(instance string, keyMapValues map[string]string) error {
	if len(keyMapValues) < 1 {
		return nil
	}

	conn := p.Pools[instance].Get()
	defer conn.Close()

	params := []interface{}{}
	for k, v := range keyMapValues {
		params = append(params, k, v)
	}

	_, err := conn.Do("MSET", params...)
	if err != nil {
		return getError(err)
	}

	return err
}

func (p *redisPoolsMap) MGet(instance string, keys []string) ([]string, error) {
	if len(keys) < 1 {
		return []string{}, nil
	}

	conn := p.Pools[instance].Get()
	defer conn.Close()

	var keysInterface []interface{}
	for _, v := range keys {
		keysInterface = append(keysInterface, v)
	}
	rs, err := redis.Strings(conn.Do("MGET", keysInterface...))
	if err != nil {
		return nil, getError(err)
	}

	return rs, nil
}

//hash
func (p *redisPoolsMap) HSet(instance string, key string, field string, value string) error {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	_, err := conn.Do("HSET", key, field, value)
	if err != nil {
		return getError(err)
	}

	return err
}

func (p *redisPoolsMap) HSetNx(instance string, key string, field string, value string) (int, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rs, err := redis.Int(conn.Do("HSETNX", key, field, value))
	if err != nil {
		return rs, getError(err)
	}

	return rs, err
}

func (p *redisPoolsMap) HGet(instance string, key string, field string) (string, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rs, err := conn.Do("HGET", key, field)
	if err != nil {
		return "", getError(err)
	}

	if rs == nil {
		return "", err
	} else {
		return redis.String(rs, err)
	}
}

func (p *redisPoolsMap) HExists(instance string, key string, field string) (bool, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rs, err := redis.Bool(conn.Do("HEXISTS", key, field))
	if err != nil {
		return false, getError(err)
	}
	return rs, err
}

func (p *redisPoolsMap) HDel(instance string, key string, fields []string) error {
	if len(fields) < 1 {
		return nil
	}

	conn := p.Pools[instance].Get()
	defer conn.Close()

	params := []interface{}{}
	params = append(params, key)
	for _, v := range fields {
		params = append(params, v)
	}

	_, err := conn.Do("HDEL", params...)
	if err != nil {
		return getError(err)
	}
	return err
}

func (p *redisPoolsMap) HLen(instance string, key string) (int, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rs, err := redis.Int(conn.Do("HLEN", key))
	if err != nil {
		return 0, getError(err)
	}

	return rs, err
}

func (p *redisPoolsMap) HIncrBy(instance string, key string, field string, increment int) (int, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rs, err := redis.Int(conn.Do("HINCRBY", key, field, increment))
	if err != nil {
		return 0, getError(err)
	}
	return rs, err
}

func (p *redisPoolsMap) HMSet(instance string, key string, fieldMapValues map[string]string) error {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	params := []interface{}{}
	params = append(params, key)
	for k, v := range fieldMapValues {
		params = append(params, k, v)
	}

	_, err := conn.Do("HMSET", params...)
	if err != nil {
		return getError(err)
	}

	return err
}

func (p *redisPoolsMap) HMGet(instance string, key string, fields []string) (map[string]string, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	params := []interface{}{}
	params = append(params, key)
	for _, v := range fields {
		params = append(params, v)
	}
	rs, err := redis.Strings(conn.Do("HMGET", params...))
	if err != nil {
		return nil, getError(err)
	}

	rspData := map[string]string{}
	for k, v := range fields {
		rspData[v] = rs[k]
	}

	return rspData, nil
}

func (p *redisPoolsMap) HGetAll(instance string, key string) (map[string]string, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rs, err := redis.StringMap(conn.Do("HGETALL", key))
	if err != nil {
		return nil, getError(err)
	}

	return rs, nil
}

//list //待实现
func (p *redisPoolsMap) LPush() {}

func (p *redisPoolsMap) RPush() {}

func (p *redisPoolsMap) LPop() {}

func (p *redisPoolsMap) RPop() {}

func (p *redisPoolsMap) LLen() {}

//set //待实现
func (p *redisPoolsMap) SAdd(instance string, key string, member string) (int, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rs, err := redis.Int(conn.Do("SADD", key, member))
	if err != nil {
		return rs, getError(err)
	}

	return rs, err
}

func (p *redisPoolsMap) SMembers(instance string, key string) ([]string, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rs, err := redis.Strings(conn.Do("SMEMBERS", key))
	if err != nil {
		return rs, getError(err)
	}
	return rs, err
}

func (p *redisPoolsMap) SIsMember(instance string, key string, member string) (int, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rs, err := redis.Int(conn.Do("SISMEMBER", key, member))
	if err != nil {
		return rs, getError(err)
	}
	return rs, err
}

func (p *redisPoolsMap) SPop() {}

func (p *redisPoolsMap) SRandMember() {}

func (p *redisPoolsMap) SRem(instance string, key string, member string) (int, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rs, err := redis.Int(conn.Do("SREM", key, member))
	if err != nil {
		return rs, getError(err)
	}
	return rs, err
}

func (p *redisPoolsMap) SCard(instance string, key string) (int, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rs, err := redis.Int(conn.Do("SCARD", key))
	if err != nil {
		return rs, getError(err)
	}
	return rs, err
}

//zset
func (p *redisPoolsMap) ZAdd(instance string, key string, score int, member string) (int, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rs, err := redis.Int(conn.Do("ZADD", key, score, member))
	if err != nil {
		return rs, getError(err)
	}

	return rs, err
}

func (p *redisPoolsMap) ZAddMembers(instance string, key string, memberScoreInt map[string]int) (int, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	if len(memberScoreInt) < 1 {
		return 0, nil
	}

	kvArr := make([]interface{}, 0)
	kvArr = append(kvArr, key)
	for k, v := range memberScoreInt {
		kvArr = append(kvArr, v)
		kvArr = append(kvArr, k)
	}

	rs, err := redis.Int(conn.Do("ZADD", kvArr...))
	if err != nil {
		return rs, getError(err)
	}

	return rs, err
}

func (p *redisPoolsMap) ZScore(instance string, key string, member string) (int, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	doRs, err := conn.Do("ZSCORE", key, member)
	if err != nil {
		return 0, getError(err)
	}
	if doRs != nil {
		rs, _ := redis.Int(doRs, err)
		return rs, nil
	} else {
		return 0, nil
	}
}

func (p *redisPoolsMap) ZIncrBy(instance string, key string, increment int, member string) (int, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rs, err := redis.Int(conn.Do("ZINCRBY", key, increment, member))
	if err != nil {
		return rs, getError(err)
	}

	return rs, err
}

func (p *redisPoolsMap) ZCard(instance string, key string) (int, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rs, err := redis.Int(conn.Do("ZCARD", key))
	if err != nil {
		return rs, getError(err)
	}

	return rs, err
}

func (p *redisPoolsMap) ZCount(instance string, key string, min interface{}, max interface{}) (int, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rs, err := redis.Int(conn.Do("ZCOUNT", key, min, max))
	if err != nil {
		return rs, getError(err)
	}

	return rs, err
}

func (p *redisPoolsMap) ZRange(instance string, key string, start int, stop int) ([]string, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rs, err := redis.Strings(conn.Do("ZRANGE", key, start, stop))
	if err != nil {
		return rs, getError(err)
	}

	return rs, err
}

func (p *redisPoolsMap) ZRangeWithScores(instance string, key string, start int, stop int) ([]map[string]string, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rspArr := []map[string]string{}
	rs, err := redis.Strings(conn.Do("ZRANGE", key, start, stop, "WITHSCORES"))
	if err != nil {
		return rspArr, getError(err)
	}
	for i := 0; i < len(rs); i += 2 {
		key := rs[i]
		value := rs[i+1]
		rspArr = append(rspArr, map[string]string{
			"member": key,
			"value":  value,
		})
	}

	return rspArr, err
}

func (p *redisPoolsMap) ZRevRange(instance string, key string, start int, stop int) ([]string, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rs, err := redis.Strings(conn.Do("ZREVRANGE", key, start, stop))
	if err != nil {
		return rs, getError(err)
	}

	return rs, err
}

func (p *redisPoolsMap) ZRevRangeWithScores(instance string, key string, start int, stop int) ([]map[string]string, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rspArr := []map[string]string{}
	rs, err := redis.Strings(conn.Do("ZREVRANGE", key, start, stop, "WITHSCORES"))
	if err != nil {
		return rspArr, getError(err)
	}
	for i := 0; i < len(rs); i += 2 {
		key := rs[i]
		value := rs[i+1]
		rspArr = append(rspArr, map[string]string{
			"member": key,
			"value":  value,
		})
	}

	return rspArr, err
}

func (p *redisPoolsMap) ZRangeByScore(instance string, key string, min interface{}, max interface{}, limitOffset int, count int) ([]string, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rs, err := redis.Strings(conn.Do("ZRANGEBYSCORE", key, min, max, "LIMIT", fmt.Sprint(limitOffset), fmt.Sprint(count)))
	if err != nil {
		return rs, getError(err)
	}

	return rs, err
}

/**
 * 返回数组，每个元素都是map 包含member和value
 */
func (p *redisPoolsMap) ZRangeByScoreWithScores(instance string, key string, min interface{}, max interface{}, limitOffset int, count int) ([]map[string]string, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rspArr := []map[string]string{}
	rs, err := redis.Strings(conn.Do("ZRANGEBYSCORE", key, min, max, "WITHSCORES", "LIMIT", fmt.Sprint(limitOffset), fmt.Sprint(count)))
	if err != nil {
		return rspArr, getError(err)
	}

	for i := 0; i < len(rs); i += 2 {
		key := rs[i]
		value := rs[i+1]
		rspArr = append(rspArr, map[string]string{
			"member": key,
			"value":  value,
		})
	}

	return rspArr, err
}

func (p *redisPoolsMap) ZRevRangeByScore(instance string, key string, max interface{}, min interface{}, limitOffset int, count int) ([]string, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rs, err := redis.Strings(conn.Do("ZREVRANGEBYSCORE", key, max, min, "LIMIT", fmt.Sprint(limitOffset), fmt.Sprint(count)))
	if err != nil {
		return rs, getError(err)
	}

	return rs, err
}

/**
 * 返回数组，每个元素都是map 包含member和value
 */
func (p *redisPoolsMap) ZRevRangeByScoreWithScores(instance string, key string, max interface{}, min interface{}, limitOffset int, count int) ([]map[string]string, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rspArr := []map[string]string{}
	rs, err := redis.Strings(conn.Do("ZREVRANGEBYSCORE", key, max, min, "WITHSCORES", "LIMIT", fmt.Sprint(limitOffset), fmt.Sprint(count)))
	if err != nil {
		return rspArr, getError(err)
	}

	for i := 0; i < len(rs); i += 2 {
		key := rs[i]
		value := rs[i+1]
		rspArr = append(rspArr, map[string]string{
			"member": key,
			"value":  value,
		})
	}

	return rspArr, err
}

func (p *redisPoolsMap) ZRank(instance string, key string, member string) (int, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	doRs, err := conn.Do("ZRANK", key, member)
	if err != nil {
		return 0, getError(err)
	} else {
		rs, _ := redis.Int(doRs, err)
		return rs, nil
	}
}

func (p *redisPoolsMap) ZRevRank(instance string, key string, member string) (int, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	doRs, err := conn.Do("ZREVRANK", key, member)
	if err != nil {
		return 0, getError(err)
	} else {
		rs, _ := redis.Int(doRs, err)
		return rs, nil
	}
}

func (p *redisPoolsMap) ZRem(instance string, key string, member string) (int, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rs, err := redis.Int(conn.Do("ZREM", key, member))
	if err != nil {
		return rs, getError(err)
	}

	return rs, err
}

func (p *redisPoolsMap) ZRemMembers(instance string, key string, members []string) (int, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rediParam := []interface{}{}
	rediParam = append(rediParam, key)
	for _, v := range members {
		rediParam = append(rediParam, v)
	}

	rs, err := redis.Int(conn.Do("ZREM", rediParam...))
	if err != nil {
		return rs, getError(err)
	}

	return rs, err
}

func (p *redisPoolsMap) ZRemRangeByRank(instance string, key string, start int, stop int) (int, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rs, err := redis.Int(conn.Do("ZREMRANGEBYRANK", key, start, stop))
	if err != nil {
		return rs, getError(err)
	}

	return rs, err
}

func (p *redisPoolsMap) ZRemRangeByScore(instance string, key string, min int, max int) (int, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rs, err := redis.Int(conn.Do("ZREMRANGEBYSCORE", key, min, max))
	if err != nil {
		return rs, getError(err)
	}

	return rs, err
}

//bit
func (p *redisPoolsMap) SetBit(instance string, key string, offset int, value int) error {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	_, err := redis.Int(conn.Do("SETBIT", key, offset, value))
	if err != nil {
		return getError(err)
	}

	return err
}

func (p *redisPoolsMap) GetBit(instance string, key string, offset int) (int, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rs, err := redis.Int(conn.Do("GETBIT", key, offset))
	if err != nil {
		return 0, getError(err)
	}

	return rs, err
}

//自动过期
func (p *redisPoolsMap) Expire(instance string, key string, seconds int) error {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	_, err := conn.Do("EXPIRE", key, seconds)
	if err != nil {
		return getError(err)
	}

	return err
}

func (p *redisPoolsMap) ExpireAt(instance string, key string, timestamp int) error {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	_, err := conn.Do("EXPIREAT", key, timestamp)
	if err != nil {
		return getError(err)
	}

	return err
}

/**
 *返回剩余生存时间
 *当 key 不存在时，返回 -2 。 当 key 存在但没有设置剩余生存时间时，返回 -1 。 否则，以秒为单位，返回 key 的剩余生存时间。
 */
func (p *redisPoolsMap) TTL(instance string, key string) (int, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rs, err := redis.Int(conn.Do("TTL", key))
	if err != nil {
		return rs, getError(err)
	}

	return rs, err
}

func (p *redisPoolsMap) Persist(instance string, key string) error {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	_, err := conn.Do("PERSIST", key)
	if err != nil {
		return getError(err)
	}

	return err
}

/**
 *以毫秒为单位
 */
func (p *redisPoolsMap) PExpire(instance string, key string, milliseconds int) error {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	_, err := conn.Do("PEXPIRE", key, milliseconds)
	if err != nil {
		return getError(err)
	}

	return err
}

func (p *redisPoolsMap) PExpireAt(instance string, key string, millisecTimestamp int) error {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	_, err := conn.Do("EXPIREAT", key, millisecTimestamp)
	if err != nil {
		return getError(err)
	}

	return err
}

func (p *redisPoolsMap) PTtl(instance string, key string) (int, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rs, err := redis.Int(conn.Do("PTTL", key))
	if err != nil {
		return rs, getError(err)
	}

	return rs, err
}

// HGetAllPipeLine hgetallpiple
func (p *redisPoolsMap) HGetAllPipeLine(instance string, keys []string) ([]map[string]string, error) {
	var err error
	conn := p.Pools[instance].Get()
	defer conn.Close()
	for _, key := range keys {
		err = conn.Send("HGETALL", key)
		if err != nil {
			return nil, err
		}
	}
	err = conn.Flush()
	if err != nil {
		return nil, err
	}

	mslice := make([]map[string]string, len(keys))
	for k := range keys {
		m, err := redis.StringMap(conn.Receive())
		if err != nil {
			return nil, err
		}
		mslice[k] = m
	}
	return mslice, nil
}

// ZAddPipeLine ZAddPipeLine
func (p *redisPoolsMap) ZAddPipeLine(instance string, data []map[string]interface{}) ([]int, error) {
	var err error
	conn := p.Pools[instance].Get()
	defer conn.Close()
	for _, d := range data {
		err = conn.Send("ZADD", d["key"], d["score"], d["member"])
		if err != nil {
			return nil, err
		}
	}
	err = conn.Flush()
	if err != nil {
		return nil, err
	}

	islice := make([]int, 0, len(data))
	for range data {
		i, err := redis.Int(conn.Receive())
		if err != nil {
			return nil, err
		}
		islice = append(islice, i)
	}
	return islice, nil
}

// ZRemPipeLine ZRemPipeLine
func (p *redisPoolsMap) ZRemPipeLine(instance string, data []map[string]interface{}) ([]int, error) {
	var err error
	conn := p.Pools[instance].Get()
	defer conn.Close()
	for _, d := range data {
		err = conn.Send("ZREM", d["key"], d["member"])
		if err != nil {
			return nil, err
		}
	}
	err = conn.Flush()
	if err != nil {
		return nil, err
	}

	islice := make([]int, 0, len(data))
	for range data {
		i, err := redis.Int(conn.Receive())
		if err != nil {
			return nil, err
		}
		islice = append(islice, i)
	}
	return islice, nil
}

func (p *redisPoolsMap) SetRange(instance string, key string, offset int, value string) error {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	_, err := redis.Int(conn.Do("SETRANGE", key, offset, value))
	if err != nil {
		return getError(err)
	}

	return err
}

func (p *redisPoolsMap) GetRange(instance string, key string, start int, end int) (value string, err error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rs, err := conn.Do("GETRANGE", key, start, end)
	if err != nil {
		value = ""
		err = getError(err)
		return
	}
	if rs == nil {
		value = ""
	} else {
		value, err = redis.String(rs, nil)
	}
	return
}

func (p *redisPoolsMap) Exists(instance string, key string) (bool, error) {
	conn := p.Pools[instance].Get()
	defer conn.Close()

	rs, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false, getError(err)
	}
	return rs, err
}
