# jwt-gowrapper
[![circleci](https://circleci.com/gh/yuchanns/jwt-gowrapper.svg?style=svg)](https://circleci.com/gh/yuchanns/jwt-gowrapper)
[![GoDoc](https://pkg.go.dev/badge/github.com/yuchanns/jwt-gowrapper)](https://pkg.go.dev/github.com/yuchanns/jwt-gowrapper)

a wrapper of [dgrijalva/jwt-go](https://github.com/dgrijalva/jwt-go) with storage interface

## example
> this example shows how this package work with a builtin redigo storage. Certainly users can implemented a custom storage by themselves.

* generate a json web token:
```go
// conn is an instance of redis.Conn from `github.com/gomodule/redigo/redis`
storage := NewStorageRedis(conn)
j := NewDefaultJwt(storage, secretKey, 2*time.Second)
token, err := j.Generate(jwt.MapClaims{
    "id":   111,
    "name": "yuchanns",
})
```
* validate a jwt:
> the validate method will expand the living time of the given token once it is valid.

```go
err := j.Validate(token)
```
* invalidate a jwt:
```go
err := j.Invalidate(token)
```

## custom storage
do implement the Storage interface below:
```go
type Storage interface {
	// Check a key if it is exist and return its remaining time to live
	TTL (key string) (duration time.Duration, err error)
	// Set key to hold a string value and to timeout after a given time.Duration
	SetEx(key, val string, duration time.Duration) error
	// Extend a key's living time duration before timeout
	ExtendKey(key string, duration time.Duration) error
	// Delete a key
	DelKey(key string) error
}
```