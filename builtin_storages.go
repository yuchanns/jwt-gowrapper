package jwt_gowrapper

import (
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"time"
)

type StorageRedis struct {
	conn redis.Conn
}

func NewStorageRedis(conn redis.Conn) Storage {
	return &StorageRedis{
		conn: conn,
	}
}

func (r StorageRedis) TTL (key string) (duration time.Duration, err error) {
	ttl, err := redis.Int64(r.conn.Do("TTL", key))
	if err != nil {
		return 0, err
	}
	ttlDuration := time.Duration(ttl)*time.Second
	if ttlDuration <= 0 {
		return 0, errors.New("the token has been expired")
	}
	return ttlDuration, nil
}

func (r StorageRedis) SetEx(key, val string, duration time.Duration) error {
	_, err := r.conn.Do("SETEX", key, duration.Seconds(), val)
	return err
}

func (r StorageRedis) ExtendKey(key string, duration time.Duration) error {
	_, err := r.conn.Do("EXPIRE", key, duration.Seconds())
	return err
}

func (r StorageRedis) DelKey(key string) error {
	_, err := r.conn.Do("DEL", key)
	return err
}
