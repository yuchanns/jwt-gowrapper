package jwt_gowrapper

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock"
	"testing"
	"time"
)

func TestJwt_Validate(t *testing.T) {
	conn := redigomock.NewConn()
	defer conn.Close()
	// mock redis
	conn.Command("SETEX").Expect(1)
	conn.Command("TTL").Expect(int64(86395))
	conn.Command("EXPIRE").Expect(1)
	// test start
	storage := NewStorageRedis(conn)
	secretKey := []byte("3mKfbv5Fj47Ujv")
	j := NewDefaultJwt(storage, secretKey, 86400*time.Second)
	token, err := j.Generate(jwt.MapClaims{
		"id":   111,
		"name": "yuchanns",
	})
	if err != nil {
		t.Errorf("failed to generate tokenString: %+v", err)
	}
	err = j.Validate(token)
	if err != nil {
		t.Errorf("token is invalid: %+v", err)
	}
}

func TestJwt_Invalidate(t *testing.T) {
	conn := redigomock.NewConn()
	defer conn.Close()
	// mock redis
	conn.Command("SETEX").Expect(1)
	conn.Command("DEL").Expect(1)
	conn.Command("TTL").Expect(nil)
	// test start
	storage := NewStorageRedis(conn)
	secretKey := []byte("3mKfbv5Fj47Ujv")
	j := NewDefaultJwt(storage, secretKey, 86400*time.Second)
	token, err := j.Generate(jwt.MapClaims{
		"id":   111,
		"name": "yuchanns",
	})
	if err != nil {
		t.Errorf("failed to generate tokenString: %+v", err)
	}
	if err := j.Invalidate(token); err != nil {
		t.Errorf("failed to invalidate tokenString: %+v", err)
	}
	err = j.Validate(token)
	if err == nil {
		t.Errorf("token should not be validate: %+v", err)
	}
}

func TestNewStorageRedis(t *testing.T) {
	pool := &redis.Pool{
		Dial: func() (conn redis.Conn, err error) {
			return redis.Dial("tcp", "redis:6379")
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	conn := pool.Get()
	defer conn.Close()
	storage := NewStorageRedis(conn)
	secretKey := []byte("3mKfbv5Fj47Ujv")
	j := NewDefaultJwt(storage, secretKey, 2*time.Second)
	token, err := j.Generate(jwt.MapClaims{
		"id":   111,
		"name": "yuchanns",
	})
	if err != nil {
		t.Errorf("failed to generate tokenString: %+v", err)
	}
	if err := j.Invalidate(token); err != nil {
		t.Errorf("failed to invalidate tokenString: %+v", err)
	}
	err = j.Validate(token)
	if err == nil {
		t.Errorf("token should not be validate: %+v", err)
	}
	token, err = j.Generate(jwt.MapClaims{
		"id":   111,
		"name": "yuchanns",
	})
	if err != nil {
		t.Errorf("failed to generate tokenString: %+v", err)
	}
	err = j.Validate(token)
	if err != nil {
		t.Errorf("token is invalid: %+v", err)
	}
	time.Sleep(3 * time.Second)
	err = j.Validate(token)
	if err == nil {
		t.Errorf("token should not be validate: %+v", err)
	}
}
