package utils

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"time"
)

// 主要是对通过redigo对redis的访问进行一次简单的封装

var (
	redisPool *redis.Pool
)

// 获取一个redis pool
func RedisPoolInit(server, password string) *redis.Pool {
	//redis pool
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			// 进行连接测试
			_, err := c.Do("PING")
			return err
		},
	}
}

// 初始化
func RedisInit(server, password string) {
	if redisPool != nil {
		return
	}
	redisPool = RedisPoolInit(server, password)
}

// 销毁redis pool相关资源
func RedisDestroy() {
	if redisPool != nil {
		if err := redisPool.Close(); err != nil {
			LogTraceE("%s", err.Error())
		} else {
			redisPool = nil
		}
	}
}

//LPOP key 从队列的左边出队一个元素
func RedisLPop(key string) (value string, err error) {
	if redisPool == nil {
		err = errors.New("the redis pool is nil")
		return
	}
	// 从redis连接池里面获取一个有效的conn
	conn := redisPool.Get()
	defer conn.Close()
	//执行LPOP操作
	value, err = redis.String(conn.Do("LPOP", key))
	// if err != nil {
	// 	LogTraceE("LPOP:%s,err:%s", key, err.Error())
	// }
	return
}

//HGET key field 返回 key 指定的哈希集中该字段所关联的值
func RedisHGet(key, field string) (value string, err error) {
	if redisPool == nil {
		err = errors.New("the redis pool is nil")
		return
	}
	// 从redis连接池里面获取一个有效的conn
	conn := redisPool.Get()
	defer conn.Close()
	//执行LPOP操作
	value, err = redis.String(conn.Do("HGet", key, field))
	if err != nil {
		LogTraceE("HGet:%s %s,err:%s", key, field, err.Error())
	}
	return
}

func RedisGet(key string) (value string, err error) {
	if redisPool == nil {
		err = errors.New("the redis pool is nil")
		return
	}
	// 从redis连接池里面获取一个有效的conn
	conn := redisPool.Get()
	defer conn.Close()
	//执行LPOP操作
	value, err = redis.String(conn.Do("Get", key))
	if err != nil {
		LogTraceE("Get:%s ,err:%s", key, err.Error())
	}
	return
}

func RedisLPush(key, value string) (err error) {
	if redisPool == nil {
		err = errors.New("the redis pool is nil")
		return
	}
	// 从redis连接池里面获取一个有效的conn
	conn := redisPool.Get()
	defer conn.Close()
	//执行LPUSH操作
	_, err = conn.Do("LPUSH", key, value)
	if err != nil {
		LogTraceE("LPUSH:%s %s,err:%s", key, value, err.Error())
	}
	return
}

func RedisRPush(key, value string) (err error) {
	if redisPool == nil {
		err = errors.New("the redis pool is nil")
		return
	}
	// 从redis连接池里面获取一个有效的conn
	conn := redisPool.Get()
	defer conn.Close()
	//执行LPUSH操作
	_, err = conn.Do("RPUSH", key, value)
	if err != nil {
		LogTraceE("RPUSH:%s %s,err:%s", key, value, err.Error())
	}
	return
}

func RedisLLen(key string) (value int, err error) {
	if redisPool == nil {
		err = errors.New("the redis pool is nil")
		return
	}
	// 从redis连接池里面获取一个有效的conn
	conn := redisPool.Get()
	defer conn.Close()
	//执行LPUSH操作
	value, err = redis.Int(conn.Do("LLEN", key))
	if err != nil {
		LogTraceE("LLEN:%s,err:%s", key, err.Error())
	}
	return
}

func RedisDo(cmd string, args ...interface{}) (err error) {
	if redisPool == nil {
		err = errors.New("the redis pool is nil")
		return
	}
	// 从redis连接池里面获取一个有效的conn
	conn := redisPool.Get()
	defer conn.Close()
	//执行LPUSH操作
	_, err = conn.Do(cmd, args...)
	if err != nil {
		LogTraceE("%s,err:%s", cmd, err.Error())
	}
	return
}
