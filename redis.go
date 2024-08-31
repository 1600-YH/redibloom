package re_bloom

import (
	"context"
	"github.com/gomodule/redigo/redis"
)

type RedisClient struct {
	pool *redis.Pool
}

func NewRedisClient(pool *redis.Pool) *RedisClient {
	return &RedisClient{
		pool: pool,
	}
}

// 执行 lua 脚本，保证复合操作的原子性
func (r *RedisClient) Eval(ctx context.Context, src string, keyCount int, keysAndArgs []interface{}) (interface{}, error) {
	args := make([]interface{}, 2+len(keysAndArgs))
	args[0] = src
	args[1] = keyCount
	copy(args[2:], keysAndArgs)

	//获取连接
	conn, err := r.pool.GetContext(ctx)
	if err != nil {
		return -1, err
	}

	//放回连接池
	defer conn.Close()

	//执行 lua 脚本
	return conn.Do("eval", args...)
}
