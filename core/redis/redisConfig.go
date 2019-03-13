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

// redis配置
package redis

import (
	"github.com/liuchonglin/go-utils"
)

var redisConfig *RedisConfig
// redis 配置
type RedisConfig struct {
	// 连接地址
	Address string `json:"address" yaml:"address"`
	// 最大空闲时间
	MaxIdle int `json:"maxIdle" yaml:"maxIdle"`
	// 最大存活时间
	MaxActive int `json:"maxActive" yaml:"maxActive"`
	// 空闲超时时间
	IdleTimeout int `json:"idleTimeout" yaml:"idleTimeout"`
	// 密码
	Password string `json:"password" yaml:"password"`
}

func GetRedisConfig() *RedisConfig {
	if redisConfig == nil {
		redisConfig = &RedisConfig{}
		redisConfig.defaultValue()
	}
	return redisConfig
}

func (r *RedisConfig) defaultValue() {
	if utils.IsEmpty(r.Address) {
		r.Address = "localhost:6379"
	}
	if r.MaxIdle == 0 {
		r.MaxIdle = 16
	}
	if r.MaxActive == 0 {
		r.MaxActive = 60
	}
	if r.IdleTimeout == 0 {
		r.IdleTimeout = 300
	}
}
