package driver

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

type RedisOption struct {
	Host               string
	Port               int
	DialConnectTimeout time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	MaxIdle            int
	MaxActive          int
	IdleTimeout        time.Duration
	Wait               bool
	MaxConnLifetime    time.Duration
	Password           string
	Namespace          string
}

func NewRedis(option RedisOption) *redis.Pool {
	dialConnectTimeoutOption := redis.DialConnectTimeout(option.DialConnectTimeout)
	readTimeoutOption := redis.DialReadTimeout(option.ReadTimeout)
	writeTimeoutOption := redis.DialWriteTimeout(option.WriteTimeout)

	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(fmt.Sprintf("redis://%s@%s:%d", option.Password, option.Host, option.Port), dialConnectTimeoutOption, readTimeoutOption, writeTimeoutOption)
			if err != nil {
				panic(fmt.Errorf("ERROR connect redis | %v", err))
			}

			if option.Password != "" {
				if _, err := c.Do("AUTH", option.Password); err != nil {
					c.Close()
					panic(fmt.Sprintf("ERROR on AUTH redis | %v", err))
				}
			}

			if _, err := c.Do("SELECT", option.Namespace); err != nil {
				c.Close()
				panic(fmt.Sprintf("ERROR on SELECT namespace redis | %v", err))
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			if err != nil {
				panic(fmt.Sprintf("ERROR on PING redis | %v", err))
			}
			return nil
		},
		MaxIdle:         option.MaxIdle,
		MaxActive:       option.MaxActive,
		IdleTimeout:     option.IdleTimeout * time.Second,
		Wait:            option.Wait,
		MaxConnLifetime: option.MaxConnLifetime * time.Second,
	}

	return pool
}
