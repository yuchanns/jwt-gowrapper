package jwt_gowrapper

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"time"
)

type Jwt struct {
	storage  Storage
	token    *jwt.Token
	secret   []byte
	duration time.Duration
}

// NewDefaultJwt returns an instance of jwt takes the signing method jwt.SigningMethodHS256
func NewDefaultJwt(storage Storage, secret []byte, duration time.Duration) *Jwt {
	return NewJwtWithMethod(storage, secret, duration, jwt.SigningMethodHS256)
}

// NewJwtWithMethod returns an instance of jwt takes the given signing method
func NewJwtWithMethod(storage Storage, secret []byte, duration time.Duration, method jwt.SigningMethod) *Jwt {
	token := jwt.New(method)
	return &Jwt{
		storage:  storage,
		token:    token,
		secret:   secret,
		duration: duration,
	}
}

// Generate generate a tokenString with given jwt.MapClaims
func (j *Jwt) Generate(m jwt.MapClaims) (string, error) {
	j.token.Claims = m
	tokenString, err := j.token.SignedString(j.secret)
	if err != nil {
		return "", err
	}
	key, err := j.getKey()
	if err != nil {
		return "", err
	}
	if err := j.storage.SetEx(key, tokenString, j.duration); err != nil {
		return "", err
	}
	return tokenString, nil
}

// Validate verify if the tokenString is valid
func (j *Jwt) Validate(tokenString string) error {
	if err := j.parseToken(tokenString); err != nil {
		return err
	}
	if !j.token.Valid {
		return errors.New("token validate failed")
	}
	key, err := j.getKey()
	if err != nil {
		return err
	}
	ttl, err := j.storage.TTL(key)
	if err != nil {
		return err
	}
	if ttl.Seconds() <= 0 {
		return errors.New("the token has been expired")
	}
	if err := j.storage.ExtendKey(key, j.duration); err != nil {
		return err
	}
	return nil
}

// Invalidate make the given tokenString invalid
func (j *Jwt) Invalidate(tokenString string) error {
	if err := j.parseToken(tokenString); err != nil {
		return err
	}
	key, err := j.getKey()
	if err != nil {
		return err
	}
	return j.storage.DelKey(key)
}

// getKey build a key string from the jwt.Claims
func (j *Jwt) getKey() (string, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	defer gz.Close()
	jsonByte, err := json.Marshal(j.token.Claims)
	if err != nil {
		return "", err
	}
	if _, err := gz.Write(jsonByte); err != nil {
		return "", err
	}
	if err := gz.Flush(); err != nil {
		return "", err
	}
	return b.String(), nil
}

// parseToken parse the given tokenString and set to jwt.Token
func (j *Jwt) parseToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})
	if err != nil {
		return err
	}
	j.token = token
	return nil
}
