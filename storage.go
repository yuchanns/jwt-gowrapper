package jwt_gowrapper

import "time"

// Storage is an interface which will be used as an arg in Jwt.
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
