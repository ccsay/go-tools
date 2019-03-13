// Copyright 2019 go-tools Authors

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// redis工具类
package redis

import (
	"github.com/gomodule/redigo/redis"
	"fmt"
	"time"
	"github.com/liuchonglin/go-utils"
)

type Redis struct {
	RedisPool *redis.Pool
}

func NewRedis(redisConfig *RedisConfig) (*Redis, error) {
	if redisConfig == nil {
		redisConfig = &RedisConfig{}
		redisConfig.defaultValue()
	}
	// MaxIdle：最大的空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态。
	// MaxActive：最大的激活连接数，表示同时最多有N个连接
	// IdleTimeout：最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭
	redisPool := &redis.Pool{
		MaxIdle:     redisConfig.MaxIdle,
		MaxActive:   redisConfig.MaxActive,
		IdleTimeout: time.Duration(redisConfig.IdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			if utils.IsEmpty(redisConfig.Password) {
				return redis.Dial("tcp", redisConfig.Address)
			} else {
				return redis.Dial("tcp", redisConfig.Address,
					redis.DialKeepAlive(1*time.Second),
					redis.DialConnectTimeout(5*time.Second),
					redis.DialReadTimeout(1*time.Second),
					redis.DialPassword(redisConfig.Password))
			}
		},
	}
	if _, err := redisPool.Dial(); err != nil {
		return nil, fmt.Errorf("redis 初始化失败: %v", err)
	}
	return &Redis{RedisPool: redisPool}, nil
}

// 获取 key 对应的 string 值
func (r *Redis) Get(key string) (value string, err error) {
	if utils.IsEmpty(key) {
		return "", fmt.Errorf("key 不能为空")
	}
	conn := r.RedisPool.Get()
	defer conn.Close()
	result, err := conn.Do("get", key)
	if result == nil && err == nil {
		return "", nil
	}
	value, err = redis.String(result, err)
	if err != nil {
		return "", err
	}
	return value, nil
}

// 设置 string 值
func (r *Redis) Set(key string, value string) error {
	return r.SetExpire(key, value, 0)
}

// 删除 key 对应的值
func (r *Redis) Del(key string) error {
	if utils.IsEmpty(key) {
		return fmt.Errorf("key 不能为空")
	}
	conn := r.RedisPool.Get()
	defer conn.Close()

	_, err := conn.Do("del", key)
	if err != nil {
		return err
	}
	return nil
}

// 设置 string 值 和 超时时间
func (r *Redis) SetExpire(key string, value string, ex int) error {
	if utils.IsEmpty(key) || utils.IsEmpty(value) {
		return fmt.Errorf("key or value 不能为空")
	}

	conn := r.RedisPool.Get()
	defer conn.Close()
	_, err := conn.Do("set", key, value)
	if err != nil {
		return err
	}
	if ex > 0 {
		_, err = conn.Do("expire", key, ex)
		if err != nil {
			return err
		}
	}
	return nil
}
